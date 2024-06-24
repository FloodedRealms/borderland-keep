CREATE TABLE campaign_activities (
id INTEGER PRIMARY KEY AUTOINCREMENT,
character_id INTEGER,
name text,
taken_at_level INTEGER,
xp_gained INTEGER,
FOREIGN KEY(character_id) REFERENCES characters(id)
);

