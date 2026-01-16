package config

import (
	"os"
	"strconv"
	"time"
)

const (
	numMinCharsTokenSecret         = 32               // naumero mínimo de caracteres para o token secret
	numMinCharsRefreshTokenSecret  = 32               // número mínimo de caracteres para o refresh token secret
	tokenTTLmin                    = 5 * time.Minute  // número mínimo de duração para o token TTL (5 minutos)
	refreshTokenTTLmin             = 30 * time.Minute // número mínimo de duração para o refresh token TTL (30 minutos)
	tokenTTLMax                    = 48 * time.Hour   // número máximo de duração para o token TTL (2 dias)
	refreshTokenTTLMax             = 720 * time.Hour  // número máximo de duração para o refresh token TTL (30 dias)
	DurationTokenTTLDefault        = "48h"            // valor padrão para duração (2 dias)
	DurationRefreshTokenTTLDefault = "168h"           // valor padrão para duração (7 dias)
)

// AuthConfig contém configurações de autenticação
type AuthConfig struct {
	tokenSecret        string // segredo para assinatura do token
	refreshTokenSecret string // segredo para assinatura do refresh token
	tokenTTL           string // tempo de vida do token
	refreshTokenTTL    string // tempo de vida do refresh token
}

// carregarAuthConfig carrega as configurações de autenticação
func carregarAuthConfig() AuthConfig {
	return AuthConfig{
		tokenSecret:        os.Getenv("TOKEN_SECRET"),
		refreshTokenSecret: os.Getenv("REFRESH_TOKEN_SECRET"),
		tokenTTL:           getEnv("TOKEN_TTL", DurationTokenTTLDefault),
		refreshTokenTTL:    getEnv("REFRESH_TOKEN_TTL", DurationRefreshTokenTTLDefault),
	}
}

// ValidarAuth valida as configurações de autenticação.
//
// Em caso de erros, retorna um erro do tipo ErrosConfig.
func (c *AuthConfig) Validar() error {
	erros := NovoErrosConfig()

	// -- token Secret --
	if c.tokenSecret == "" {
		erros.Add("TOKEN_SECRET", "é obrigatório")
	}

	if len(c.tokenSecret) < numMinCharsTokenSecret {
		erros.Add("TOKEN_SECRET",
			"deve ter no mínimo "+strconv.Itoa(numMinCharsTokenSecret)+" caracteres",
		)
	}

	// -- refresh Token Secret --
	if c.refreshTokenSecret == "" {
		erros.Add("REFRESH_TOKEN_SECRET", "é obrigatório")
	}

	if len(c.refreshTokenSecret) < numMinCharsRefreshTokenSecret {
		erros.Add("REFRESH_TOKEN_SECRET",
			"deve ter no mínimo "+strconv.Itoa(numMinCharsRefreshTokenSecret)+" caracteres",
		)
	}

	// -- token TTL --
	tokenTTL, err := time.ParseDuration(c.tokenTTL)
	if err != nil || tokenTTL < tokenTTLmin {
		erros.Add("TOKEN_TTL",
			"inválido: deve ter uma duração maior ou igual a "+tokenTTLmin.String(),
		)
	}

	// -- refresh Token TTL --
	refreshTokenTTL, err := time.ParseDuration(c.refreshTokenTTL)
	if err != nil || refreshTokenTTL < refreshTokenTTLmin {
		erros.Add("REFRESH_TOKEN_TTL",
			"inválido: deve ter uma duração maior ou igual a "+refreshTokenTTLmin.String(),
		)
	}

	// -- Limites máximos --
	if tokenTTL > tokenTTLMax {
		erros.Add("TOKEN_TTL",
			"muito alto: máximo recomendado é "+tokenTTLMax.String(),
		)
	}

	if refreshTokenTTL > refreshTokenTTLMax {
		erros.Add("REFRESH_TOKEN_TTL",
			"muito alto: máximo recomendado é "+refreshTokenTTLMax.String(),
		)
	}

	if erros.HaErros() {
		return erros
	}

	return nil
}

// Métodos de acesso

func (a AuthConfig) TokenSecret() string        { return a.tokenSecret }
func (a AuthConfig) RefreshTokenSecret() string { return a.refreshTokenSecret }
func (a AuthConfig) TokenTTL() string           { return a.tokenTTL }
func (a AuthConfig) RefreshTokenTTL() string    { return a.refreshTokenTTL }
