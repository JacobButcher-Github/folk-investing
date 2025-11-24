--TABLES
CREATE TABLE IF NOT EXISTS user_stocks (
  user_id INTEGER NOT NULL,
  stock_id INTEGER NOT NULL,
  quantity INTEGER NOT NULL,
  PRIMARY KEY (user_id, stock_id),
  FOREIGN KEY (user_id) REFERENCES users (id),
  FOREIGN KEY (stock_id) REFERENCES stocks (id)
);
