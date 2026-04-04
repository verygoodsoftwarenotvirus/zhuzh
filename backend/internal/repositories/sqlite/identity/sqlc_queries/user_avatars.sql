-- name: ArchiveUserAvatar :exec
UPDATE user_avatars SET
	archived_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
WHERE belongs_to_user = sqlc.arg(belongs_to_user)
	AND archived_at IS NULL;

-- name: CreateUserAvatar :exec
INSERT INTO user_avatars (
	id,
	belongs_to_user,
	uploaded_media_id
) VALUES (
	sqlc.arg(id),
	sqlc.arg(belongs_to_user),
	sqlc.arg(uploaded_media_id)
);
