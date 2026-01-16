package banner

import (
	_ "embed"
	"strings"
)

//go:embed banner.txt
var bannerContent string

// CarregarBanner retorna o banner da aplicação com a versão fornecida
func CarregarBanner(versao string) string {
	banner := strings.ReplaceAll(bannerContent, "{{VERSAO}}", versao)
	bannerColorido := "\033[1;36m" + banner + "\033[0m"
	return bannerColorido
}
