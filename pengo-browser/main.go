package main

import (
	"bytes"
	"embed"
	"net/http"
	"pengo-proto/pengo"
	"strings"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"golang.org/x/net/html"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed frontend/src/pages/connection-error.html
var connectionErrorPage []byte

const pengoPrefix = "/pengo/"
const pengoScheme = "pengo://"

type PengoHandler struct {
	app *App
}

func getTitle(response pengo.Response) string {
	tokenizer := html.NewTokenizer(bytes.NewReader(response.Body))
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			break
		}
		tokenTag, _ := tokenizer.TagName()
		if string(tokenTag) == "title" {
			titleNext := tokenizer.Next()
			if titleNext == html.TextToken {
				titleText := tokenizer.Text()				
				return string(titleText)
			}
		}
	}
	return ""
}

func (h PengoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host, requestPath := splitHostAndPath(strings.TrimPrefix(r.URL.Path, pengoPrefix))
	requestUri := "pengo://" + strings.TrimPrefix(r.URL.Path, pengoPrefix)
	response, err := pengo.Fetch("fetch", host, requestPath, nil)
	if err != nil {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusBadGateway)
		w.Write(connectionErrorPage)
		return
	}
	
	if isNavigationRequest(r) {
		//trigger an event to get navigated data in frontend
		h.app.emitNavigated(requestUri, getTitle(response))
	}
	w.Header().Set("Content-Type", response.ContentType)
	w.Write(response.Body)
}

// "welcome/about"  ->  host "welcome", path "/about"
func splitHostAndPath(hostAndPath string) (host string, requestPath string) {
	uriBody := strings.SplitN(hostAndPath, "/", 2)
	host = uriBody[0]
	requestPath = "/"
	if len(uriBody) > 1 && uriBody[1] != "" {
		requestPath = "/" + uriBody[1]
	}
	return host, requestPath
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
