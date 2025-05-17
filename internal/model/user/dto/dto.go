package dto

import (
	"time"

	"github.com/farganamar/evv-service/helpers"
	"github.com/farganamar/evv-service/helpers/failure"
	"github.com/farganamar/evv-service/helpers/logger"
	model "github.com/farganamar/evv-service/internal/model/user"
	"github.com/gofrs/uuid"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	IssuedAt     time.Time `json:"issued_at"`
}

type CreateNewUser struct {
	Email            string    `json:"email" validate:"required,email"`
	FirstName        string    `json:"first_name" validate:"required"`
	LastName         string    `json:"last_name"`
	PhoneNumber      string    `json:"phone_number" validate:"required"`
	PhoneCountryCode string    `json:"phone_country_code" validate:"required"`
	Password         string    `json:"password" validate:"required,min=8,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()"`
	Dob              time.Time `json:"dob" validate:"required,lt=now"`
	Gender           string    `json:"gender" validate:"required"`
	Interests        []string  `json:"interests" validate:"required,min=1,dive,required"`
}

type UpdateUser struct {
	ID               uuid.UUID `json:"id" validate:"required"`
	FirstName        string    `json:"first_name" validate:"required"`
	LastName         string    `json:"last_name"`
	PhoneNumber      string    `json:"phone_number" validate:"required"`
	PhoneCountryCode string    `json:"phone_country_code" validate:"required"`
	Dob              time.Time `json:"dob" validate:"required,lt=now"`
	IsVerified       bool      `json:"is_verified"`
}

func (d *CreateNewUser) Validate() (err error) {
	validator := helpers.GetValidator()
	return validator.Struct(d)
}

func (d CreateNewUser) ToModel(salt string) (res model.User, errr error) {
	id, _ := uuid.NewV7()
	isVerified := false

	password, err := model.GenerateHashedPassword(d.Password, "scrypt")
	if err != nil {
		logger.ErrorWithStack(err)
		failure.InternalError(err)
		return
	}

	res = model.User{
		ID:               id,
		Email:            d.Email,
		FirstName:        &d.FirstName,
		LastName:         &d.LastName,
		PhoneNumber:      &d.PhoneNumber,
		PhoneCountryCode: d.PhoneCountryCode,
		Password:         &password,
		Dob:              d.Dob,
		Gender:           &d.Gender,
		IsVerified:       &isVerified,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	return
}

func (d *UpdateUser) Validate() (err error) {
	validator := helpers.GetValidator()
	return validator.Struct(d)
}

func (d UpdateUser) ToModel() model.User {
	return model.User{
		ID:               d.ID,
		FirstName:        &d.FirstName,
		LastName:         &d.LastName,
		PhoneNumber:      &d.PhoneNumber,
		PhoneCountryCode: d.PhoneCountryCode,
		Dob:              d.Dob,
		UpdatedAt:        time.Now(),
		IsVerified:       &d.IsVerified,
	}
}

func (d *LoginRequest) Validate() (err error) {
	validator := helpers.GetValidator()
	return validator.Struct(d)
}
