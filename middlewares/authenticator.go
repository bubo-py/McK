package middlewares

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bubo-py/McK/contextHelpers"
	"github.com/bubo-py/McK/customErrors"
	"github.com/bubo-py/McK/users/service"
)

var unauthenticatedReturn = customErrors.ReturnError{
	ErrorType:    customErrors.ErrUnauthenticated.ErrorType,
	ErrorMessage: customErrors.ErrUnauthenticated.Error(),
}

func Authenticate(bl service.BusinessLogicInterface) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			login, pwd, ok := r.BasicAuth()
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				err := json.NewEncoder(w).Encode(unauthenticatedReturn)
				if err != nil {
					log.Println(err)
				}
				return
			}

			err := bl.LoginUser(r.Context(), login, pwd)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				err = json.NewEncoder(w).Encode(unauthenticatedReturn)
				if err != nil {
					log.Println(err)
				}
				return
			}

			user, err := bl.GetUserByLogin(r.Context(), login)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				err = json.NewEncoder(w).Encode(unauthenticatedReturn)
				if err != nil {
					log.Println(err)
				}
				return
			}

			r = r.WithContext(contextHelpers.WriteLoginToContext(r.Context(), user.Login))
			r = r.WithContext(contextHelpers.WriteTimezoneToContext(r.Context(), user.Timezone))

			next.ServeHTTP(w, r)
		})
	}
}
