package handlers

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/dylanmacdonald/shorten/service"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func InitHandlers(logger logrus.FieldLogger, s service.Service) *mux.Router {
	r := mux.NewRouter()

	shortener := Shorten(logger.WithField("handler", "shortener"), s)
	redirect := Redirect(logger.WithField("handler", "redirect"), s)

	r.Handle("/shorten", handlers.MethodHandler{
		"GET": shortener,
	})
	r.Handle("/health", handlers.MethodHandler{
		"GET": http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	})
	r.PathPrefix("/").Handler(redirect)

	return r
}
