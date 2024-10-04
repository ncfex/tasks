-- name: CreateTask :one
INSERT INTO tasks (id, description, is_completed, created_at, due_date)
VALUES (
    gen_random_uuid(),
    $1,
    FALSE,
    NOW(),
    $2
)
RETURNING *;

-- name: GetAllTasks :many
SELECT * 
FROM tasks;

-- name: GetTaskById :one
SELECT *
FROM tasks
WHERE id = $1;

-- name: GetTaskByPartialId :one
SELECT *
FROM tasks
WHERE id::text LIKE $1 || '%'
LIMIT 1;

-- name: CompleteTask :one
UPDATE tasks
SET is_completed = TRUE
WHERE id = $1
RETURNING *;

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = $1;