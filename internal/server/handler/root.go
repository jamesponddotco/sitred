package handler

import (
	"log/slog"
	"math/rand"
	"net/http"

	"git.sr.ht/~jamesponddotco/sitred/internal/fetch"
	"git.sr.ht/~jamesponddotco/sitred/internal/sitemap"
	"git.sr.ht/~jamesponddotco/xstd-go/xnet/xhttp"
)

// RootHandler is the HTTP handler for the root endpoint.
type RootHandler struct {
	fetchClient *fetch.Client
	logger      *slog.Logger
	sitemapURL  string
}

// NewRootHandler returns a new RootHandler instance.
func NewRootHandler(fetchClient *fetch.Client, logger *slog.Logger, sitemapURL string) *RootHandler {
	return &RootHandler{
		fetchClient: fetchClient,
		logger:      logger,
		sitemapURL:  sitemapURL,
	}
}

// ServeHTTP handles HTTP requests for the root endpoint.
func (h *RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data, err := h.fetchClient.Remote(r.Context(), h.sitemapURL)
	if err != nil {
		h.logger.LogAttrs(
			r.Context(),
			slog.LevelError,
			"error fetching sitemap",
			slog.String("url", h.sitemapURL),
			slog.String("error", err.Error()),
		)

		response := xhttp.ResponseError{
			Message: "Failed to fetch sitemap.",
			Code:    http.StatusInternalServerError,
		}

		response.Write(r.Context(), h.logger, w)

		return
	}
	defer data.Body.Close()

	uris, err := sitemap.Parse(data.Body)
	if err != nil {
		h.logger.LogAttrs(
			r.Context(),
			slog.LevelError,
			"error parsing sitemap",
			slog.String("url", h.sitemapURL),
			slog.String("error", err.Error()),
		)

		response := xhttp.ResponseError{
			Message: "Failed to parse sitemap.",
			Code:    http.StatusInternalServerError,
		}

		response.Write(r.Context(), h.logger, w)

		return
	}

	if len(uris) == 0 {
		h.logger.LogAttrs(
			r.Context(),
			slog.LevelError,
			"no URLs available for redirect",
			slog.String("url", h.sitemapURL),
		)

		response := xhttp.ResponseError{
			Message: "No URLs available for redirect.",
			Code:    http.StatusInternalServerError,
		}

		response.Write(r.Context(), h.logger, w)

		return
	}

	uri := RandomURL(uris)

	http.Redirect(w, r, uri, http.StatusFound)
}

// RandomURL returns a random URL from the provided slice of URLs.
func RandomURL(urls []string) string {
	index := rand.Intn(len(urls)) //nolint:gosec // we don't need cryptographic randomness here

	return urls[index]
}
