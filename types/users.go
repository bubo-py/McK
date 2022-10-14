package types

type User struct {
	ID       int64  `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Timezone string `json:"timezone"` // E.g. Africa/Abidjan, Europe/London, Asia/Tokyo
}

// time.LoadLocation("EST")
