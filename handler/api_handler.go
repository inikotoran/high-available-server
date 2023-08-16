package handler

import (
	"encoding/json"
	"fmt"
	"github.com/inikotoran/high-available-server/model"
	"net/http"
	"sync"
	"time"
)

func NewHandler() *Handler {
	return &Handler{
		users: make(map[string]model.User),
	}
}

type Handler struct {
	users     map[string]model.User
	usersLock sync.Mutex
}

func (h *Handler) Put(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Path[len("/hello/"):]
	if username == "" {
		http.Error(w, "Username not provided", http.StatusBadRequest)
		return
	}

	var userData struct {
		DateOfBirth string `json:"dateOfBirth"`
	}

	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	dateOfBirth, err := time.Parse("2006-01-02", userData.DateOfBirth)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	h.usersLock.Lock()
	defer h.usersLock.Unlock()

	h.users[username] = model.User{
		Username:    username,
		DateOfBirth: dateOfBirth,
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Path[len("/hello/"):]
	if username == "" {
		http.Error(w, "Username not provided", http.StatusBadRequest)
		return
	}

	h.usersLock.Lock()
	user, found := h.users[username]
	h.usersLock.Unlock()

	if !found {
		http.NotFound(w, r)
		return
	}

	daysUntilBirthday := daysUntil(user.DateOfBirth)
	message := fmt.Sprintf("Hello, %s! Your birthday is in %d day(s)", username, daysUntilBirthday)
	if daysUntilBirthday == 0 {
		message = fmt.Sprintf("Hello, %s! Happy birthday!", username)
	}

	responseData := map[string]string{
		"message": message,
	}

	responseJSON, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

func daysUntil(target time.Time) int {
	now := time.Now()
	target = time.Date(now.Year(), target.Month(), target.Day(), 0, 0, 0, 0, now.Location())
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	if target.Before(now) {
		target = target.AddDate(1, 0, 0)
	}
	return int(target.Sub(now).Hours() / 24)
}
