package sqlstore

type AuthRepository struct {
	store *Store
}

func (r *AuthRepository) KeyExists(token string) (exists bool) {
	r.store.db.QueryRow(
		"SELECT id FROM auth WHERE `api-key` = ?",
		token,
	).Scan(&exists)

	return
}
