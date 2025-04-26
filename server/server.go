package main

import (
	"log"
	"net/http"

	"github.com/gilmarvgs/db"
	"github.com/gilmarvgs/handler"
)

func main() {

	database, err := db.InitializeDatabase()

	if err != nil {
		log.Fatalf("Erro ao inicializar o banco de dados: %v", err)
	}

	defer database.Close()

	// Configura o endpoint e passa o banco de dados para o handler
	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		handler.Handler(database, w, r)
	})

	http.ListenAndServe(":8080", nil)
}
