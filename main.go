package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type User struct {
	Username    string    `json:"-"`
	DateOfBirth time.Time `json:"dateOfBirth"`
}

var (
	users     = make(map[string]User)
	usersLock sync.Mutex
)

func putHandler(w http.ResponseWriter, r *http.Request) {
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

	usersLock.Lock()
	defer usersLock.Unlock()

	users[username] = User{
		Username:    username,
		DateOfBirth: dateOfBirth,
	}

	w.WriteHeader(http.StatusNoContent)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Path[len("/hello/"):]
	if username == "" {
		http.Error(w, "Username not provided", http.StatusBadRequest)
		return
	}

	usersLock.Lock()
	user, found := users[username]
	usersLock.Unlock()

	if !found {
		http.NotFound(w, r)
		return
	}

	daysUntilBirthday := daysUntil(user.DateOfBirth)
	message := fmt.Sprintf("Hello, %s! Your birthday is in %d day(s)", username, daysUntilBirthday)

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
	target = target.AddDate(now.Year(), 0, 0) // Set target's year to current year
	if target.Before(now) {
		target = target.AddDate(1, 0, 0) // If birthday already passed this year, set target's year to next year
	}
	return int(target.Sub(now).Hours() / 24)
}

func main() {
	http.HandleFunc("/hello/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			putHandler(w, r)
		case http.MethodGet:
			getHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server listening on :8080...")
	http.ListenAndServe(":8080", nil)
}
