package auth

import "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/model"

func IsAllowed(userPerm string, allowed ...model.Permissao) bool {
	for _, p := range allowed {
		if string(p) == userPerm {
			return true
		}
	}
	return false
}
