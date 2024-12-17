-- +goose Up
CREATE TABLE rsi (
  rsi_id UUID UNIQUE NOT NULL,
  ticker TEXT NOT NULL,
  rsi FLOAT NOT NULL,
  date TIMESTAMP NOT NULL,
  PRIMARY KEY (rsi_id)
);

-- +goose Down
DROP TABLE IF EXISTS rsi;
