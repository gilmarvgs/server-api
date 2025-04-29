package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitializeDatabase() (*sql.DB, error) {

	// Conecta ao banco de dados SQLite (ou cria o arquivo se não existir)
	db, err := sql.Open("sqlite3", "./one_database.db")
	if err != nil {
		return nil, err
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS cotacao (
		data TIMESTAMP NOT NULL,
		valor DECIMAL(10, 2) NOT NULL
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	log.Println("banco de dados SQLite configurado com sucesso!")
	return db, nil
}

func InsertCotacao(ctx context.Context, db *sql.DB, data string, valor float64) error {
	// Configura o log para incluir milissegundos
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	log.Println("inserindo cotação no banco de dados...")
	// Prepara a consulta para inserção
	insertSQL := `INSERT INTO cotacao (data, valor) VALUES (?, ?)`
	stmt, err := db.PrepareContext(ctx, insertSQL)
	if err != nil {
		return fmt.Errorf("erro ao preparar a consulta: %w", err)
	}
	defer stmt.Close()

	// Executa a inserção com o contexto
	_, err = stmt.ExecContext(ctx, data, valor)
	log.Println("fim inserção")
	if err != nil {
		// Verifica se o erro foi causado por timeout
		if ctx.Err() != nil {
			if ctx.Err() == context.DeadlineExceeded {
				return fmt.Errorf("operação excedeu o limite de tempo")
			}
		}
		return fmt.Errorf("erro ao executar a inserção: %w", err)
	}

	return nil
}
