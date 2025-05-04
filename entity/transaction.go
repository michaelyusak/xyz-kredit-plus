package entity

type Transaction struct {
	Id                int64   `json:"-"`
	AccountId         int64   `json:"-"`
	ContactNumber     string  `json:"contact_number" binding:"required"`
	OTR               float64 `json:"otr" binding:"required,gt=0"`
	InstallmentMonths int     `json:"installment_months" binding:"required,gte=1"`
	AdminFee          float64 `json:"admin_fee" binding:"required,gt=0"`
	TotalInstallemnt  float64 `json:"total_installment" binding:"required,gt=0"`
	TotalInterest     float64 `json:"total_interest" binding:"required,gt=0"`
	AssetName         string  `json:"asset_name" binding:"required"`
	CreatedAt         int64   `json:"created_at"`
	UpdatedAt         int64   `json:"updated_at"`
	DeletedAt         *int64  `json:"-"`
}
