package types

type User struct {
	ID       int64  `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Timezone string `json:"timezone"`
}

// time.LoadLocation("EST")
