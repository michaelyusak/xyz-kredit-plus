package entity

type Transaction struct {
	Id               int64   `json:"-"`
	AccountId        int64   `json:"-"`
	ContactNumber    string  `json:"contact_number"`
	OTR              float64 `json:"otr"`
	AdminFee         float64 `json:"admin_fee"`
	TotalInstallemnt float64 `json:"total_installemnt"`
	TotalInterest    float64 `json:"total_interest"`
	AssetName        string  `json:"asset_name"`
	CreatedAt        int64   `json:"created_at"`
	UpdatedAt        int64   `json:"updated_at"`
	DeletedAt        *int64  `json:"-"`
}
