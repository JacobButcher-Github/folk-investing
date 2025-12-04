--TABLES
CREATE TABLE IF NOT EXISTS site_settings (
  id INTEGER PRIMARY KEY CHECK (id = 1),
  number_of_events_visible INTEGER DEFAULT 10 NOT NULL,
  value_symbol TEXT DEFAULT '$' NOT NULL,
  event_label TEXT DEFAULT 'instance' NOT NULL,
  -- 0 = false (not locked), 1 = true (locked)
  lockout INTEGER DEFAULT 0 NOT NULL,
  lockout_time_start DATETIME DEFAULT current_timestamp NOT NULL
  --TODO: make this into a cronstring with a parser for every X amount of time.
  -- For now, this will work.
);
