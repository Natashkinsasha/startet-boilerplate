package handler

type TokenBody struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type tokenOutput struct {
	Body TokenBody
}
