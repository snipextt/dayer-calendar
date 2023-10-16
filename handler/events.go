package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/snipextt/dayer/models"
)

func SyncEventsForMsCalendar(c *fiber.Ctx) error {
	token := "0.AXEAMAKsx3BI2UyGDbo8iDO5aYkHTeAFFx1CvQlMxSmKvMJxAPI.AgABAAEAAAD--DLA3VO7QrddgJg7WevrAgDs_wUA9P_R-kJgVHxlmHbWLZHfMu6-xX_eqtEETM5Nxo9O-i_dVO0lGfGzAh3FH4MbymGD5wMerSzUUdfqvQ-EYl6u8_9R97u1EA16IAdIM-nwZC5kAPqs3exOJkRBLOAiM9Cc7Cq0aT4vzULVpFjU7iKzOUKc11p2iAn9EHc3bvgzYaI87HBCOTRF37eB_6OLU80CHTAY2YcRmrkRImVfK_xkStbTRqx3IgVx3_trybigadE4H31n7sPQukb1uMZr3hrGACOp6RB5JYIzDdPLd1q9XxMEcDZCV8Y8AfjGns-AltmqHJqW1-UGtQNeE3zmXVGfN5CKYYw8ndpHxyMmjueBKcO3I18sV61PFwUhh-5XWv9T9BkUXWgFHvVIrQtuBRGkj_PFH9E-9h7xwJa-N96SiXTBJOsoVxppiW1HL5MogvXSE9Weuw7ANP6EGmAhP14jvWUX2PsyoNHyRAHD1bgve-mfbjQi58TBTvwzC3W3DDG5H3eaxLWgxhpD6zHvlJv1r2h40447nzF_8DRrVpFzCU06KJzbYlxkD1VtXWtrcVjuCS8H1P-sOSj-bLcqwZPnPgFQGdv0wKwlS17tvVsncIXYJDFrAD-_nkYzMxn-Fc2A_LA0u5ocaEQrSwfU_pLyVvpatSyF_dUR6SXoG4AmBV25Pg0YmVsAFdVkWKDyc_28_z94DOFTFJMu1wU1VbuxwgQWECorfBody6q-4TNAQtwOzwn-82CyMtLwSIG10Fdmzhs9js8qTyxk8k8Topk2k4_K5GbJBff4mjJM4ZxaD25PKhrGikJbS6NU31UQjTQj8ZIHK-aPrSafkPLBYK_nhzZi-QorBbStiC5fLXM75yTC122h1kBBZI_8pl3liTHzExDVGMeDtm-lQ1-fggUOBxkf4YhXby6XDpI6Z5xsihdOTtcw6trDI80yHO3MovywfY2OMWtf1r8Nb7WG_bfnQZg-rV0sdII3SblhGN_XuXxw0IDxfuy5qiBGasa0EOFLATSlE6yRWwD_Q-axwBRnhoTqpyhgvGloXqimQ_U_4o7sOZs7B370bRIdpEjYP9OE04hzqdP6jI-AQsX-KUTRv2bp4G9IrFGa3ec0VfCfZjpricbWyAeJS7EyLiyzl5GR_kQ4eN7aA97cyx2Uh1G7MgWlwRlWjfFmXnjlE8BaRWinbRcs7NZNiF4MPYC2aMctSBHx-WmI6pygP6wBJxwS_70dTNJqnke_EUms5Oy2jbgZLjTDxuwYGQPHRdAl7zfcHXk1MR1v7NfZAU-glw7zzK1redsz1qdTaFkLG0meiRqHXncNEO2jDtpMS_7zvO2mrWGB9aq7UxTZOA4wElAnoQow1S6HQ6mDLBkhVemscNYDkdNDaWI6uidifjmdGku6jFi4cpNcznaebTChsNad0bmSVi_tVjfWqrWNoYQjYZpAKZXc4k7U6WaPGc8FnptP3ocifUVorxAlqumWCT-3Odu_2rdAGGQF0L3rpV8PHS3UTke5zqRH2o912CdklZdYp6uBF9FxCxGkAOU11jQAZ_jdK3wtYilGqu5G"
	tokenUrl := "https://login.microsoftonline.com/c7ac0230-4870-4cd9-860d-ba3c8833b969/oauth2/v2.0/token"

	res, err := http.Post(tokenUrl, "application/x-www-form-urlencoded", strings.NewReader(url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {token},
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
	token = authres.AccessToken

	calApi := "https://graph.microsoft.com/v1.0/me/calendarview?startdatetime=2023-04-08T20:04:31.713Z&enddatetime=2023-04-15T20:04:31.713Z"

	req, err := http.NewRequest("GET", calApi, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	val, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return c.JSON(string(val))
}
