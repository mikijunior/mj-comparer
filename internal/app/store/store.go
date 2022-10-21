package store

type Store interface {
	Auth() AuthRepository
}