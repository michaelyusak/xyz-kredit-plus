package entity

type Account struct {
	Id        int64  `json:"-"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"-"`
	DeletedAt *int64  `json:"-"`
}
