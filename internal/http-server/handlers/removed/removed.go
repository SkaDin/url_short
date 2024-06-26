package removed

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	resp "url_short/internal/lib/api/response"
	"url_short/internal/storage"
)

//go:generate go run github.com/vektra/mockery/v2@v2.43.2 --name=URLDel
type URLDel interface {
	GetURL(alias string) (string, error)
	DeleteURL(alias string) error
}

func New(log *slog.Logger, deleteURL URLDel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handle.removed.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		alias := chi.URLParam(r, "alias")

		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		if _, err := deleteURL.GetURL(alias); errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)

			render.JSON(w, r, resp.Error("not found"))

			return
		}
		err := deleteURL.DeleteURL(alias)
		if err == nil {
			log.Error("failed to removed url")

			render.JSON(w, r, resp.Error("internal error"))

			return
		}
		log.Info("got removed url", slog.String("url", alias))
		render.JSON(w, r, resp.OK())
	}
}
