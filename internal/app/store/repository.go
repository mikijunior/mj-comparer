package store

import (
)

type AuthRepository interface {
	KeyExists(token string) bool
}
