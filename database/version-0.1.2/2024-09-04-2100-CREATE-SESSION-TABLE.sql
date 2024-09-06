CREATE TABLE sessions (
index_id INTEGER not null primary key,
uuid TEXT not null,
user_id INTEGER not null,
user_name TEXT not null,
expiry_time DATETIME not null,
FOREIGN KEY(user_id) REFERENCES users(id)
);