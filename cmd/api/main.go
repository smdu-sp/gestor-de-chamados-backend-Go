package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/db"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/config"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/http/router"
	_ "github.com/smdu-sp/gestor-de-chamados-backend-Go/docs"
)

// @title Gestor de Chamados API
// @version 1.0
// @description API para gestão de chamados
// @host localhost:8080
// @BasePath /
func main() {
	if err := run(); err != nil {
		log.Fatalf("erro ao iniciar a aplicação: %v", err)
	}
}

func run() error {
	// Carrega configuração
	cfg := config.Load()

	// Conecta ao banco de dados passando a configuração
	dbConn, err := db.OpenMySQL(cfg)
	if err != nil {
		return err
	}
	// Garante que a conexão será fechada ao final
	defer dbConn.Close()

	// Monta o router
	r := router.Build(cfg, dbConn)

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
			log.Fatalf("erro no servidor: %v", err)
		}
	}()

	// Aguarda sinal de interrupção para desligar o servidor graciosamente
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Desligando o servidor...")

	// Cria um contexto com timeout para o desligamento
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return srv.Shutdown(ctx)
}
