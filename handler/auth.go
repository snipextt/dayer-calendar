package handler

import (
	"encoding/base64"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/snipextt/dayer/utils"
	"golang.org/x/oauth2"
)

func GoogleAuthUrl(c *fiber.Ctx) error {
	defer HandleInternalServerError(c)

	cb := c.Query("cb")
	uid := c.Locals("uid").(string)
	state, err := utils.SetOAuthState(uid)
	utils.PanicOnError(err)

	state = uid + ":" + state + ";" + cb
	state = base64.StdEncoding.EncodeToString([]byte(state))
	authURL := utils.GoogleOauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)

	return HandleSuccess(c, nil, fiber.Map{
		"authURI": authURL,
	})
}

func MsAuthUrl(c *fiber.Ctx) error {
	clientID := "e04d0789-1705-421d-bd09-4cc5298abcc2"
	redirectURI := "http://localhost:3000/auth/microsoft/callback"
	scopes := []string{"Calendars.ReadWrite", "offline_access"}

	authURL := "https://login.microsoftonline.com/c7ac0230-4870-4cd9-860d-ba3c8833b969/oauth2/v2.0/authorize?" + url.Values{
		"response_type": {"code"},
		"response_mode": {"form_post"},
		"client_id":     {clientID},
		"redirect_uri":  {redirectURI},
		"state":         {"snipextt-was-here"},
	}.Encode() + "&scope=" + strings.Join(scopes, " ")

	return c.JSON(fiber.Map{
		"url": authURL,
	})
}
