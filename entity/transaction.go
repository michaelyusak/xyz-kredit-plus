package entity

type Transaction struct {
	Id                int64   `json:"-"`
	AccountId         int64   `json:"-"`
	ContactNumber     string  `json:"contact_number" binding:"required"`
	OTR               float64 `json:"otr" binding:"required,gt=0"`
	InstallmentMonths int     `json:"installment_months" binding:"required,gte=1,lte=4"`
	AdminFee          float64 `json:"admin_fee" binding:"required,gt=0"`
	TotalInstallemnt  float64 `json:"total_installment" binding:"required,gt=0"`
	TotalInterest     float64 `json:"total_interest" binding:"required,gt=0"`
	AssetName         string  `json:"asset_name" binding:"required"`
	CreatedAt         int64   `json:"created_at"`
	UpdatedAt         int64   `json:"updated_at"`
	DeletedAt         *int64  `json:"-"`
}

type CreateTransactionReq struct {
	AccountId         int64   `json:"-"`
	ContactNumber     string  `json:"contact_number" example:"081312341234" binding:"required"`
	OTR               float64 `json:"otr" example:"1000000" binding:"required,gt=0"`
	InstallmentMonths int     `json:"installment_months" example:"2" binding:"required,gte=1,lte=4"`
	AdminFee          float64 `json:"admin_fee" example:"100000" binding:"required,gt=0"`
	TotalInstallemnt  float64 `json:"total_installment" example:"1000000000" binding:"required,gt=0"`
	TotalInterest     float64 `json:"total_interest" example:"999000000" binding:"required,gt=0"`
	AssetName         string  `json:"asset_name" example:"dog house" binding:"required"`
}