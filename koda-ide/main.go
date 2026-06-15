package main

import (
	"embed"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()
	if len(os.Args) > 1 {
		app.SetInitialWorkspace(os.Args[1])
	}

	err := wails.Run(&options.App{
		Title:            "Koda Studio",
		Width:            1280,
		Height:           800,
		MinWidth:         900,
		MinHeight:        560,
		BackgroundColour: &options.RGBA{R: 19, G: 25, B: 32, A: 1},
		Debug: options.Debug{
			OpenInspectorOnStartup: os.Getenv("KODA_STUDIO_DEVTOOLS") == "1",
		},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:  app.startup,
		OnShutdown: app.shutdown,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
