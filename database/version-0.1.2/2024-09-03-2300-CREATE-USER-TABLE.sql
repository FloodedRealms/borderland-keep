
CREATE TABLE users(
id INTEGER not null primary key,
name TEXT not null,
email TEXT not null,
campaigns_limited INTEGER NOT NULL CHECK (campaigns_limited IN (0,1)),
password TEXT not null,
salt TEXT not null
);