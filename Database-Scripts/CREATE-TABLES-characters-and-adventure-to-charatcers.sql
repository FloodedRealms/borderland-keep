CREATE TABLE characters (
id INTEGER NOT NULL PRIMARY KEY,
campaign_id INTEGER,
name TEXT,
current_xp INTEGER,
prime_req_percent INTEGER,
character_level INTEGER,
character_class TEXT,
created_at DATETIME NOT NULL,
updated_at DATETIME NOT NULL,
FOREIGN KEY(campaign_id) REFERENCES campaigns(id)
);

CREATE TABLE adventures_to_characters (
id INTEGER NOT NULL PRIMARY KEY,
adventure_id INTEGER,
character_id INTEGER,
half_share INTEGER NOT NULL CHECK (half_share IN (0,1)),
FOREIGN KEY(adventure_id) REFERENCES adventures(id),
FOREIGN KEY(character_id) REFERENCES characters(id),
UNIQUE(adventure_id, character_id)
);