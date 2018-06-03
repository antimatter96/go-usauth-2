package handlers

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/asaskevich/govalidator"
	
	"golang.org/x/crypto/bcrypt"
)

// SignupHandlerGet
func SignupHandlerGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	signupTemplate.Execute(w, nil)
}

// SignupHandlerPost
func SignupHandlerPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	errParseForm := r.ParseForm()
	if errParseForm != nil {
		fmt.Println(errParseForm)
		signupTemplate.Execute(w, constErrInternalError)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || !govalidator.IsEmail(email) {
		signupTemplate.Execute(w, constErrEmailMissing)
		return
	}

	if password == "" {
		signupTemplate.Execute(w, constErrPasswordMissing)
		return
	}

	userPresent, errPresent := userService.CheckUser(email)
	if errPresent != nil {
		signupTemplate.Execute(w, constErrInternalError)
		return
	}

	if userPresent {
		signupTemplate.Execute(w, constErrEmailTaken)
		return
	}

	hashedString, errBcrypt := getHashedPassword(password)
	if errBcrypt != nil {
		signupTemplate.Execute(w, constErrInternalError)
		return
	}

	if err := userService.AddUser(email, hashedString); err!=nil{
		signupTemplate.Execute(w, constErrInternalError)
		return
	}
	http.Redirect(w, r, "./login?success=true", http.StatusSeeOther)
}

func getHashedPassword(password string) (string, error) {
	passwordBytes := []byte(password)
	hashedBytes, errBcrypt := bcrypt.GenerateFromPassword(passwordBytes, bcryptCost)
	if errBcrypt != nil {
		fmt.Printf("Bcrypt error : %v", errBcrypt)
		return "", errBcrypt
	}
	return string(hashedBytes), nil
}