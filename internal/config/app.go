package config

import (
	"strconv"
	"strings"
)

const (
	portMinimo         = 1                            // Porta mínima válida
	portMaximo         = 65535                        // Porta máxima válida
	portValidoMensagem = "deve estar entre 1 e 65535" // Mensagem de validação de porta
	appNameDefault     = "GestorDeChamadosAPI"        // Nome da aplicação
	versionDefault     = "1.0.0"                      // Versão da aplicação
	logLevelDefault    = "INFO"                       // Nível de log padrão
	PortDefault        = "8080"                       // Porta padrão da aplicação
	environmentDefault = "development"                // Ambiente padrão da aplicação
)

// envsValidos contém os ambientes válidos
var envsValidos = map[string]struct{}{
	"local":       {},
	"development": {},
	"production":  {},
}

// logSValidos contém os níveis de log válidos
var logSValidos = map[string]struct{}{
	"DEBUG": {},
	"INFO":  {},
	"WARN":  {},
	"ERROR": {},
}

// AppConfig contém configurações gerais da aplicação
type AppConfig struct {
	name        string // Nome da aplicação
	version     string // Versão da aplicação
	environment string // Ambiente da aplicação
	logLevel    string // Nível de log
	port        string // Porta da aplicação
	corsOrigin  string // Origem permitida para CORS
}

// carregarAppConfig carrega as configurações da aplicação
func carregarAppConfig() AppConfig {
	return AppConfig{
		name:        getEnv("APP_NAME", appNameDefault),
		version:     getEnv("VERSION", versionDefault),
		environment: getEnv("ENVIRONMENT", environmentDefault),
		logLevel:    getEnv("LOG_LEVEL", logLevelDefault),
		port:        getEnv("PORT", PortDefault),
		corsOrigin:  getEnv("CORS_ORIGIN", stringVazia),
	}
}

// ValidarApp valida as configurações da aplicação.
//
// Em caso de erros de validação, retorna um erro do tipo ErrosConfig.
func (a *AppConfig) Validar() error {
	erros := NovoErrosConfig()

	// -- Environment --
	if _, ok := envsValidos[a.environment]; !ok {
		erros.Add("ENVIRONMENT", "inválido, os válidos são: "+EnvironmentsValidos())
	}

	// -- Port --
	port, err := strconv.Atoi(a.port)
	if err != nil {
		erros.Add("PORT", "inválido: deve ser um número inteiro")
	}
	
	if port < portMinimo || port > portMaximo {
		erros.Add("PORT", portValidoMensagem)
	}

	// CORS
	if a.environment == production && a.corsOrigin == stringVazia {
		erros.Add("CORS_ORIGIN", "é obrigatório em produção")
	}

	// -- Log Level --
	if _, ok := logSValidos[a.logLevel]; !ok {
		erros.Add("LOG_LEVEL", "inválido, os válidos são: "+LogsValidos())
	}

	if erros.HaErros() {
		return erros
	}

	return nil
}

// EnvironmentsValidos retorna todos os ambientes válidos formatados.
func EnvironmentsValidos() string {
	envs := make([]string, 0, len(envsValidos))
	for e := range envsValidos {
		envs = append(envs, e)
	}
	return strings.Join(envs, ", ")
}

// LogsValidos retorna todos os níveis de log válidos formatados.
func LogsValidos() string {
	levels := make([]string, 0, len(logSValidos))
	for l := range logSValidos {
		levels = append(levels, l)
	}
	return strings.Join(levels, ", ")
}

// Métodos de acesso

func (c AppConfig) Name() string        { return c.name }
func (c AppConfig) Version() string     { return c.version }
func (c AppConfig) Environment() string { return c.environment }
func (c AppConfig) LogLevel() string    { return c.logLevel }
func (c AppConfig) Port() string        { return c.port }
func (c AppConfig) CORSOrigin() string  { return c.corsOrigin }
