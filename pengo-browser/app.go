package main

import (
	"context"
	"fmt"
	"pengo-proto/pengo"

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

func (a *App) PengoFetch(input string) (pengo.Response, error) {
	return pengo.Fetch(input)
}

func (a *App) emitNavigated(uri string, title string) {
	if a.ctx == nil {
		return
	}
	runtime.EventsEmit(a.ctx, navigatedEvent, uri, title)
}
