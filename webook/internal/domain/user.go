package domain

import "time"

type User struct {
	Id          int64
	Nickname    string
	Email       string
	Password    string
	Birthday    time.Time
	Gender      int
	Description string
}

type Address struct {
}
