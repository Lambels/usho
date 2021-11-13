package handlers

import (
	"context"
	"net/http"

	"github.com/Lambels/usho/encoding/json"
	"github.com/Lambels/usho/repo"
)

func MiddlewareValidateIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		in := repo.URLRequest{}

		err := json.Decode(&in, r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.Encode(map[string]string{"error": err.Error()}, w)
			return
		}

		ctx := context.WithValue(r.Context(), repo.URLKey{}, in)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
