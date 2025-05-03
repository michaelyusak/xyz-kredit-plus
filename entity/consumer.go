package entity

type Consumer struct {
	Id                int64  `json:"-"`
	AccountId         int64  `json:"-"`
	IdentityNumber    string `json:"nik" binding:"required" validate:"required"`
	FullName          string `json:"full_name" binding:"required" validate:"required"`
	LegalName         string `json:"legal_name" binding:"required" validate:"required"`
	PlaceOfBirth      string `json:"place_of_birth" binding:"required" validate:"required"`
	DateOfBirth       string `json:"date_of_birth" binding:"required" validate:"required"`
	Salary            int64  `json:"salary" binding:"required" validate:"required"`
	IdentityCardPhoto Media  `json:"identity_card_photo"`
	SelfiePhoto       Media  `json:"selfie_photo"`
	CreatedAt         int64  `json:"-"`
	UpdatedAt         int64  `json:"-"`
	DeletedAt         *int64 `json:"-"`
}
