package zap_logger

import (
	"fmt"
	logger "go_project/infrastructure/logging"
	"os"
	"path/filepath"
	"strings"
	"time"

	"bytes"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type prefixWriteSyncer struct {
	underlying zapcore.WriteSyncer
	prefix     []byte
	mu         sync.Mutex
}

func newPrefixWriteSyncer(
	ws zapcore.WriteSyncer,
	serviceName string,
) zapcore.WriteSyncer {
	return &prefixWriteSyncer{
		underlying: ws,
		prefix:     []byte(fmt.Sprintf("[%s] ", serviceName)),
	}
}

func (p *prefixWriteSyncer) Write(data []byte) (int, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	lines := bytes.SplitAfter(data, []byte("\n"))
	var output bytes.Buffer

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		output.Write(p.prefix)
		output.Write(line)
	}

	_, err := p.underlying.Write(output.Bytes())
	if err != nil {
		return 0, err
	}

	return len(data), nil
}

func (p *prefixWriteSyncer) Sync() error {
	return p.underlying.Sync()
}

type zapAdapter struct {
	logger *zap.SugaredLogger // Используем SugaredLogger для поддержки args ...any (Debugw, Infow)
	file   *os.File
}

func NewLogger(config logger.Config) (logger.Logger, error) {
	// 1. Уровень логирования
	zapLvl := zap.NewAtomicLevel()
	if err := zapLvl.UnmarshalText([]byte(config.Level)); err != nil {
		return nil, fmt.Errorf("unmarshal log level: %w", err)
	}

	// 2. Создание директории и файла для логов (если указана папка)
	var logFile *os.File
	var err error
	if config.Folder != "" {
		if err := os.MkdirAll(config.Folder, 0755); err != nil {
			return nil, fmt.Errorf("mkdir logFolder: %w", err)
		}

		timestamp := time.Now().UTC().Format("2006-01-02T15-04-05.000000")
		logFilePath := filepath.Join(
			config.Folder,
			fmt.Sprintf("%s.log", timestamp),
		)
		logFile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("open log file: %w", err)
		}
	}

	// 3. Выбор энкодера (JSON vs CONSOLE)
	var encoder zapcore.Encoder

	if strings.ToUpper(config.Format) == "JSON" {
		jsonConfig := zap.NewProductionEncoderConfig()
		jsonConfig.TimeKey = "timestamp"
		jsonConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		encoder = zapcore.NewJSONEncoder(jsonConfig)
	} else {
		// Консольный формат: имя сервиса в начале строки
		consoleConfig := zap.NewDevelopmentEncoderConfig()

		consoleConfig.EncodeTime = func(
			t time.Time,
			enc zapcore.PrimitiveArrayEncoder,
		) {
			enc.AppendString(t.Format("2006-01-02T15:04:05.000000"))
		}

		encoder = zapcore.NewConsoleEncoder(consoleConfig)
	}

	// 4. Настройка приёмников вывода (stdout + файл)
	var cores []zapcore.Core

	stdout := zapcore.AddSync(os.Stdout)

	if strings.ToUpper(config.Format) != "JSON" {
		stdout = newPrefixWriteSyncer(
			stdout,
			strings.ToUpper(config.ServiceName),
		)
	}

	cores = append(
		cores,
		zapcore.NewCore(encoder, stdout, zapLvl),
	)

	if logFile != nil {
		fileWriter := zapcore.AddSync(logFile)

		if strings.ToUpper(config.Format) != "JSON" {
			fileWriter = newPrefixWriteSyncer(
				fileWriter,
				strings.ToUpper(config.ServiceName),
			)
		}

		cores = append(
			cores,
			zapcore.NewCore(encoder, fileWriter, zapLvl),
		)
	}

	// if logFile != nil {
	// 	cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(logFile), zapLvl))
	// }

	core := zapcore.NewTee(cores...)

	// 5. Создание базвого логгера
	baseLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	// Если выбран JSON, добавляем имя сервиса как отдельное поле в JSON-структуру
	if strings.ToUpper(config.Format) == "JSON" {
		baseLogger = baseLogger.With(zap.String("service", config.ServiceName))
	}

	return &zapAdapter{
		logger: baseLogger.Sugar(),
		file:   logFile,
	}, nil
}
func (z *zapAdapter) Debug(msg string, args ...any) { z.logger.Debugw(msg, args...) }
func (z *zapAdapter) Info(msg string, args ...any)  { z.logger.Infow(msg, args...) }
func (z *zapAdapter) Warn(msg string, args ...any)  { z.logger.Warnw(msg, args...) }
func (z *zapAdapter) Error(msg string, args ...any) { z.logger.Errorw(msg, args...) }

// Принимаем variadic args ...any и возвращаем интерфейс logger.Logger
func (z *zapAdapter) With(args ...any) logger.Logger {
	return &zapAdapter{
		logger: z.logger.With(args...),
		file:   z.file,
	}
}
func (z *zapAdapter) Close() error {
	_ = z.Sync()
	if z.file != nil {
		return z.file.Close()
	}
	return nil
}
func (z *zapAdapter) Sync() error {
	return z.logger.Sync()
}
