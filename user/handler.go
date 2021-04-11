package user

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	Validate *validator.Validate
	Storage  Storage
	logger   *logrus.Logger
}

func NewUserHandler(v *validator.Validate, s Storage, l *logrus.Logger) *UserHandler {
	return &UserHandler{
		Validate: v,
		Storage:  s,
		logger:   l,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		respondBadRequest(w, "The data you sent is not json formatted.")
		return
	}
	defer r.Body.Close()

	if err := h.Validate.Struct(&newUser); err != nil {
		respondValidationError(w, err.(validator.ValidationErrors))
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 14)
	if err != nil {
		respondInternalServerError(w)
		h.logger.WithFields(logrus.Fields{
			"password": newUser.Password,
			"err":      err.Error(),
		}).Warn("Failed to create password hash")
		return
	}
	newUser.Password = string(hash)

	if err := h.Storage.CreateUser(&newUser); err != nil {
		respondInternalServerError(w)
		h.logger.WithFields(logrus.Fields{
			"user": newUser.Email,
			"err":  err.Error(),
		}).Warn("Failed to create new user in the storage")
		return
	}

	respondCreated(w, "Your account has been successfully created.")
}

type validationErrorResponse struct {
	Code   int               `json:"code"`
	Errors map[string]string `json:"errors"`
}

func respondValidationError(w http.ResponseWriter, vErrors validator.ValidationErrors) {
	formatted := map[string]string{}
	for _, v := range vErrors {
		formatted[v.Field()] = v.Tag()
	}
	w.WriteHeader(http.StatusBadRequest)
	resp := validationErrorResponse{
		Code:   http.StatusBadRequest,
		Errors: formatted,
	}
	json.NewEncoder(w).Encode(resp)
}

type response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func respondBadRequest(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	resp := response{
		Code:    http.StatusBadRequest,
		Message: msg,
	}
	json.NewEncoder(w).Encode(resp)
}

func respondInternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	resp := response{
		Code:    http.StatusInternalServerError,
		Message: "Sorry, we are having a problem... retry later!",
	}
	json.NewEncoder(w).Encode(resp)
}

func respondCreated(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusCreated)
	resp := response{
		Code:    http.StatusCreated,
		Message: msg,
	}
	json.NewEncoder(w).Encode(resp)
}

type emailAvailableResponse struct {
	Valid bool   `json:"valid"`
	Msg   string `json:"msg"`
	Taken bool   `json:"taken"`
}

func (h *UserHandler) CheckEmailAvailable(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")

	if !isValidEmail(email) {
		resp := emailAvailableResponse{
			Valid: false,
			Msg:   "This email is invalid.",
			Taken: false,
		}
		respondEmailAvailable(w, resp)
		return
	}

	user, err := h.Storage.GetUserByEmail(email)
	if err != nil {
		log.Println(err)
		respondInternalServerError(w)
		return
	}

	if user != nil {
		resp := emailAvailableResponse{
			Valid: false,
			Msg:   "Email has already been taken.",
			Taken: true,
		}
		respondEmailAvailable(w, resp)
		return
	}

	resp := emailAvailableResponse{
		Valid: true,
		Msg:   "Available!",
		Taken: false,
	}
	respondEmailAvailable(w, resp)
}

func respondEmailAvailable(w http.ResponseWriter, resp emailAvailableResponse) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func isValidEmail(email string) bool {
	return true
}
