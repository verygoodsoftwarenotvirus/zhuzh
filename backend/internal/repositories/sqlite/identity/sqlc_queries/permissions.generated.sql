-- name: CreatePermission :exec
INSERT INTO permissions (
	id,
	name,
	description
) VALUES (
	sqlc.arg(id),
	sqlc.arg(name),
	sqlc.arg(description)
);

-- name: GetPermissionByID :one
SELECT
	permissions.id,
	permissions.name,
	permissions.description,
	permissions.created_at,
	permissions.last_updated_at,
	permissions.archived_at
FROM permissions
WHERE permissions.archived_at IS NULL
	AND permissions.id = sqlc.arg(id);

-- name: GetPermissionByName :one
SELECT
	permissions.id,
	permissions.name,
	permissions.description,
	permissions.created_at,
	permissions.last_updated_at,
	permissions.archived_at
FROM permissions
WHERE permissions.archived_at IS NULL
	AND permissions.name = sqlc.arg(name);

-- name: GetPermissions :many
SELECT
	permissions.id,
	permissions.name,
	permissions.description,
	permissions.created_at,
	permissions.last_updated_at,
	permissions.archived_at,
	(
		SELECT COUNT(permissions.id)
		FROM permissions
		WHERE permissions.archived_at IS NULL
			AND
			permissions.created_at > COALESCE(sqlc.narg(created_after), strftime('%Y-%m-%dT%H:%M:%fZ', 'now', '-999 years'))
			AND permissions.created_at < COALESCE(sqlc.narg(created_before), strftime('%Y-%m-%dT%H:%M:%fZ', 'now', '+999 years'))
			AND (
				permissions.last_updated_at IS NULL
				OR permissions.last_updated_at > COALESCE(sqlc.narg(updated_before), strftime('%Y-%m-%dT%H:%M:%fZ', 'now', '-999 years'))
			)
			AND (
				permissions.last_updated_at IS NULL
				OR permissions.last_updated_at < COALESCE(sqlc.narg(updated_after), strftime('%Y-%m-%dT%H:%M:%fZ', 'now', '+999 years'))
			)
	) AS filtered_count,
	(
		SELECT COUNT(permissions.id)
		FROM permissions
		WHERE permissions.archived_at IS NULL
	) AS total_count
FROM permissions
WHERE permissions.archived_at IS NULL
	AND permissions.created_at > COALESCE(sqlc.narg(created_after), strftime('%Y-%m-%dT%H:%M:%fZ', 'now', '-999 years'))
	AND permissions.created_at < COALESCE(sqlc.narg(created_before), strftime('%Y-%m-%dT%H:%M:%fZ', 'now', '+999 years'))
	AND (
		permissions.last_updated_at IS NULL
		OR permissions.last_updated_at > COALESCE(sqlc.narg(updated_after), strftime('%Y-%m-%dT%H:%M:%fZ', 'now', '-999 years'))
	)
	AND (
		permissions.last_updated_at IS NULL
		OR permissions.last_updated_at < COALESCE(sqlc.narg(updated_before), strftime('%Y-%m-%dT%H:%M:%fZ', 'now', '+999 years'))
	)
	AND permissions.id > COALESCE(sqlc.narg(cursor), '')
ORDER BY permissions.id ASC
LIMIT COALESCE(sqlc.narg(result_limit), 50);

-- name: ArchivePermission :execrows
UPDATE permissions SET archived_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE archived_at IS NULL AND id = sqlc.arg(id);
