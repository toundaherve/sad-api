package onboarding

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type InMemoryStore interface {
	SaveVerificationCode(code, email string) error
	GetVerificationCode(email string) (string, error)
}

type CodeStore map[string]string

func (s CodeStore) SaveVerificationCode(email, code string) error {
	s[email] = code
	return nil
}

func (s CodeStore) GetVerificationCode(email string) (string, error) {
	return s[email], nil
}

type OnboardingHandler struct {
	InMemoryStore InMemoryStore
	Logger        *logrus.Logger
}

func New(i InMemoryStore, l *logrus.Logger) *OnboardingHandler {
	onboardingHandler := &OnboardingHandler{
		InMemoryStore: i,
		Logger:        l,
	}

	if i == nil {
		onboardingHandler.InMemoryStore = CodeStore{}
	}

	return onboardingHandler
}

func (o *OnboardingHandler) BeginVerification(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	code := strconv.Itoa(generateCode())

	if err := o.InMemoryStore.SaveVerificationCode(email, code); err != nil {
		respondInternalServerError(w)
		o.Logger.WithField("email", email).Warn("Failed to save verification code")
		return
	}

	respondNoContent(w)

	fmt.Printf("Email: %s \nCode:  %s\n\n", email, code)
}

type verificationData struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

func (o *OnboardingHandler) VerifyCode(w http.ResponseWriter, r *http.Request) {
	var data verificationData

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		respondBadRequest(w, "The data you sent is not valid json format.")
		return
	}

	actualCode, err := o.InMemoryStore.GetVerificationCode(data.Email)
	if err != nil {
		respondInternalServerError(w)
		o.Logger.WithField("email", data.Email).Warn("Failed to get verification code")
		return
	}

	if actualCode == "" {
		respondBadRequest(w, "This email has no verification code. Try signing up again")
		return
	}

	if actualCode != data.Code {
		respondCode(w, false, "Wrong code")
		return
	}

	respondCode(w, true, "Good code")
}

func generateCode() int {
	min := 600000
	max := 699999
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func respondInternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w)
}

func respondBadRequest(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(struct {
		Code int
		Msg  string
	}{
		Code: http.StatusBadRequest,
		Msg:  msg,
	})
}

func respondNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
	fmt.Fprint(w)
}

func respondCode(w http.ResponseWriter, valid bool, msg string) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Valid bool   `json:"valid"`
		Msg   string `json:"msg"`
	}{
		Valid: valid,
		Msg:   msg,
	})
}
