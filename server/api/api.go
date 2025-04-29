package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func GetCotacao(ctx context.Context) (string, error) {

	// Cria a requisição HTTP com o contexto
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return "", fmt.Errorf("erro ao criar a requisição: %w", err)
	}

	// Executa a requisição
	resp, err := http.DefaultClient.Do(req)
	if err != nil {

		if ctx.Err() != nil {
			return "", fmt.Errorf("requisição expirada: %w", ctx.Err())
		}
		return "", fmt.Errorf("erro ao executar a requisição: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("erro ao ler o corpo da resposta: %w", err)
	}

	return string(body), nil
}
