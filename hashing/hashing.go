package hashing

type Hashing interface {
	HashPassword(password string) (string, error)
	ComparePasswordWithHash(password, hash string) bool
}
