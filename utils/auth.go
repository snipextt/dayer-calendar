package utils

import (
	"golang.org/x/oauth2"
	goauth "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

func GetConnectedGoogleEmail(token string) (email string, err error) {
	ctx, cancel := NewContext()
	defer cancel()

	client := GoogleOauthConfig.Client(ctx, &oauth2.Token{
		RefreshToken: token,
	})
	srv, err := goauth.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return
	}

	u, err := srv.Userinfo.V2.Me.Get().Do()
	if err != nil {
		return
	}

	email = u.Email

	return
}
