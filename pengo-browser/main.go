package main

import (
	"embed"
	"net/http"
	"pengo-proto/pengo"
	"strings"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

const pengoPrefix = "/pengo/"

type PengoHandler struct {
	app *App
}

func (h PengoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uri := "pengo://" + strings.TrimPrefix(r.URL.Path, pengoPrefix)
	if isNavigationRequest(r) {
		h.app.emitNavigated(uri)
	}
	response, err := pengo.Fetch(uri)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", response.ContentType)
	w.Write(response.Body)
}

func isNavigationRequest(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept"), "text/html")
}

func pengoMiddleware(app *App) assetserver.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, pengoPrefix) {
				PengoHandler{app: app}.ServeHTTP(w, r)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "pengo-browser",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets:     assets,
			Middleware: pengoMiddleware(app),
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
