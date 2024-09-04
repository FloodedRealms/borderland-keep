CREATE TABLE api_users (
id TEXT NOT NULL PRIMARY KEY,
api_key TEXT NOT NULL DEFAULT "",
friendly_name TEXT NOT NULL DEFAULT "",
campaign_number_limited INTEGER NOT NULL CHECK (campaign_number_limited IN (0,1)) DEFAULT 1,
salt TEXT,
UNIQUE(api_key)
);
CREATE TABLE adventures (
id INTEGER NOT NULL PRIMARY KEY,
campaign_id INTEGER,
name TEXT,
adventure_date DATETIME,
created_at DATETIME NOT NULL,
updated_at DATETIME NOT NULL, copper INTEGER DEFAULT 0, silver INTEGER DEFAULT 0, electrum INTEGER DEFAULT 0, gold INTEGER DEFAULT 0, platinum INTEGER DEFAULT 0, duration INT DEFAULT 1,
FOREIGN KEY(campaign_id) REFERENCES "_table1_old"(id)
);
CREATE TABLE adventures_to_characters (
id INTEGER NOT NULL PRIMARY KEY,
adventure_id INTEGER,
character_id INTEGER,
half_share INTEGER NOT NULL CHECK (half_share IN (0,1)),
xp_gained INTEGER NOT NULL DEFAULT 0,
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
CREATE TABLE systems(
id int not null primary key,
system_name text not null DEFAULT ""
);
CREATE TABLE classes(
id int not null primary key,
system_id int,
class_name text not null default "",
FOREIGN KEY(system_id) REFERENCES systems(id)
);
CREATE TABLE class_level_thresholds(
id int not null primary key,
class_id int,
xp_level int not null default 10000000,
xp_amount int not null default 10000000,
FOREIGN KEY(class_id) REFERENCES classes(id)
);
CREATE TABLE campaigns
(
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
system_id int DEFAULT 1, password TEXT, salt TEXT,
FOREIGN KEY(api_user_id) REFERENCES api_users(id)
FOREIGN KEY(system_id) REFERENCES systems(id)

);
CREATE TABLE campaign_activities (
id INTEGER PRIMARY KEY AUTOINCREMENT,
character_id INTEGER,
name text,
taken_at_level INTEGER,
xp_gained INTEGER,
FOREIGN KEY(character_id) REFERENCES characters(id)
);
CREATE TABLE sqlite_sequence(name,seq);
CREATE VIEW adventures_to_character_name AS SELECT atc.*, c.name FROM adventures_to_characters atc
JOIN "characters" c ON c.id = atc.character_id
/* adventures_to_character_name(id,adventure_id,character_id,half_share,xp_gained,name) */;
CREATE VIEW possible_characters_for_adventure AS
SELECT 
	a.id as adventure_id, 
    ch.id as character_id,
    ch.name as character_name, -- assuming there's a characterName field
    CASE 
        WHEN atc.id IS NOT NULL THEN 'Yes'
        ELSE 'No'
    END AS on_adventure
FROM adventures a 
LEFT JOIN "characters" ch on ch.campaign_id = a.campaign_id 
LEFT JOIN 
    adventures_to_characters atc ON atc.character_id = ch.id AND atc.adventure_id = a.id
ORDER BY a.id ASC
/* possible_characters_for_adventure(adventure_id,character_id,character_name,on_adventure) */;
CREATE VIEW campaign_to_class_options AS 
SELECT c.id as campaign_id, cl.id as class_id, cl.class_name as class_name
FROM campaigns c 
JOIN classes cl ON cl.system_id = c.system_id
/* campaign_to_class_options(campaign_id,class_id,class_name) */;
CREATE TABLE IF NOT EXISTS "character_statuses" (
id INTEGER NOT NULL PRIMARY KEY,
name TEXT
);
CREATE TABLE IF NOT EXISTS "characters" (
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
CREATE VIEW "character_campaign_view" AS 
WITH character_total_xp AS (
SELECT
c.campaign_id,
c.id, 
c.name, 
c.status_id,
c.class_id,
c.prime_req_percent as preq,
COALESCE(SUM(atc.xp_gained), 0) as adventure_xp, 
COALESCE(SUM(ca.xp_gained), 0) as campaign_xp,
COALESCE(SUM(atc.xp_gained), 0) + COALESCE(SUM(ca.xp_gained), 0) as total_xp
FROM "characters" c
LEFT JOIN adventures_to_characters atc on c.id = atc.character_id
LEFT JOIN campaign_activities ca on c.id = ca.character_id
GROUP BY c.id
)
SELECT 
t.campaign_id,
t.id, 
t.name, 
cs.name as status, 
cl.class_name,
t.preq,
t.adventure_xp, 
t.campaign_xp, 
t.total_xp, 
COALESCE(MAX(clt.xp_level), 1) AS level
FROM character_total_xp t
LEFT JOIN class_level_thresholds clt on t.class_id = clt.class_id AND t.total_xp >= clt.xp_amount
LEFT JOIN classes cl ON t.class_id = cl.id
LEFT JOIN character_statuses cs ON t.status_id = cs.id
GROUP BY t.id
/* character_campaign_view(campaign_id,id,name,status,class_name,preq,adventure_xp,campaign_xp,total_xp,level) */;
