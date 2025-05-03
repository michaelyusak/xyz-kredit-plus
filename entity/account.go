package entity

type Account struct {
	Id        int64  `json:"-"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"-"`
	DeletedAt int64  `json:"-"`
}
