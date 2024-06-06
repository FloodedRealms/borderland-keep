CREATE TABLE adventures (
id INTEGER NOT NULL PRIMARY KEY,
campaign_id INTEGER,
name TEXT,
adventure_date DATETIME,
created_at DATETIME NOT NULL,
updated_at DATETIME NOT NULL,
FOREIGN KEY(campaign_id) REFERENCES campaigns(id)
);