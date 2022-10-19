package middlewares

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/bubo-py/McK/users/service"
)

func Authenticate(bl service.BusinessLogic) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			login, pwd, ok := r.BasicAuth()
			if r.URL.Path != "/api/users" && r.Method != "POST" {
				if !ok {
					err := json.NewEncoder(w).Encode("please provide your credentials")
					if err != nil {
						log.Println(err)
					}
					return
				}

				err := bl.LoginUser(r.Context(), login, pwd)
				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					err = json.NewEncoder(w).Encode(err.Error())
					if err != nil {
						log.Println(err)
					}
					return
				}
			}

			user, _ := bl.Db.GetUserByLogin(r.Context(), login)

			ctxWithUser := context.WithValue(r.Context(), "login", user.Login)
			ctxWithUser = context.WithValue(r.Context(), "timezone", user.Timezone)

			rWithUser := r.WithContext(ctxWithUser)
			next.ServeHTTP(w, rWithUser)
		})
	}
}
