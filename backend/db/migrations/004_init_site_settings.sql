--TABLES
CREATE TABLE IF NOT EXISTS site_settings (
  number_of_events_visible INTEGER DEFAULT 10,
  value_symbol TEXT DEFAULT "$",
  event_label TEXT DEFAULT "instance",
  lockout_time_start DATETIME DEFAULT current_timestamp
);
