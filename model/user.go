package model

import "time"

type User struct {
	Username    string    `json:"-"`
	DateOfBirth time.Time `json:"dateOfBirth"`
}
