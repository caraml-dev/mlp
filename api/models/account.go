package models

type Account struct {
	Id           Id          `json:"id"`
	UserId       Id          `json:"user_id`
	AccountType  AccountType `json:"account_type"`
	AccessToken  string      `json:"access_token"`
	TokenType    string      `json:"token_type"`
	RefreshToken string      `json:"refresh_token"`
	CreatedUpdated
}
