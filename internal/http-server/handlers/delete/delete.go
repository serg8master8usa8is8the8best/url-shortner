package delete

import (
	"errors"
	"log/slog"
	"net/http"

	"sergey/url-shortner/internal/lib/api/response"
	"sergey/url-shortner/internal/storage"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type Response struct {
	response.Response
	Alias string `json:"alias, omitempty"`
}

type UrlDeleter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDeleter UrlDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.url.delete.new"

		alias := chi.URLParam(r, "alias")

		if alias == "" {
			log.Info("alias are empty")

			render.JSON(w, r, response.Error("alias are empty"))

			return
		}

		log = log.With(
			slog.String("op", op),
			slog.String("alias", alias),
		)

		err := urlDeleter.DeleteURL(alias)

		if errors.Is(err, storage.ErrURLNotFund) {

			log.Info("alias not found", "alias:", alias)

			render.JSON(w, r, response.Error("alias not found:"))
			return
		}

		if err != nil {
			log.Error("failed to delete alias", alias)

			render.JSON(w, r, response.Error("bd error"))
		}

		responseOK(w, r, alias)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: response.OK(),
		Alias:    alias,
	})
}
