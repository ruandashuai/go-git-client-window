package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

// App struct
type App struct {
	ctx        context.Context
	httpServer *http.Server
	serverPort int
	isRunning  bool
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		serverPort: 8080,
	}
}

// Greet 欢迎语
func (a *App) Greet() string {
	return fmt.Sprintf("欢迎使用 Go Git Client!")
}

// OpenInBrowser 在默认浏览器（Chrome）中打开指定URL
func (a *App) OpenInBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		// 尝试使用Chrome打开
		chromePaths := []string{
			"C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe",
			"C:\\Program Files (x86)\\Google\\Chrome\\Application\\chrome.exe",
			"C:\\Users\\" + os.Getenv("USERNAME") + "\\AppData\\Local\\Google\\Chrome\\Application\\chrome.exe",
		}
		chromePath := ""
		for _, path := range chromePaths {
			if _, err := os.Stat(path); err == nil {
				chromePath = path
				break
			}
		}
		if chromePath != "" {
			cmd = exec.Command(chromePath, url)
		} else {
			// 如果找不到Chrome，使用默认浏览器
			cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
		}
	case "darwin":
		cmd = exec.Command("open", "-a", "Google Chrome", url)
	default: // linux
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Start()
}

// StartHTTPServer 启动HTTP服务器
func (a *App) StartHTTPServer(port int) (string, error) {
	if a.isRunning {
		return fmt.Sprintf("HTTP服务器已在运行，端口: %d", a.serverPort), nil
	}

	a.serverPort = port
	mux := http.NewServeMux()

	// 静态文件服务
	mux.Handle("/", http.FileServer(http.FS(assets)))

	// API路由
	mux.HandleFunc("/api/status", a.handleGitStatus)
	mux.HandleFunc("/api/init", a.handleGitInit)
	mux.HandleFunc("/api/clone", a.handleGitClone)
	mux.HandleFunc("/api/add", a.handleGitAdd)
	mux.HandleFunc("/api/commit", a.handleGitCommit)
	mux.HandleFunc("/api/server-info", a.handleServerInfo)

	a.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTP服务器错误: %v\n", err)
		}
	}()

	a.isRunning = true
	return fmt.Sprintf("HTTP服务器已启动，访问地址: http://localhost:%d", port), nil
}

// StopHTTPServer 停止HTTP服务器
func (a *App) StopHTTPServer() (string, error) {
	if !a.isRunning {
		return "HTTP服务器未运行", nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.httpServer.Shutdown(ctx); err != nil {
		return "", fmt.Errorf("停止服务器失败: %v", err)
	}

	a.isRunning = false
	return "HTTP服务器已停止", nil
}

// GetServerInfo 获取服务器信息
func (a *App) GetServerInfo() (map[string]interface{}, error) {
	return map[string]interface{}{
		"running": a.isRunning,
		"port":    a.serverPort,
		"url":     fmt.Sprintf("http://localhost:%d", a.serverPort),
	}, nil
}

// API响应结构
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 发送JSON响应
func sendJSONResponse(w http.ResponseWriter, success bool, message string, data interface{}) {
	response := APIResponse{
		Success: success,
		Message: message,
		Data:    data,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// 处理服务器信息请求
func (a *App) handleServerInfo(w http.ResponseWriter, r *http.Request) {
	info, _ := a.GetServerInfo()
	sendJSONResponse(w, true, "服务器信息", info)
}

// 处理Git状态请求
func (a *App) handleGitStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendJSONResponse(w, false, "只支持POST请求", nil)
		return
	}

	var req struct {
		Path string `json:"path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONResponse(w, false, "请求解析失败", err.Error())
		return
	}

	result, err := a.GitStatus(req.Path)
	if err != nil {
		sendJSONResponse(w, false, err.Error(), nil)
		return
	}

	sendJSONResponse(w, true, "获取状态成功", result)
}

// 处理Git初始化请求
func (a *App) handleGitInit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendJSONResponse(w, false, "只支持POST请求", nil)
		return
	}

	var req struct {
		Path string `json:"path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONResponse(w, false, "请求解析失败", err.Error())
		return
	}

	result, err := a.GitInit(req.Path)
	if err != nil {
		sendJSONResponse(w, false, err.Error(), nil)
		return
	}

	sendJSONResponse(w, true, "初始化成功", result)
}

// 处理Git克隆请求
func (a *App) handleGitClone(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendJSONResponse(w, false, "只支持POST请求", nil)
		return
	}

	var req struct {
		RepoURL    string `json:"repoUrl"`
		TargetPath string `json:"targetPath"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONResponse(w, false, "请求解析失败", err.Error())
		return
	}

	result, err := a.GitClone(req.RepoURL, req.TargetPath)
	if err != nil {
		sendJSONResponse(w, false, err.Error(), nil)
		return
	}

	sendJSONResponse(w, true, "克隆成功", result)
}

// 处理Git添加请求
func (a *App) handleGitAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendJSONResponse(w, false, "只支持POST请求", nil)
		return
	}

	var req struct {
		Path  string `json:"path"`
		Files string `json:"files"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONResponse(w, false, "请求解析失败", err.Error())
		return
	}

	result, err := a.GitAdd(req.Path, req.Files)
	if err != nil {
		sendJSONResponse(w, false, err.Error(), nil)
		return
	}

	sendJSONResponse(w, true, "添加成功", result)
}

// 处理Git提交请求
func (a *App) handleGitCommit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendJSONResponse(w, false, "只支持POST请求", nil)
		return
	}

	var req struct {
		Path    string `json:"path"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONResponse(w, false, "请求解析失败", err.Error())
		return
	}

	result, err := a.GitCommit(req.Path, req.Message)
	if err != nil {
		sendJSONResponse(w, false, err.Error(), nil)
		return
	}

	sendJSONResponse(w, true, "提交成功", result)
}

// GitInit 初始化Git仓库
func (a *App) GitInit(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("路径不能为空")
	}
	cmd := exec.Command("git", "init")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("初始化失败: %s", string(output))
	}
	return string(output), nil
}

// GitStatus 获取Git状态
func (a *App) GitStatus(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("路径不能为空")
	}
	cmd := exec.Command("git", "status")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("获取状态失败: %s", string(output))
	}
	return string(output), nil
}

// GitClone 克隆Git仓库
func (a *App) GitClone(repoURL, targetPath string) (string, error) {
	if repoURL == "" {
		return "", fmt.Errorf("仓库URL不能为空")
	}
	if targetPath == "" {
		return "", fmt.Errorf("目标路径不能为空")
	}
	cmd := exec.Command("git", "clone", repoURL, targetPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("克隆失败: %s", string(output))
	}
	return string(output), nil
}

// GitAdd 添加文件到暂存区
func (a *App) GitAdd(path, files string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("路径不能为空")
	}
	if files == "" {
		files = "."
	}
	cmd := exec.Command("git", "add", files)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("添加文件失败: %s", string(output))
	}
	return string(output), nil
}

// GitCommit 提交更改
func (a *App) GitCommit(path, message string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("路径不能为空")
	}
	if message == "" {
		return "", fmt.Errorf("提交信息不能为空")
	}
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("提交失败: %s", string(output))
	}
	return string(output), nil
}

// GitBranch 获取Git分支列表

func (a *App) GitBranch(path string) (string, error) {

	log.Printf("[GitBranch] 开始获取分支，路径: %s", path)

	if path == "" {

		log.Println("[GitBranch] 错误: 路径不能为空")

		return "", fmt.Errorf("路径不能为空")

	}

	cmd := exec.Command("git", "branch", "-a")

	cmd.Dir = path

	output, err := cmd.CombinedOutput()

	log.Printf("[GitBranch] 命令执行完成，输出长度: %d", len(output))

	if err != nil {

		log.Printf("[GitBranch] 错误: %v, 输出: %s", err, string(output))

		return "", fmt.Errorf("获取分支失败: %s", string(output))

	}

	log.Printf("[GitBranch] 成功获取分支，分支数: %d", len(string(output)))

	return string(output), nil

}

// GitLog 获取Git提交历史
func (a *App) GitLog(path string, limit int) (string, error) {
	if path == "" {
		return "", fmt.Errorf("路径不能为空")
	}
	if limit <= 0 {
		limit = 50
	}
	cmd := exec.Command("git", "log", "--all", "--graph", "--decorate", "--oneline", "--pretty=format:%h|%d|%s|%an|%ad", "--date=short", "-n", fmt.Sprintf("%d", limit))
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("获取日志失败: %s", string(output))
	}
	return string(output), nil
}

// GitShowBranchTree 获取分支树结构
func (a *App) GitShowBranchTree(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("路径不能为空")
	}
	cmd := exec.Command("git", "log", "--all", "--graph", "--decorate", "--oneline")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("获取分支树失败: %s", string(output))
	}
	return string(output), nil
}

// GitCheckout 切换分支
func (a *App) GitCheckout(path, branch string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("路径不能为空")
	}
	if branch == "" {
		return "", fmt.Errorf("分支名不能为空")
	}
	cmd := exec.Command("git", "checkout", branch)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("切换分支失败: %s", string(output))
	}
	return string(output), nil
}

// GitCreateBranch 创建新分支
func (a *App) GitCreateBranch(path, branch string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("路径不能为空")
	}
	if branch == "" {
		return "", fmt.Errorf("分支名不能为空")
	}
	cmd := exec.Command("git", "checkout", "-b", branch)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("创建分支失败: %s", string(output))
	}
	return string(output), nil
}

// GitGetCurrentBranch 获取当前分支
func (a *App) GitGetCurrentBranch(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("路径不能为空")
	}
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("获取当前分支失败: %s", string(output))
	}
	return string(output), nil
}

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Go Git Client",
		Width:  800,
		Height: 600,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		OnDomReady:       app.domReady,
		OnShutdown:       app.shutdown,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

// startup is called when the app starts. The context is saved
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// domReady is called when the front end is ready
func (a *App) domReady(ctx context.Context) {
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
}
