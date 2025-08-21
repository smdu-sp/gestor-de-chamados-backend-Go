package auth

import "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/model"

// IsAllowed verifica se a permissão do usuário (userPerm)
// está dentro da lista de permissões permitidas (allowed...).
func IsAllowed(userPerm string, allowed ...model.Permissao) bool {
	for _, p := range allowed {
		if string(p) == userPerm {
			return true
		}
	}
	return false
}
