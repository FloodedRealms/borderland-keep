CREATE TABLE campaigns (
id INTEGER NOT NULL PRIMARY KEY,
name TEXT,
recruitment INTEGER,
judge TEXT,
timekeeping TEXT,
cadence TEXT,
created_at DATETIME NOT NULL,
updated_at DATETIME NOT NULL,
last_adventure DATETIME, 
api_user_id TEXT NOT NULL DEFAULT "",
FOREIGN KEY api_user_id REFERENCES api_users(id);