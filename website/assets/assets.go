package assets

import (
	"embed"
)

//go:embed css img js
var AssetFiles embed.FS

// func RegisterAssetRoutes() http.Handler {
// 	router := chi.NewRouter()
// 	router.Handle("/*", http.StripPrefix("/assets/", http.FileServerFS(assets)))

// 	return router
// }
