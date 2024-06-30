package main

import (
	"belajar/controller/auth"
	"belajar/controller/pakaian"
	"belajar/controller/pembeli"
	"belajar/database"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	database.InitDB()
	fmt.Println("Hello World")

	router := mux.NewRouter()

	router.HandleFunc("/regis", auth.Registration).Methods("POST")
	router.HandleFunc("/login", auth.Login).Methods("POST")

	// Router handler pakaian
	router.HandleFunc("/pakaian", pakaian.GetPakaian).Methods("GET")
	router.HandleFunc("/pakaian", auth.JWTAuth(pakaian.PostPakaian)).Methods("POST")
	router.HandleFunc("/pakaian/{id}", auth.JWTAuth(pakaian.PutPakaian)).Methods("PUT")
	router.HandleFunc("/pakaian/{id}", auth.JWTAuth(pakaian.DeletePakaian)).Methods("DELETE")

	// Router handler pembeli
	router.HandleFunc("/pembeli", pembeli.GetPembeli).Methods("GET")
	router.HandleFunc("/pembeli", auth.JWTAuth(pembeli.PostPembeli)).Methods("POST")
	router.HandleFunc("/pembeli/{id}", auth.JWTAuth(pembeli.PutPembeli)).Methods("PUT")
	router.HandleFunc("/pembeli/{id}", auth.JWTAuth(pembeli.DeletePembeli)).Methods("DELETE")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://127.0.0.1:5500", "http://localhost:8000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		Debug:            true,
	})

	handler := c.Handler(router)

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
