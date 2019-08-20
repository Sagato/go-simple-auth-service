package bcrypt

import "golang.org/x/crypto/bcrypt"

type BCrypt struct {
	Cost int
}

func New(cost int) *BCrypt {
	return &BCrypt{
		Cost: cost,
	}
}

func (b *BCrypt) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (b *BCrypt) ComparePasswordWithHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

