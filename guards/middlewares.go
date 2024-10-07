package guards

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Tutuacs/internal/auth"
	"github.com/dgrijalva/jwt-go"

	"github.com/Tutuacs/pkg/logs"
	"github.com/Tutuacs/pkg/resolver"
)

type contextKey string

const UserKey contextKey = "user"

func AutenticatedRoute(handlerFunc http.HandlerFunc, roles ...int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tokenString := resolver.GetTokenFromRequest(r)

		token, err := auth.ValidateJWT(tokenString)
		if err != nil {
			logs.ErrorLog(fmt.Sprintf("failed to validate token: %v", err))
			permissionDenied(w)
			return
		}

		if !token.Valid {
			log.Println("invalid token")
			permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		str := claims["userID"].(string)

		userID, err := strconv.Atoi(str)
		if err != nil {
			logs.ErrorLog(fmt.Sprintf("failed to convert userID to int: %v", err))
			permissionDenied(w)
			return
		}

		store, err := auth.NewStore()
		if err != nil {
			resolver.WriteResponse(w, http.StatusInternalServerError, map[string]string{"Error": err.Error()})
			return
		}

		u, err := store.GetUserByID(userID)
		if err != nil {
			logs.ErrorLog(fmt.Sprintf("failed to get user by id: %v", err))
			permissionDenied(w)
			return
		}

		// Verificar roles se roles != vazio
		if len(roles) > 0 {
			isAuthorized := false
			for _, role := range roles {
				if u.Role == role { // Supondo que o usuário tenha o campo Role no seu struct
					isAuthorized = true
					break
				}
			}
			if !isAuthorized {
				resolver.WriteResponse(w, http.StatusUnauthorized, map[string]string{"Error": "unauthorized"})
				return
			}
		}

		// Adiciona o usuário ao contexto
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, u)
		r = r.WithContext(ctx)

		// Executa o handler se o token for válido
		handlerFunc(w, r)
	}
}

func permissionDenied(w http.ResponseWriter) {
	resolver.WriteResponse(w, http.StatusForbidden, fmt.Sprintf("permission denied"))
}
