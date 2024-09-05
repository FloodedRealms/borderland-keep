DROP TABLE campaigns;

CREATE TABLE campaigns
(
id INTEGER NOT NULL PRIMARY KEY,
user_id INTEGER NOT NULL,
name TEXT,
recruitment INTEGER,
judge TEXT,
timekeeping TEXT,
cadence TEXT,
created_at DATETIME NOT NULL,
updated_at DATETIME NOT NULL,
last_adventure DATETIME, 
system_id int DEFAULT 1,
FOREIGN KEY(user_id) REFERENCES users(id),
FOREIGN KEY(system_id) REFERENCES systems(id)

);