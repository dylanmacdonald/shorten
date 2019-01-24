package handlers

import (
	"net/http"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/dylanmacdonald/shorten/service"
)

func Redirect(logger logrus.FieldLogger, s service.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path, err := strconv.Atoi(r.URL.EscapedPath()[1:])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logger.WithError(err).Error()
			return
		}
		uri, err := s.Redirect(r.Context(), int64(path))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			logger.WithError(err).Error()
			return
		}

		http.Redirect(w, r, uri.String(), http.StatusTemporaryRedirect)
	})
}
