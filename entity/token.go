package entity

type Token struct {
	Token     string `json:"token"`
	ExpiredAt int64  `json:"expired_at"`
}

type TokenData struct {
	AccessToken  Token `json:"access_token"`
	RefreshToken Token `json:"refresh_token"`
}
