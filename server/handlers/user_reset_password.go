package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Pharmeum/pharmeum-users-api/db"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/twinj/uuid"
)

type ResetPasswordRequest struct {
	Email string `json:"email"`
}

func (r ResetPasswordRequest) Validate() error {
	return validation.Validate(&r.Email, is.Email, validation.Required)
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	log := Log(r).WithField("handler", "reset_password")
	resetPasswordRequest := &ResetPasswordRequest{}
	if err := json.NewDecoder(r.Body).Decode(resetPasswordRequest); err != nil {
		w.Write(ErrResponse(http.StatusBadRequest, errors.New("not valid request body")))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := resetPasswordRequest.Validate(); err != nil {
		w.Write(ErrResponse(400, errors.New("invalid email length")))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := DB(r).GetUser(resetPasswordRequest.Email)
	if err != nil {
		log.WithError(err).Error("failed to get user")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(ErrResponse(http.StatusInternalServerError, err))
		return
	}

	//check if user no exist
	if user == nil {
		//return this by security reasons
		w.WriteHeader(http.StatusAccepted)
		return
	}

	token := uuid.NewV4()

	emailToken := &db.Token{
		Email:      resetPasswordRequest.Email,
		Token:      token.String(),
		LastSentAt: time.Now(),
	}

	if err := DB(r).CreateToken(emailToken); err != nil {
		Log(r).WithField("db", "email_tokens").WithError(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(ErrResponse(500, err))
		return
	}

	//link to web app new password form
	link := fmt.Sprintf("%s/recovery?token=%s", WebApp(r).String(), token.String())

	//skip err for Email client
	if err := EmailClient(r).Forgot(user.Email, link); err != nil {
		Log(r).WithField("google", "smtp").Error(fmt.Sprintf("failed to send forgot password email %s", err))
	}

	w.WriteHeader(http.StatusAccepted)
}
