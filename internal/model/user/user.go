package model

import (
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID               uuid.UUID  `db:"id"`
	Email            string     `db:"email"`
	FirstName        *string    `db:"first_name"`
	LastName         *string    `db:"last_name"`
	PhoneNumber      *string    `db:"phone_number"`
	PhoneCountryCode string     `db:"phone_country_code"`
	Password         *string    `db:"password"`
	ProfilePicture   *string    `db:"profile_picture"`
	IsVerified       *bool      `db:"is_verified"`
	CreatedAt        time.Time  `db:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at"`
	DeletedAt        *time.Time `db:"deleted_at"`
	Dob              time.Time  `db:"dob"`
	Gender           *string    `db:"gender"`
	VerifyAt         *time.Time `db:"verify_at"`
}

func GenerateHashedPassword(password string, alg string) (res string, err error) {
	var b []byte
	switch alg {
	case "bcrypt":
		b, err = bcrypt.GenerateFromPassword([]byte(password), 0)
		if err != nil {
			return "", err
		}
		res = string(b)
	default:
	}

	return
}

func (r *User) IsValidCredential(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(*r.Password), []byte(password))
	return err == nil
}
