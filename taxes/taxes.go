package taxes

import (
	"database/sql"
	"log"
	"strings"
)

// Tax holds information about taxes of a user
type Tax struct {
	Amount float32
}

// ForUser returns taxes information about a user based on the user email.
func ForUser(email string) (Tax, bool) {

	email = strings.TrimSpace(email)
	if len(email) <= 0 {
		// TODO: What do we do here?
		return Tax{}, false
	}

	db, err := sql.Open("sqlite3", "taxes.db")
	if err != nil {
		log.Fatal(err)
		return Tax{}, false
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT amount FROM taxes WHERE email = ?")
	if err != nil {
		log.Fatal(err)
		return Tax{}, false
	}
	defer stmt.Close()

	var amount float32
	err = stmt.QueryRow(email).Scan(&amount)
	if err != nil {
		// TODO: Handle email not found
		log.Fatal(err)
		return Tax{}, false
	}

	return Tax{Amount: amount}, true
}
