package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"

	"github.com/Pharmeum/pharmeum-users-api/db"

	"github.com/twinj/uuid"
	"golang.org/x/crypto/bcrypt"
)

type SignupRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	DateOfBirth string `json:"date_of_birth"`
}

func (u SignupRequest) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.Required),
		validation.Field(&u.Name, validation.Required),
		validation.Field(&u.Phone, validation.Required),
	)
}

func Signup(w http.ResponseWriter, r *http.Request) {
	log := Log(r).WithField("handler", "user_signup")
	signupRequest := &SignupRequest{}
	if err := json.NewDecoder(r.Body).Decode(signupRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(ErrResponse(400, err))
		return
	}

	if err := signupRequest.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(ErrResponse(400, err))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signupRequest.Password), 8)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(ErrResponse(400, err))
		return
	}

	createdUser, err := DB(r).GetUser(signupRequest.Email)
	if err != nil && err != sql.ErrNoRows {
		log.WithError(err).Error("failed to get signupRequest")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(ErrResponse(500, errors.New("Internal error")))
		return
	}

	if createdUser != nil {
		_, _ = w.Write(ErrResponse(400, errors.New("signupRequest with this email is already registered")))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dbUser := &db.User{
		Name:        signupRequest.Name,
		Email:       signupRequest.Email,
		Password:    string(hashedPassword),
		Phone:       signupRequest.Phone,
		DateOfBirth: signupRequest.DateOfBirth,
	}

	if err := DB(r).CreateUser(dbUser); err != nil {
		Log(r).WithField("table", "users").Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := DB(r).GetUser(dbUser.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(ErrResponse(http.StatusBadRequest, errors.New("invalid email address")))
			return
		}

		log.WithError(err).Errorf("failed to get user by email %s", dbUser.Email)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token := uuid.NewV4()
	confirmToken := &db.Token{
		UserID:     user.ID,
		Token:      token.String(),
		LastSentAt: time.Now(),
	}

	if err := DB(r).CreateToken(confirmToken); err != nil {
		Log(r).WithError(err).Error("failed to create token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//TODO: wait until Zain provide valid HTML template
	//link := fmt.Sprintf("%s/confirm/confirm?token=%s&email=%s", HTTP(r).String(), token.String(), signupRequest.Email)
	//err = EmailClient(r).Signup(signupRequest.Email, link)
	//if err != nil {
	//	Log(r).WithField("smtp", "client").Error(err)
	//}

	w.WriteHeader(http.StatusCreated)
}
