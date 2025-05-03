package entity

type JwtClaims struct {
	AccountId      int64  `json:"account_id"`
	Email          string `json:"email"`
	IsKycCompleted bool   `json:"is_kyc_completed"`
}
