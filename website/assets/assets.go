package assets

import (
	"embed"
)

//go:embed css img js
var AssetFiles embed.FS

// TODO: setup asset storage using cloudflare R2 and CDN

// func RegisterAssetRoutes() http.Handler {
// 	router := chi.NewRouter()
// 	router.Handle("/*", http.StripPrefix("/assets/", http.FileServerFS(assets)))

// 	return router
// }
