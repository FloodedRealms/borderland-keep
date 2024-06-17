CREATE TABLE api_users (
id TEXT NOT NULL PRIMARY KEY,
api_key TEXT NOT NULL DEFAULT "",
friendly_name TEXT NOT NULL DEFAULT "",
campaign_number_limited INTEGER NOT NULL CHECK (campaign_number_limited IN (0,1)) DEFAULT 1,
salt TEXT,
UNIQUE(api_key)
);

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
FOREIGN KEY(api_user_id) REFERENCES api_users(id)
);

CREATE TABLE adventures (
id INTEGER NOT NULL PRIMARY KEY,
campaign_id INTEGER,
name TEXT,
adventure_date DATETIME,
created_at DATETIME NOT NULL,
updated_at DATETIME NOT NULL, copper INTEGER DEFAULT 0, silver INTEGER DEFAULT 0, electrum INTEGER DEFAULT 0, gold INTEGER DEFAULT 0, platinum INTEGER DEFAULT 0, duration INT DEFAULT 1,
FOREIGN KEY(campaign_id) REFERENCES campaigns(id)
);

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

CREATE TABLE gems (
id INTEGER NOT NULL PRIMARY KEY,
adventure_id INTEGER,
name TEXT,
description TEXT,
value INTEGER,
total INTEGER,
FOREIGN KEY(adventure_id) REFERENCES adventures(id)
);

CREATE TABLE jewellery (
id INTEGER NOT NULL PRIMARY KEY,
adventure_id INTEGER,
name TEXT,
description TEXT,
value INTEGER,
total INTEGER,
FOREIGN KEY(adventure_id) REFERENCES adventures(id)
);

CREATE TABLE magic_items (
id INTEGER NOT NULL PRIMARY KEY,
adventure_id INTEGER,
name TEXT,
description TEXT,
apparent_value INTEGER, actual_value INTEGER DEFAULT 0,
FOREIGN KEY(adventure_id) REFERENCES adventures(id)
);

CREATE TABLE monster_groups (
id INTEGER NOT NULL PRIMARY KEY,
adventure_id INTEGER,
monster_name TEXT,
number_defeated INTEGER,
xp_per_monster INTEGER,
total_xp INTEGER,
FOREIGN KEY(adventure_id) REFERENCES adventures(id)
);