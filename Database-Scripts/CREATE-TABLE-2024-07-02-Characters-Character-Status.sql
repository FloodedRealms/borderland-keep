CREATE TABLE "character_statuses" (
id INTEGER NOT NULL PRIMARY KEY,
name TEXT
);

CREATE TABLE "characters" (
id INTEGER NOT NULL PRIMARY KEY,
campaign_id INTEGER,
name TEXT,
status_id INTEGER,
prime_req_percent INTEGER,
class_id INTEGER,
created_at DATETIME NOT NULL,
updated_at DATETIME NOT NULL,
FOREIGN KEY(campaign_id) REFERENCES campaigns(id),
FOREIGN KEY(class_id) REFERENCES classes(id),
FOREIGN KEY(status_id) REFERENCES "character_statuses"(id)
);