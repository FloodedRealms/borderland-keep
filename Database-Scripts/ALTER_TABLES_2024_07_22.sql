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

INSERT INTO systems(id, system_name) VALUES(1, 'Adventurer Conqueror King System');
INSERT INTO classes(id, system_id, class_name) VALUES(1, 1, 'Fighter');
INSERT INTO classes(id, system_id, class_name) VALUES(2, 1, 'Explorer');
INSERT INTO classes(id, system_id, class_name) VALUES(3, 1, 'Thief');
INSERT INTO classes(id, system_id, class_name) VALUES(4, 1, 'Mage');
INSERT INTO classes(id ,system_id, class_name) VALUES(5, 1, 'Crusader');
INSERT INTO classes(id, system_id, class_name) VALUES(6, 1, 'Venturer');
INSERT INTO classes(id, system_id, class_name) VALUES(7, 1, 'Assassin');
INSERT INTO classes(id, system_id, class_name) VALUES(8, 1, 'Barbarian');
INSERT INTO classes(id, system_id, class_name) VALUES(9, 1, 'Bard');
INSERT INTO classes(id, system_id, class_name) VALUES(10, 1, 'Bladedancer');
INSERT INTO classes(id, system_id, class_name) VALUES(11, 1, 'Paladin');
INSERT INTO classes(id, system_id, class_name) VALUES(12, 1, 'Preistess');
INSERT INTO classes(id, system_id, class_name) VALUES(13, 1, 'Shaman');
INSERT INTO classes(id, system_id, class_name) VALUES(14, 1, 'Warlock');
INSERT INTO classes(id, system_id, class_name) VALUES(15, 1, 'Witch');
INSERT INTO classes(id, system_id, class_name) VALUES(16, 1, 'Dwarven Craftpriest');
INSERT INTO classes(id, system_id, class_name) VALUES(17, 1, 'Dwarven Vaultguard');
INSERT INTO classes(id, system_id, class_name) VALUES(18, 1, 'Elven Nightblade');
INSERT INTO classes(id, system_id, class_name) VALUES(19, 1, 'Elven Spellsword');
INSERT INTO classes(id, system_id, class_name) VALUES(20, 1, 'Nobiran Wonderworker');
INSERT INTO classes(id, system_id, class_name) VALUES(21, 1, 'Zaharan Ruinguard');
PRAGMA foreign_keys=off;


ALTER TABLE campaigns RENAME TO _table1_old;

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
system_id int,
FOREIGN KEY(api_user_id) REFERENCES api_users(id)
FOREIGN KEY(system_id) REFERENCES systems(id)

);

INSERT INTO campaigns(id, name, recruitment, judge, timekeeping, cadence, created_at, updated_at, last_adventure, api_user_id) SELECT * FROM _table1_old;
UPDATE campaigns SET system_id =1;
DROP TABLE _table1_old;

PRAGMA foreign_keys=on;

