package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const baseURL = "http://localhost:42069"
const contentType = "application/json"

// LoginRequest holds data needed for login request creation
type LoginRequest struct {
	SuccessURL string `json:"successUrl"`
	CancelURL  string `json:"cancelUrl"`
}

// LoginResponse contains information about procceding with a login request
type LoginResponse struct {
	RedirectURL string `json:"redirectUrl"`
}

// CreateLogin creates a login request with easyID API
func CreateLogin(body LoginRequest) (LoginResponse, error) {
	// TODO: Validate body data
	bodyJSON, _ := json.Marshal(body)

	url := fmt.Sprintf("%s/%s", baseURL, "login-request")
	res, err := http.Post(url, contentType, bytes.NewBuffer(bodyJSON))
	if err != nil {
		log.Fatal("Error reading response. ", err)
		return LoginResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatal("easyID responded with non-200 code")
		return LoginResponse{}, err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Error reading body. ", err)
		return LoginResponse{}, err
	}

	var redirectData LoginResponse
	err = json.Unmarshal(resBody, &redirectData)
	if err != nil {
		log.Fatal("Error parsing response body. ", err)
		return LoginResponse{}, err
	}
	return redirectData, nil
}
