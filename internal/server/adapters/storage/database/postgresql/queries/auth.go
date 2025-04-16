package queries

const (
	InsertUser = `INSERT INTO users (id, login, password_hash) VALUES ($1, $2, $3);`
	GetUser    = `SELECT id, login, password_hash FROM users WHERE login = $1;`
)
