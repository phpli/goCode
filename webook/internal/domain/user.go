package domain

import "time"

type User struct {
	Id          int64
	Nickname    string
	Email       string
	Phone       string
	Password    string
	Birthday    time.Time
	Gender      int
	Description string
	Ctime       time.Time
}

type Address struct {
}
