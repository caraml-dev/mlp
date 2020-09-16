package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gojek/mlp/models"
	rand "github.com/gojek/mlp/util"
)

type UsersController struct {
	*AppContext
}

var oauthStateString string

//AuthorizeUser redirects user to the OAuth2 URL for the user to authorize the application to be granted the required scopes.
//Return redirect URL with the authorization code
func (c *UsersController) AuthorizeUser(r *http.Request, args map[string]string, _ interface{}) *ApiResponse {
	oauthStateString = rand.String(20)
	url := c.GitlabService.Authorize(oauthStateString)
	return Ok(url)
}

//GenerateToken exchanges the authorization code with access token
func (c *UsersController) GenerateToken(r *http.Request, args map[string]string, body interface{}) *ApiResponse {
	ctx := r.Context()
	r.ParseForm()

	//TODO: fix validation on oauthStateString
	//      - using the current validation, two user can authorize at the same time and failed on generating tokens
	//state := r.Form.Get("state")
	//if state != oauthStateString {
	//	return Error(http.StatusInternalServerError, "Invalid Oauth State" + state  + oauthStateString)
	//}

	code := r.Form.Get("code")
	if code == "" {
		return Error(http.StatusBadRequest, "Code not found")
	}

	token, err := c.GitlabService.GenerateToken(ctx, code)
	if err != nil {
		fmt.Println(err)
		return Error(http.StatusInternalServerError, "Code exchange failed")
	}

	//Store generated token here
	user, err := c.GitlabService.GetUserInfo(token.AccessToken)
	savedUser, err := c.UsersService.Save(user)
	if savedUser == nil {
		return Error(http.StatusInternalServerError, "User is already present in the database")
	}
	if err != nil {
		return Error(http.StatusInternalServerError, err.Error())
	}

	//Build the user account
	userAccount := &models.Account{
		UserId:       savedUser.Id,
		AccessToken:  token.AccessToken,
		AccountType:  models.AccountTypes.Gitlab,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
	}

	_, err = c.AccountService.Save(userAccount)
	if err != nil {
		return Error(http.StatusInternalServerError, err.Error())
	}

	return Ok("Authorized")
}

func (c *UsersController) RetrieveToken(r *http.Request, _ map[string]string, body interface{}) *ApiResponse {
	ctx := context.Background()
	userEmail := r.Header.Get("User-Email")
	if userEmail == "" {
		return Error(http.StatusInternalServerError, "User's email is not provided in the request.")
	}

	user, err := c.UsersService.FindByEmail(userEmail)
	if err != nil {
		return Error(http.StatusInternalServerError, err.Error())
	}

	account, err := c.AccountService.FindByUserId(user.Id)
	if err != nil {
		return Error(http.StatusInternalServerError, err.Error())
	}

	accessToken, err := c.GitlabService.GetUserAccessToken(ctx, account)

	if err != nil {
		return Error(http.StatusInternalServerError, err.Error())
	}
	return Ok(accessToken)
}
