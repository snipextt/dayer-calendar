package timedoctor

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type TimeDoctorLoginResponseData struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expiresAt"`
}

type TimeDoctorLoginResponse struct {
	Data    TimeDoctorLoginResponseData `json:"data"`
	Error   string                      `json:"error"`
	Message string                      `json:"message"`
}

const baseurl = "https://api2.timedoctor.com/api/1.0"

func Login(email, password, totp string) (response *TimeDoctorLoginResponse, err error) {
	url := baseurl + "/login"
	body := map[string]string{
		"email":       email,
		"password":    password,
		"totpCode":    totp,
		"permissions": "read",
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return
	}

	rsp, err := http.Post(url, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return
	}

	defer rsp.Body.Close()

	err = json.NewDecoder(rsp.Body).Decode(&response)
	return
}
