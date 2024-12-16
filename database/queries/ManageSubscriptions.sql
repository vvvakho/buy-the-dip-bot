-- name: GetSubscription :one
SELECT * FROM user_subscriptions
WHERE user_id = $1 AND ticker = $2
LIMIT 1;

-- name: GetSubscriptionsByTicker :many
SELECT sub_id, user_id FROM user_subscriptions
WHERE ticker = $1 AND active = TRUE;

-- name: ReactivateSubscription :exec
UPDATE user_subscriptions
SET active = true, date_subscribed = NOW()
WHERE user_id = $1 AND ticker = $2;

-- name: AddSubscription :exec
INSERT INTO user_subscriptions (sub_id, user_id, ticker)
VALUES ($1, $2, $3);

-- name: Unsubscribe :exec
UPDATE user_subscriptions
SET active = false
WHERE user_id = $1 AND ticker = $2;

-- name: TempGetAllUsers :many
SELECT DISTINCT user_id FROM user_subscriptions;
