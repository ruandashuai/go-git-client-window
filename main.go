package main

import (
	"context"
	"embed"
	"log/slog"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"go-git-client-window/git"
	"go-git-client-window/utils"
)

var log = slog.New(slog.NewTextHandler(os.Stdout, nil))

//go:embed all:frontend/dist
var assets embed.FS

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// Greet 欢迎语
func (a *App) Greet() string {
	return "欢迎使用 Go Git Client!"
}

// OpenInBrowser 在默认浏览器（Chrome）中打开指定URL
func (a *App) OpenInBrowser(url string) error {
	return utils.OpenInBrowser(url)
}

// GitInit 初始化Git仓库
func (a *App) GitInit(path string) (string, error) {
	return git.GitInit(path)
}

// GitStatus 获取Git状态
func (a *App) GitStatus(path string) (string, error) {
	return git.GitStatus(path)
}

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:             "Go Git Client",
		Width:             800,
		Height:            600,
		DisableResize:     false,
		Fullscreen:        true,
		Frameless:         false,
		MinWidth:          800,
		MinHeight:         800,
		MaxWidth:          0,
		MaxHeight:         0,
		StartHidden:       false,
		HideWindowOnClose: false,
		AlwaysOnTop:       false,
		BackgroundColour:  &options.RGBA{R: 27, G: 54, B: 38, A: 1},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Menu:               &menu.Menu{},
		Logger:             nil,
		LogLevel:           0,
		LogLevelProduction: 0,
		OnStartup:          app.startup,
		OnDomReady:         app.domReady,
		OnShutdown:         app.shutdown,
		OnBeforeClose:      nil,
		Bind: []interface{}{
			app,
		},
		EnumBind:                         nil,
		WindowStartState:                 0,
		ErrorFormatter:                   nil,
		CSSDragProperty:                  "",
		CSSDragValue:                     "",
		EnableDefaultContextMenu:         false,
		EnableFraudulentWebsiteDetection: false,
		SingleInstanceLock:               nil,
	})

	if err != nil {
		log.Error("应用运行失败", "error", err)
		os.Exit(1)
	}
}

// startup is called when the app starts. The context is saved
// so we can call the context's lifecycle event methods.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) domReady(ctx context.Context) {
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
}
