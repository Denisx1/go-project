BEGIN;

-- Удаляем колонку request_id из таблицы orders.outbox, если она существует
ALTER TABLE orders.outbox 
DROP COLUMN IF EXISTS request_id;

COMMIT;