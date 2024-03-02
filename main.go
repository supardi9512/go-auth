package main

import (
	"fmt"
	authcontroller "go-auth/controllers"
	"net/http"
)

func main() {
	http.HandleFunc("/", authcontroller.Index)
	http.HandleFunc("/login", authcontroller.Login)

	fmt.Println("Server is running on : http://localhost:3000")
	http.ListenAndServe(":3000", nil)
}
