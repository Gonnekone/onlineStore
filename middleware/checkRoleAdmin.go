package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/Gonnekone/onlineStore/models"
)

func CheckRoleAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user")

		if user != nil {
			if u, ok := user.(models.User); ok {
				if u.Role == "ADMIN" {
					next.ServeHTTP(w, r)
					return
				} else {
					response := map[string]string{"error": "Ты не админ"}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(response)
					return
				}
			} else {
				response := map[string]string{"error": "Как нахуй так, что ты вытащил по ключу user не user!?"}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(response)
				return
			}
		} else {
			response := map[string]string{"error": "user = nil"}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)
			return
		}
	})
}
