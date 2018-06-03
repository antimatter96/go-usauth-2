package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	
	"./db"
	"./constants"
	"./handlers"
	
	"github.com/julienschmidt/httprouter"
	gorillaHandlers "github.com/gorilla/handlers"
)

func main(){
	
	config := flag.String("config", "config", "config file")
	flag.Parse()

	if err := constants.Init(*config); err != nil{
		fmt.Printf("cant initialize constants : %v", err)
	}

	db.Init()
	handlers.Init()

	router := httprouter.New()

	router.GET("/login", handlers.LoginHandlerGet)
	router.POST("/login", handlers.LoginHandlerPost)

	router.ServeFiles("/static/*filepath", http.Dir("./template/static/"))
	
	output,_ := constants.Value("output").(string)
	
	file, err := os.Create(output)
	if err != nil {
		fmt.Printf("could not create file %s : %v", output, err)
	}

	loggedRouter := gorillaHandlers.LoggingHandler(file, router)

	port,_ := constants.Value("port").(string)

	if err := http.ListenAndServe(port, loggedRouter); err!=nil{
		fmt.Printf("error starting server: %v", err)
	}
}