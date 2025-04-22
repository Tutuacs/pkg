package guards

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"

	"github.com/Tutuacs/pkg/enums"
	JWT "github.com/Tutuacs/pkg/jwt"
	"github.com/Tutuacs/pkg/logs"
	"github.com/Tutuacs/pkg/resolver"
)

type contextKey string

const UserKey contextKey = "user"

type Guards struct {
	store *Store
}

func UseGuard(conn *sql.DB) *Guards {
	return &Guards{store: &Store{db: conn}}
}

func (g *Guards) AutenticatedRoute(handlerFunc http.HandlerFunc, roles ...enums.Role) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tokenString := resolver.GetTokenFromRequest(r)

		token, err := JWT.ValidateJWT(tokenString)
		if err != nil {
			logs.ErrorLog(fmt.Sprintf("failed to validate token: %v", err))
			g.permissionDenied(w)
			return
		}

		if !token.Valid {
			log.Println("invalid token")
			g.permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		str := claims["userID"].(string)

		userID, err := strconv.Atoi(str)
		if err != nil {
			logs.ErrorLog(fmt.Sprintf("failed to convert userID to int: %v", err))
			g.permissionDenied(w)
			return
		}

		str = fmt.Sprintf("%v", claims["role"])

		tokenRole, err := strconv.ParseInt(str, 10, 8)
		if err != nil {
			logs.ErrorLog(fmt.Sprintf("failed to convert role to int: %v", err))
			g.permissionDenied(w)
			return
		}

		store := g.store

		u, err := store.GetUserByID(userID)
		if err != nil {
			logs.ErrorLog(fmt.Sprintf("failed to get user by id: %v", err))
			g.permissionDenied(w)
			return
		}

		// Verificar roles se roles != vazio
		if len(roles) > 0 {
			notAuthorized := true
			for _, role := range roles {
				if tokenRole == int64(role) { // Supondo que o usuário tenha o campo Role no seu struct
					notAuthorized = false
					break
				}
			}
			if notAuthorized {
				resolver.WriteResponse(w, http.StatusUnauthorized, map[string]string{"Error": "unauthorized"})
				return
			}

			if tokenRole != int64(u.Role) {
				resolver.WriteResponse(w, http.StatusUnauthorized, map[string]string{"Error": "role mismatch, make login again"})
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

func (g *Guards) AuthenticatedUrlRoute(handlerFunc http.HandlerFunc, roles ...enums.Role) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tokenString := resolver.GetQueryParam(r, "token")

		token, err := JWT.ValidateJWT(tokenString)
		if err != nil {
			logs.ErrorLog(fmt.Sprintf("failed to validate token: %v", err))
			g.permissionDenied(w)
			return
		}

		if !token.Valid {
			log.Println("invalid token")
			g.permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		str := claims["userID"].(string)

		userID, err := strconv.Atoi(str)
		if err != nil {
			logs.ErrorLog(fmt.Sprintf("failed to convert userID to int: %v", err))
			g.permissionDenied(w)
			return
		}

		str = fmt.Sprintf("%v", claims["role"])

		tokenRole, err := strconv.ParseInt(str, 10, 8)
		if err != nil {
			logs.ErrorLog(fmt.Sprintf("failed to convert role to int: %v", err))
			g.permissionDenied(w)
			return
		}

		store := g.store

		u, err := store.GetUserByID(userID)
		if err != nil {
			logs.ErrorLog(fmt.Sprintf("failed to get user by id: %v", err))
			g.permissionDenied(w)
			return
		}

		// Verificar roles se roles != vazio
		if len(roles) > 0 {
			notAuthorized := true
			for _, role := range roles {
				if tokenRole == int64(role) { // Supondo que o usuário tenha o campo Role no seu struct
					notAuthorized = false
					break
				}
			}
			if notAuthorized {
				resolver.WriteResponse(w, http.StatusUnauthorized, map[string]string{"Error": "unauthorized"})
				return
			}

			if tokenRole != int64(u.Role) {
				resolver.WriteResponse(w, http.StatusUnauthorized, map[string]string{"Error": "role mismatch, make login again"})
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

func (g *Guards) permissionDenied(w http.ResponseWriter) {
	resolver.WriteResponse(w, http.StatusForbidden, fmt.Sprintf("permission denied"))
}
