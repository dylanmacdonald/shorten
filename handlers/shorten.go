package handlers

import (
	"net/http"

	"github.com/dylanmacdonald/shorten/api"
	"github.com/dylanmacdonald/shorten/service"
	"github.com/gorilla/schema"

	"github.com/Sirupsen/logrus"
)

func Shorten(logger logrus.FieldLogger, s service.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, err := decodeShortenRequest(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logger.WithError(err).Error()
			return
		}

		url, err := s.Shorten(r.Context(), req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.WithError(err).Error()
			return
		}

		_, err = w.Write([]byte(url.String()))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.WithError(err).Error()
			return
		}
	})
}

func decodeShortenRequest(r *http.Request) (*api.ShortenRequest, error) {
	req := &api.ShortenRequest{}
	return req, schema.NewDecoder().Decode(req, r.URL.Query())
}
