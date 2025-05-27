package database

func (q *Query) IfPrefixExists(prefix string) bool {
	var exists bool
	q.db.QueryRow("SELECT EXISTS(SELECT 1 FROM warehouse WHERE prefix = $1)", prefix).Scan(&exists)
	return exists
}
