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
	//不要组合，以后可能有叮叮info
	WechatInfo WechatInfo
}

//type Address struct {
//}
