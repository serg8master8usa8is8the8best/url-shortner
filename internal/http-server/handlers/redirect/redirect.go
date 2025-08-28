package redirect

import (
	"errors"
	"log/slog"
	"net/http"

	"sergey/url-shortner/internal/lib/api/response"
	"sergey/url-shortner/internal/storage"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.url.redirect.new"

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, response.Error("invalid request"))

			return
		}

		log = log.With(
			slog.String("op", op),
			slog.String("alias", alias),
		)

		resUrl, err := urlGetter.GetURL(alias)
		if resUrl == "" {
			log.Info("url is empty  ")

			render.JSON(w, r, response.Error("empty url"))

			return
		}

		if errors.Is(err, storage.ErrURLNotFund) {
			log.Info("url not found", "alias:", alias)
			render.JSON(w, r, response.Error("url not found"))
			return
		}
		if err != nil {
			log.Info("failed to get url ")

			render.JSON(w, r, response.Error("internal error"))

			return
		}

		http.Redirect(w, r, resUrl, http.StatusSeeOther)
	}
}
