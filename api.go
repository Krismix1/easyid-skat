// https://tutorialedge.net/golang/creating-restful-api-with-golang/
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/gorilla/mux"

	jose "github.com/dvsekhvalnov/jose2go"
	Rsa "github.com/dvsekhvalnov/jose2go/keys/rsa"
)

type taxAccount struct {
	Currency string
	Amount   float32
}

type jwtClaim struct {
	Email string
	Iat   int32
}

var accounts = map[string]taxAccount{
	"a@a.com": taxAccount{Currency: "DKK", Amount: 5000},
	"b@a.com": taxAccount{Currency: "DKK", Amount: 1000},
}

func printError(err error) {
	fmt.Println("Some error happened: ", reflect.TypeOf(err), err)
}

func sendErrorRes(w http.ResponseWriter, status int, msg string) (int, error) {
	response := map[string]string{"error": msg}
	resBytes, _ := json.Marshal(response)
	w.WriteHeader(status)
	return w.Write(resBytes)
}

func getAccountInfo(keyPath string) func(w http.ResponseWriter, r *http.Request) {
	keyBytes, err := ioutil.ReadFile(keyPath)

	if err != nil {
		panic("invalid key file")
	}

	publicKey, err := Rsa.ReadPublic(keyBytes)

	if err != nil {
		panic("invalid key format")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		authHeaders, ok := r.Header["Authorization"]
		if !ok {
			fmt.Println("No token found")
			sendErrorRes(w, 401, "Authorization required")
			return
		}
		// The header is only present if it contains non-whitespace characters
		// Therefore, we should not need to check that length of auth is > 0.
		bearer := authHeaders[0]
		token := strings.Split(bearer, " ")[1]

		payload, _, err := jose.Decode(token, publicKey)

		if err != nil {
			printError(err)
			sendErrorRes(w, 401, "Invalid JWT token")
			return
		}

		var jwtPayload jwtClaim
		json.Unmarshal([]byte(payload), &jwtPayload)
		email := jwtPayload.Email

		if account, ok := accounts[email]; ok {
			accountValue, _ := json.Marshal(account)
			w.Write(accountValue)
			return
		}

		sendErrorRes(w, 404, "Email/username not found")
	}
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/account", getAccountInfo("jwtRS256.key.pub")).Methods("GET")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	handleRequests()
}
