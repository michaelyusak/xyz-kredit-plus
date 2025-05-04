package entity

type Account struct {
	Id        int64  `json:"-"`
	Email     string `json:"email" example:"user@example.com" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	DeletedAt *int64  `json:"-"`
}

type LoginRegisterReq struct {
	Email     string `json:"email" example:"user@example.com" binding:"required,email"`
	Password  string `json:"password" example:"@abcD1234" binding:"required"`
}