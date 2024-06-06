CREATE TABLE coins (
id INTEGER NOT NULL PRIMARY KEY,
adventure_id INTEGER,
copper INTEGER,
silver INTEGER,
electrum INTEGER,
gold INTEGER,
platinum INTEGER,
FOREIGN KEY(adventure_id) REFERENCES adventures(id)
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
apparent_value INTEGER,
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