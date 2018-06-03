package handlers

import (
	"html/template"
)

var loginTemplate *template.Template
var signupTemplate *template.Template

//var resetTemplate *template.Template

func parseTemplates() {
	loginTemplate = template.Must(template.ParseFiles("./template/login.html"))
	signupTemplate = template.Must(template.ParseFiles("./template/signup.html"))
}
