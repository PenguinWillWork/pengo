package main

import (
	"context"
	"fmt"
	"pengo/protocol"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const navigatedEvent = "pengo:navigated"

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// Frontend entry point: takes a pengo:// uri, e.g. pengo://welcome/favicon.ico
func (a *App) PengoFetch(input string) (protocol.Response, error) {
	host, requestPath := splitHostAndPath(strings.TrimPrefix(input, pengoScheme))
	return protocol.Fetch("fetch", host, requestPath, nil)
}

func (a *App) emitNavigated(uri string, title string) {
	if a.ctx == nil {
		return
	}
	runtime.EventsEmit(a.ctx, navigatedEvent, uri, title)
}
