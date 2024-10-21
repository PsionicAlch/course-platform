package assets

import (
	"embed"
	"net/http"

	"github.com/go-chi/chi/v5"
)

//go:embed css img js
var assets embed.FS

func RegisterAssetRoutes() http.Handler {
	router := chi.NewRouter()
	router.Handle("/*", http.StripPrefix("/assets/", http.FileServerFS(assets)))

	return router
}
