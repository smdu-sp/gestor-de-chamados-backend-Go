package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/smdu-sp/gestor-de-chamados-backend-Go/docs" // Swagger docs
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/config"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/infra/db"
	_ "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/interface/handler"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/interface/router"
)

// @title Gestor de Chamados API
// @version 1.0
// @description API para gestão de chamados
// @host localhost:8080
// @BasePath /
func main() {
	if err := run(); err != nil {
		log.Fatalf("[main] erro ao iniciar a aplicação: %v", err)
	}
}

func run() error {
	// Carrega configuração
	cfg := config.Load()

	// Conecta ao banco de dados passando a configuração
	dbConn, err := db.ConectarMySQL(cfg)
	if err != nil {
		return fmt.Errorf("[main.run]: %w", err)
	}
	// Garante que a conexão será fechada ao final
	defer dbConn.Close()

	// Monta o router
	r := router.InicializarRoteadorHTTP(cfg, dbConn)

	// Cria o servidor HTTP
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Inicia o servidor em goroutine
	go func() {
		log.Printf("API rodando em http://localhost:%s", cfg.Port)
		log.Printf("Swagger em http://localhost:%s/swagger/index.html", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[main] erro no servidor: %v", err)
		}
	}()

	// Aguarda sinal de interrupção para desligar o servidor graciosamente
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("[main] Desligando o servidor...")

	// Cria um contexto com timeout para o desligamento
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return srv.Shutdown(ctx)
}
