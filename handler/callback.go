package handler

import (
	// "encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/snipextt/dayer/models"
	"github.com/snipextt/dayer/models/connection"
	"github.com/snipextt/dayer/utils"
)

func GoogleAuthCallback(c *fiber.Ctx) error {
	ctx, cancel := utils.NewContext()
	defer cancel()
	defer catchInternalServerError(c)

	// state := c.Query("state")
	// if state == "" {
	// 	return response.HandleBadRequest(c, "State is required")
	// }

	// decoded, err := base64.StdEncoding.DecodeString(state)
	// utils.PanicOnError(err)

	// state = string(decoded)
	// split := strings.Split(state, ";")
	//
	// uid, code := strings.Split(split[0], ":")[0], strings.Split(split[0], ":")[1]
	// codedb, err := utils.GetOAuthState(uid)
	// utils.PanicOnError(err)

	// if code != codedb {
	// 	return response.HandleUnauthorized(c, "Invalid state")
	// }

	tok, err := utils.GoogleOauthConfig.Exchange(ctx, c.Query("code"))
	utils.CheckError(err)

	email, err := utils.GetConnectedGoogleEmail(tok.RefreshToken)
	utils.CheckError(err)

	conn := connection.NewCalendarConnection(c.Locals("uid").(string), email, "google", tok.RefreshToken)
	err = conn.Save()
	utils.CheckError(err)

	return success(c, "Successfully connected to Google Calendar", conn.Id.Hex())
}

func MsAuthCallback(c *fiber.Ctx) error {
	code := c.FormValue("code")
	tokenUrl := "https://login.microsoftonline.com/c7ac0230-4870-4cd9-860d-ba3c8833b969/oauth2/v2.0/token"

	res, err := http.Post(tokenUrl, "application/x-www-form-urlencoded", strings.NewReader(url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"client_id":     {"e04d0789-1705-421d-bd09-4cc5298abcc2"},
		"access_type":   {"offline"},
		"redirect_uri":  {"http://localhost:3000/auth/microsoft/callback"},
		"client_secret": {"8tS8Q~oS.Z1SBfasrN5ErcVhPDedNleuwr2SNaSa"},
	}.Encode()))
	if err != nil {
		return err
	}
	var authres models.OuthResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&authres)
	if err != nil {
		return err
	}
	return c.JSON(authres)
}
