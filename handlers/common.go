package handlers

import(
	"crypto/rand"
	"encoding/base64"

	"../db"
	"../constants"
)

// All the different errors
var (
	constErrEmailMissing        string = "Email Not Present"
	constErrPasswordMissing     string = "Password Not Present"
	constErrNotRegistered       string = "No records found"
	constErrInternalError       string = "An Error Occured"
	constErrPasswordMatchFailed string = "Passwords do not match"
	constErrEmailTaken          string = "Email Taken"
)

var userService db.UserAuthService

// Init is used to initialize all things
func Init(){
	initCookie()
	parseTemplates()

	userService = db.NewUserInterface()

	bcryptCost = constants.Value("bcrypt-cost").(int)
}

func generateRandomString(length int) (string, error) {
	x := make([]byte, length)
	_, err := rand.Read(x)

	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(x), err
}

var bcryptCost int