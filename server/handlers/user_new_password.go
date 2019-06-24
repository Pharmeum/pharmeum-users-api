package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func NewPassword(w http.ResponseWriter, r *http.Request) {
	log := Log(r).WithField("db", "email_tokens")
	token := r.URL.Query().Get("token")
	if token == "" {
		w.Write(ErrResponse(400, errors.New("failed to found token in request params")))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if token == "" {
		_, _ = w.Write(ErrResponse(http.StatusBadRequest, errors.New("empty token parameter in request")))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userToken, err := DB(r).GetUserByToken(token)
	if err != nil {
		log.WithError(err).Error("failed to get user token")
		w.Write(ErrResponse(http.StatusBadRequest, err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if userToken == nil {
		w.Write(ErrResponse(400, errors.New("no such user with provided userToken address")))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newCreds := &struct {
		Password string `json:"password"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(newCreds); err != nil {
		w.Write(ErrResponse(400, err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if newCreds.Password == "" {
		w.Write(ErrResponse(400, errors.New("empty password")))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newCreds.Password), 8)
	if err != nil {
		w.Write(ErrResponse(400, err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := DB(r).SetUserNewPassword(userToken.Email, string(hashedPassword)); err != nil {
		Log(r).WithField("user", "new password").Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := DB(r).DeleteToken(token); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//notify user about password changing
	if err := EmailClient(r).NewPassword(userToken.Email); err != nil {
		Log(r).WithField("email_client", "notification").Error("failed to send notification about new password to", userToken)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
