-- name: AddRSI :one
INSERT INTO rsi (rsi_id, ticker, rsi, date)
VALUES ($1, $2, $3, $4)
RETURNING *;


-- name: CheckRSI :one
SELECT rsi_id, ticker, rsi, date
FROM rsi
WHERE ticker = $1 AND date::DATE = $2::DATE
LIMIT 1;

