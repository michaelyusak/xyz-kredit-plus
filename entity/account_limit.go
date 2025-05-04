package entity

type AccountLimit struct {
	Id        int64   `json:"-"`
	AccountId int64   `json:"account_id"`
	Limit1M   float64 `json:"limit_1_m"`
	Limit2M   float64 `json:"limit_2_m"`
	Limit3M   float64 `json:"limit_3_m"`
	Limit4M   float64 `json:"limit_4_m"`
	CreatedAt int64   `json:"created_at"`
	UpdatedAt int64   `json:"updated_at"`
	DeletedAt *int64  `json:"-"`
}
