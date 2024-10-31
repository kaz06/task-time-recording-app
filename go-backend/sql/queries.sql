-- name: GetUsers :many
SELECT * FROM users;

-- name: GetUserByUID :one
SELECT id, uid, name, email FROM users WHERE uid = $1;

-- name: CreateUser :one
INSERT INTO users (uid, name, email)
VALUES ($1, $2, $3)
RETURNING id, uid, name, email;

-- name: GetTasksByUserID :many
SELECT * FROM tasks WHERE user_id = $1;

-- name: CreateTask :one
INSERT INTO tasks (title, task_time, task_finish_date, user_id) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetTasks :many
SELECT * FROM tasks;

-- name: GetTaskByID :one
SELECT * FROM tasks WHERE id = $1;

-- name: UpdateTask :one
UPDATE tasks SET title = $1, task_time = $2, task_finish_date = $3 WHERE id = $4 RETURNING *;

-- name: DeleteTask :one
DELETE FROM tasks WHERE id = $1 RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetTaskByUserIDAndTaskID :one
SELECT * FROM tasks WHERE user_id = $1 AND id = $2;

-- name: GetTags :many
SELECT * FROM tags;

-- name: GetTagNameByID :one
SELECT name FROM tags WHERE id = $1;

-- name: CreateTag :one
INSERT INTO tags (name) VALUES ($1) RETURNING *;

-- name: GetTagIdByName :one
SELECT id FROM tags WHERE name = $1;

-- name: CreateTaskTag :one
INSERT INTO task_tags (task_id, tag_id) VALUES ($1, $2) RETURNING *;

-- name: GetTaskTags :many
SELECT * FROM task_tags;

