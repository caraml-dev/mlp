package service

import (
	"context"
	"fmt"

	"github.com/xanzy/go-gitlab"
	"golang.org/x/oauth2"

	"github.com/gojek/mlp/api/models"
)

type GitlabService interface {
	Authorize(state string) string
	GenerateToken(ctx context.Context, code string) (*oauth2.Token, error)
	GetUserAccessToken(ctx context.Context, userAccount *models.Account) (string, error)
	GetUserInfo(accessToken string) (*models.User, error)
}

type gitlabService struct {
	host string
	cfg  *oauth2.Config
}

func NewGitlabService(host string, cfg *oauth2.Config) GitlabService {
	return &gitlabService{host: host, cfg: cfg}
}

//Authorize User
func (service *gitlabService) Authorize(state string) string {
	url := service.cfg.AuthCodeURL(state)
	return url
}

//Generate access token after receiving the authorization code from Gitlab server
func (service *gitlabService) GenerateToken(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := service.cfg.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (service *gitlabService) GetUserAccessToken(ctx context.Context, account *models.Account) (string, error) {
	//Check that the accessToken is still valid
	userToken := &oauth2.Token{
		AccessToken:  account.AccessToken,
		TokenType:    account.TokenType,
		RefreshToken: account.RefreshToken,
	}

	//Refresh token automatically when it is expired
	validTokenSource := service.cfg.TokenSource(ctx, userToken)
	validToken, _ := validTokenSource.Token()

	return validToken.AccessToken, nil
}

//Retrieve user's information
func (service *gitlabService) GetUserInfo(accessToken string) (*models.User, error) {
	//Refactored
	git := gitlab.NewOAuthClient(nil, accessToken)
	git.SetBaseURL(fmt.Sprintf("%s/api/v4", service.host))

	currentUser, _, err := git.Users.CurrentUser()
	user := models.User{
		Username: currentUser.Username,
		Email:    currentUser.Email,
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
