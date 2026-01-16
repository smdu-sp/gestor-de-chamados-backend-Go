package config

import "time"

const (
	serverReadTimeoutDefault         = "15s" // Tempo máximo padrão para ler a requisição
	serverWriteTimeoutDefault        = "15s" // Tempo máximo padrão para escrever a resposta
	serverIdleTimeoutDefault         = "60s" // Tempo máximo padrão para conexões ociosas
	serverShutdownTimeoutDefault     = "10s" // Tempo máximo padrão para shutdown gracioso
	serverSignalChannelBufferDefault = 1     // Buffer padrão do canal de sinais
	serverErrorChannelBufferDefault  = 1     // Buffer padrão do canal de erros do servidor
)

type ServerConfig struct {
	readTimeout             time.Duration // Tempo máximo para ler a requisição
	writeTimeout            time.Duration // Tempo máximo para escrever a resposta
	idleTimeout             time.Duration // Tempo máximo para conexões ociosas
	shutdownTimeout         time.Duration // Tempo máximo para shutdown gracioso
	signalChannelBufferSize int           // Buffer do canal de sinais
	errorChannelBufferSize  int           // Buffer do canal de erros do servidor
}

// carregarServerConfig carrega as configurações do servidor
func carregarServerConfig() ServerConfig {
	return ServerConfig{
		readTimeout:             getEnvDuration("SERVER_READ_TIMEOUT", serverReadTimeoutDefault),
		writeTimeout:            getEnvDuration("SERVER_WRITE_TIMEOUT", serverWriteTimeoutDefault),
		idleTimeout:             getEnvDuration("SERVER_IDLE_TIMEOUT", serverIdleTimeoutDefault),
		shutdownTimeout:         getEnvDuration("SERVER_SHUTDOWN_TIMEOUT", serverShutdownTimeoutDefault),
		signalChannelBufferSize: getEnvInt("SERVER_SIGNAL_CHANNEL_BUFFER_SIZE", serverSignalChannelBufferDefault),
		errorChannelBufferSize:  getEnvInt("SERVER_ERROR_CHANNEL_BUFFER_SIZE", serverErrorChannelBufferDefault),
	}
}

// Validar valida as configurações do servidor.
//
// Em caso de erros de validação, retorna um erro do tipo ErrosConfig.
func (s *ServerConfig) Validar() error {
	erros := NovoErrosConfig()
	errMaiorQueZero := "deve ser maior que 0"
	errMaiorIgualUm := "deve ser no mínimo 1"

	// Timeouts
	if s.readTimeout <= 0 {
		erros.Add("SERVER_READ_TIMEOUT", errMaiorQueZero)
	}

	if s.writeTimeout <= 0 {
		erros.Add("SERVER_WRITE_TIMEOUT", errMaiorQueZero)
	}

	if s.idleTimeout <= 0 {
		erros.Add("SERVER_IDLE_TIMEOUT", errMaiorQueZero)
	}

	if s.shutdownTimeout <= 0 {
		erros.Add("SERVER_SHUTDOWN_TIMEOUT", errMaiorQueZero)
	}

	if s.shutdownTimeout > s.readTimeout+s.writeTimeout {
		erros.Add(
			"SERVER_SHUTDOWN_TIMEOUT",
			"não deve ser maior que a soma de READ + WRITE timeout",
		)
	}

	// Buffers de canais
	if s.signalChannelBufferSize < 1 {
		erros.Add(
			"SERVER_SIGNAL_CHANNEL_BUFFER_SIZE",
			errMaiorIgualUm,
		)
	}

	if s.errorChannelBufferSize < 1 {
		erros.Add(
			"SERVER_ERROR_CHANNEL_BUFFER_SIZE",
			errMaiorIgualUm,
		)
	}

	if erros.HaErros() {
		return erros
	}

	return nil
}

// Métodos de acesso

func (s ServerConfig) ReadTimeout() time.Duration     { return s.readTimeout }
func (s ServerConfig) WriteTimeout() time.Duration    { return s.writeTimeout }
func (s ServerConfig) IdleTimeout() time.Duration     { return s.idleTimeout }
func (s ServerConfig) ShutdownTimeout() time.Duration { return s.shutdownTimeout }
func (s ServerConfig) SignalChannelBufferSize() int   { return s.signalChannelBufferSize }
func (s ServerConfig) ErrorChannelBufferSize() int    { return s.errorChannelBufferSize }
