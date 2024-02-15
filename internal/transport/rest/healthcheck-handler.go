package rest

import (
	"FriendsAdvice/internal/transport"
	"net/http"
)

func Liveness(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func Readiness(controller transport.IController) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if dbReady := controller.IsStorageReady(); !dbReady {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
