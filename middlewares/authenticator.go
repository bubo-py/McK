package middlewares

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bubo-py/McK/contextHelpers"
	"github.com/bubo-py/McK/users/service"
)

func Authenticate(bl service.BusinessLogicInterface) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			login, pwd, ok := r.BasicAuth()
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

			user, err := bl.GetUserByLogin(r.Context(), login)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				err = json.NewEncoder(w).Encode(err.Error())
				if err != nil {
					log.Println(err)
				}
				return
			}

			ctxWithUser := contextHelpers.WriteLoginToContext(r.Context(), user.Login)
			ctxWithUser = contextHelpers.WriteTimezoneToContext(r.Context(), user.Timezone)

			rWithUser := r.WithContext(ctxWithUser)
			next.ServeHTTP(w, rWithUser)
		})
	}
}
