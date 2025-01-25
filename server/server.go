package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("request iniciada")
	defer log.Println("request finalizada")
	select {
	case <-time.After(10 * time.Second):
		// imprime o status da requisição no console
		log.Println("request processada com sucesso")
		// imprime a resposta no navegador
		w.Write([]byte("request processada com sucesso"))
	case <-ctx.Done():
		log.Println("request cancelada pelo cliente", ctx.Err())
		http.Error(w, "request cancelada pelo cliente", http.StatusRequestTimeout)
	}
}
