-- +goose Up
CREATE TABLE user_subscriptions (
  sub_id UUID UNIQUE NOT NULL,
  user_id BIGINT NOT NULL,
  ticker TEXT NOT NULL,
  date_subscribed TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  last_notified TIMESTAMP,
  active BOOLEAN DEFAULT TRUE,
  PRIMARY KEY (sub_id)
);

CREATE UNIQUE INDEX idx_user_ticker ON user_subscriptions (user_id, ticker);

-- +goose Down
DROP INDEX IF EXISTS idx_user_ticker;
DROP TABLE user_subscriptions;
