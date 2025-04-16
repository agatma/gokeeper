package queries

const (
	GetAllDataByUserID = `
		SELECT
			id,
			type,
			data,
			meta,
			saved_at
		FROM private
		WHERE user_id = $1
		LIMIT $2 OFFSET $3;
	`
	InsertData = `
		INSERT INTO private (id, type, data, meta, saved_at, user_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id, id)
		DO UPDATE SET
			type = $2,
			data = $3,
			meta = $4,
			saved_at = $5,
			updated_at = CURRENT_TIMESTAMP
		;
	`
	DeleteData  = `DELETE FROM private WHERE user_id = $1 AND id = $2;`
	GetDataByID = `
		SELECT
			type,
			data,
			meta,
			saved_at
		FROM private
		WHERE user_id = $1 AND id = $2;
	`
)
