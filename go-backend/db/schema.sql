CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    uid VARCHAR(128) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL
);


CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    task_time TIME NOT NULL,
    task_finish_date DATE NOT NULL,
    user_id INT NOT NULL REFERENCES users(id)
);

CREATE TABLE task_tags (
    task_id INT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    tag_id INT NOT NULL REFERENCES tags(id),
    PRIMARY KEY (task_id, tag_id)
);
