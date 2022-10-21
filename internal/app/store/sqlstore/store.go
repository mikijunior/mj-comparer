package sqlstore

import (
	"database/sql"
	"mjcomparer/internal/app/store"
)

type Store struct {
	db             *sql.DB
	authRepository *AuthRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) Auth() store.AuthRepository {
	if s.authRepository != nil {
		return s.authRepository
	}

	s.authRepository = &AuthRepository{
		store: s,
	}

	return s.authRepository
}
