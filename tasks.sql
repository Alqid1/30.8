DROP TABLE IF EXISTS tasks_labels,labels,tasks,users;

CREATE TABLE users(
id SERIAL PRIMARY KEY,
name TEXT NOT NULL
);

CREATE TABLE tasks(
id SERIAL PRIMARY KEY,
opened BIGINT NOT NULL DEFAULT extract(epoch from now()),
closed BIGINT DEFAULT 0,
author_id BIGINT REFERENCES users(id),
assigned_id BIGINT REFERENCES users(id),
title TEXT NOT NULL,
content TEXT NOT NULL
);

CREATE TABLE labels(
id SERIAL PRIMARY KEY,
name TEXT NOT NULL
);

CREATE TABLE tasks_labels(
task_id INTEGER REFERENCES tasks(id),
label_id INTEGER REFERENCES labels(id)
);

INSERT INTO users (id, name) VALUES (0, 'default');

select * from users;
