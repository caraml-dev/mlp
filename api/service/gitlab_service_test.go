package service

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"

	"github.com/gojek/mlp/models"
)

type mockGitlabService struct {
	mock.Mock
}

func (m *mockGitlabService) GetUserInfo(accessToken string) (*models.User, error) {
	args := m.Called(accessToken)
	return args.Get(0).(*models.User), args.Error(1)
}

func Test_gitlabService_GetUserInfo(t *testing.T) {
	type fields struct {
		host string
		cfg  *oauth2.Config
	}
	type args struct {
		accessToken string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.User
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "success",
			fields: fields{
				host: "https://gitlab.com",
				cfg: &oauth2.Config{
					ClientID:     "TestClient",
					ClientSecret: "TestSecret",
					Scopes:       []string{"read_user"},
					RedirectURL:  "TestRedirect",
					Endpoint: oauth2.Endpoint{
						AuthURL:  "https://gitlab.com/oauth/authorize",
						TokenURL: "https://gitlab.com/oauth/token",
					},
				},
			},
			args: args{
				accessToken: "abcde",
			},
			want: &models.User{
				Id:       1,
				Username: "user",
				Email:    "user@test.com",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := new(mockGitlabService)
			service.On("GetUserInfo", "abcde").Return(&models.User{Id: 1, Username: "user", Email: "user@test.com"}, nil)

			got, err := service.GetUserInfo(tt.args.accessToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("gitlabService.GetUserInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("gitlabService.GetUserInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_gitlabService_Authorize(t *testing.T) {
	type fields struct {
		host string
		cfg  *oauth2.Config
	}
	type args struct {
		state string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
		{
			"success",
			fields{
				host: "https://gitlab.com",
				cfg: &oauth2.Config{
					ClientID:     "TestClient",
					ClientSecret: "TestSecret",
					Scopes:       []string{"read_user"},
					RedirectURL:  "TestRedirect",
					Endpoint: oauth2.Endpoint{
						AuthURL:  "https://gitlab.com/oauth/authorize",
						TokenURL: "https://gitlab.com/oauth/token",
					},
				},
			},
			args{
				state: "xyz",
			},
			"https://gitlab.com/oauth/authorize?client_id=TestClient&redirect_uri=TestRedirect&response_type=code&scope=read_user&state=xyz",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &gitlabService{
				host: tt.fields.host,
				cfg:  tt.fields.cfg,
			}
			if got := service.Authorize(tt.args.state); got != tt.want {
				t.Errorf("gitlabService.Authorize() = %v, want %v", got, tt.want)
			}
		})
	}
}
