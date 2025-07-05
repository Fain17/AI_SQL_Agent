-- name: CreateFile :one
INSERT INTO files (filename, content, embedding)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetFile :one
SELECT * FROM files WHERE id = $1;

-- name: GetAllFiles :many
SELECT * FROM files ORDER BY id DESC;

-- name: GetFilesByFilename :many
SELECT * FROM files
WHERE filename ILIKE '%' || $1 || '%'
ORDER BY id DESC;

-- name: GetFileMetadata :many
SELECT id, filename, LENGTH(content) AS size, created_at
FROM files
ORDER BY created_at DESC;


-- name: GetFilesByDateRange :many
SELECT * FROM files
WHERE created_at BETWEEN $1 AND $2
ORDER BY created_at DESC;


-- name: UpdateFile :one
UPDATE files
  SET filename = $2, content = $3, embedding = $4
WHERE id = $1
RETURNING *;

-- name: DeleteFile :exec
DELETE FROM files WHERE id = $1;

-- name: SoftDeleteFile :exec
UPDATE files SET deleted = TRUE WHERE id = $1;

-- name: UndoSoftDelete :exec
UPDATE files SET deleted = FALSE WHERE id = $1;

-- name: GetDeletedFiles :many
SELECT * FROM files WHERE deleted = TRUE ORDER BY created_at DESC;

-- name: CountTotalFiles :one
SELECT COUNT(*) FROM files;