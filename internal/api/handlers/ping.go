package handlers

import "net/http"

func PingHandler(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Если запись не удалась - возвращаем ошибку
	if _, err := w.Write([]byte(`{"status":"ok", "handler": "ping"}`)); err != nil {
		return err
	}

	return nil
}
