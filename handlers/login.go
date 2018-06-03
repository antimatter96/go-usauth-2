package handlers

import (
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/asaskevich/govalidator"
	"github.com/julienschmidt/httprouter"
)

// LoginHandlerGet is userd asd
func LoginHandlerGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()	
	urlNext := ""
	if(len(r.Form["url_next"]) > 0){
		urlNext = r.Form["url_next"][0]
	}
	loginTemplate.Execute(w, map[string]interface{}{
		"url_next": urlNext,
    })
}

// LoginHandlerPost us
func LoginHandlerPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	
	errParseForm := r.ParseForm()
	
	urlNext := ""
	if(len(r.Form["url_next"]) > 0){
		urlNext = r.Form["url_next"][0]
	}

	if errParseForm != nil {
		fmt.Println(errParseForm)
		loginTemplate.Execute(w, map[string]interface{}{
			"error": constErrInternalError,
			"url_next": urlNext,
		})
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || !govalidator.IsEmail(email) {
		loginTemplate.Execute(w, map[string]interface{}{
			"error": constErrEmailMissing,
			"url_next": urlNext,
		})
		return
	}

	if password == "" {
		loginTemplate.Execute(w, map[string]interface{}{
			"error": constErrPasswordMissing,
			"url_next": urlNext,
		})
		//loginTemplate.Execute(w, constErrPasswordMissing)
		return
	}

	present, userID, passwordHash, errGetPassword := userService.GetPasswordHash(email)

	if errGetPassword!=nil{
		loginTemplate.Execute(w, map[string]interface{}{
			"error": constErrInternalError,
			"url_next": urlNext,
		})
		return
	}

	if !present {
		//fmt.Println("User not there")
		//loginTemplate.Execute(w, constErrNotRegistered)
		loginTemplate.Execute(w, map[string]interface{}{
			"error": constErrNotRegistered,
			"url_next": urlNext,
		})
		return
	}

	verified := checkPassword(&password, &passwordHash)

	if !verified {
		fmt.Println("Wrong passwword")
		loginTemplate.Execute(w, map[string]interface{}{
			"error": constErrPasswordMatchFailed,
			"url_next": urlNext,
		})
		//loginTemplate.Execute(w, )
		return
	}

	fmt.Println(userID)
	
	newSessionID, errGen := generateRandomString(48)
	if errGen != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	plainValue := map[string]string{"sessionid": newSessionID}
	encodedValue, errCookieEncode := cookie.Encode("sessionid", plainValue)
	
	if errCookieEncode != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else{
		newCookie := &http.Cookie{
			Name:  "sessionid",
			Value: encodedValue,
			Path:  "/",
			HttpOnly : true,
		} 
		http.SetCookie(w, newCookie)
		
		/*
		z := cache.SetSessionValue(newSessionId,"userId",userID)
		if z!=nil{
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		*/
	
		http.Redirect(w, r, urlNext, http.StatusSeeOther)

	}
	
	loginTemplate.Execute(w, nil)
}

func checkPassword(userPassword, savedPassword *string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(*savedPassword), []byte(*userPassword))
	if err != nil {
		return false
	}
	return true
}