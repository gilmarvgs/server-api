package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gilmarvgs/api"
	"github.com/gilmarvgs/db"
)

type Cotacao struct {
	USDBRL struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

func Handler(database *sql.DB, w http.ResponseWriter, r *http.Request) {
	log.Println(">>>>>>>>>>>>>request iniciada")

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 150*time.Millisecond)
	defer cancel()

	defer log.Println("request finalizada<<<<<<<<<<<<<")
	cotacao, err := api.GetCotacao(ctx)

	if err != nil {
		log.Println("erro na chamada da api de Cotacao:", err)
		http.Error(w, fmt.Sprintf("erro na chamada da api de Cotacao: %v", err), http.StatusInternalServerError)
		return
	}

	select {
	case <-time.After(1 * time.Millisecond):

		bid, err := processaCotacao(database, cotacao)

		if err != nil {
			log.Println("erro ao processar cotação:", err)
			http.Error(w, "erro ao processar cotação", http.StatusInternalServerError)
			return
		} else {

			response := map[string]float64{"bid": bid}
			responseJSON, err := json.Marshal(response)
			if err != nil {
				log.Println("erro ao criar resposta JSON:", err)
				http.Error(w, "Erro ao criar resposta JSON", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(responseJSON)
			log.Println("cotação processada com sucesso! Bid:", bid)
			return
		}
	case <-ctx.Done():
		log.Println("request cancelada por timeout", ctx.Err())
		http.Error(w, "Tempo de processamento excedido, tente novamente", http.StatusRequestTimeout)
		return
	}
}

func processaCotacao(database *sql.DB, cotacao string) (float64, error) {

	var cotacaoData Cotacao
	err := json.Unmarshal([]byte(cotacao), &cotacaoData)

	if err != nil {
		return 0, fmt.Errorf("erro ao fazer o unmarshal: %w", err)
	}

	bid, _ := strconv.ParseFloat(cotacaoData.USDBRL.Bid, 64)

	currentTime := time.Now().Format("2006-01-02 15:04:05.999999999")

	ctxInsert := context.Background()
	ctxInsert, cancelInsert := context.WithTimeout(ctxInsert, 10*time.Millisecond)
	defer cancelInsert()

	start := time.Now()
	err = db.InsertCotacao(ctxInsert, database, currentTime, bid)
	elapsed := time.Since(start)

	log.Printf("Tempo de execução da inserção: %v\n", elapsed)

	if err != nil {
		return 0, fmt.Errorf("erro ao inserir registro no banco de dados: %w", err)
	}

	if ctxInsert.Err() != nil {
		log.Println("erro no contexto de inserção:", ctxInsert.Err())
		return 0, fmt.Errorf("erro no contexto de inserção: %w", ctxInsert.Err())
	}

	return bid, nil
}
