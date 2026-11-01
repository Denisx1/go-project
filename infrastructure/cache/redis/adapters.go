package redisclient

import (
	"fmt"

	"github.com/bytedance/sonic"
	redis "github.com/redis/go-redis/v9"
)

type stringCMD struct {
	cmd *redis.StringCmd
}

type intCMD struct {
	cmd *redis.IntCmd
}
type statusCMD struct {
	status *redis.StatusCmd
}

func (s stringCMD) Result() (string, error) {
	val, err := s.cmd.Result()
	if err != nil {
		return "", mapErrors(err)
	}
	return val, nil
}
func (s statusCMD) Err() error {

	return mapErrors(s.status.Err())
}

func (s stringCMD) Bytes() ([]byte, error) {
	val, err := s.cmd.Bytes()
	if err != nil {
		return nil, mapErrors(err)
	}
	return val, nil
}

func (s stringCMD) Scan(dest any) error {
	bytes, err := s.Bytes()
	if err != nil {
		return mapErrors(err)
	}
	if err := sonic.Unmarshal(bytes, dest); err != nil {
		return fmt.Errorf("failed to scan redis payload: %w", err)
	}
	return nil
}

// Err проверяет наличие ошибки без получения значения
func (s stringCMD) Err() error {

	return mapErrors(s.cmd.Err())
}

func (i intCMD) Result() (int64, error) {

	val, err := i.cmd.Result()
	if err != nil {
		return 0, mapErrors(err)
	}
	return val, nil
}
func (i intCMD) Uint64() (uint64, error) {

	val, err := i.cmd.Uint64()
	if err != nil {
		return 0, mapErrors(err)
	}
	return val, nil
}

func (i intCMD) Err() error {

	return mapErrors(i.cmd.Err())
}

func (a statusCMD) Result() (string, error) {
	data, err := a.status.Result()
	if err != nil {
		return "", mapErrors(err)
	}

	return data, nil
}
