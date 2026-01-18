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

	"go-git-client-window/core"
	"go-git-client-window/internal/gitcmd"
	"go-git-client-window/models"
	"go-git-client-window/utils"
)

var log = slog.New(slog.NewTextHandler(os.Stdout, nil))

//go:embed all:frontend/dist
var assets embed.FS

// App struct
type App struct {
	ctx        context.Context
	gitService *core.GitCoreService
}

// NewApp creates a new App application struct
func NewApp() *App {
	gitCmdService := gitcmd.NewGitService()
	gitCoreService := core.NewGitCoreService(gitCmdService)
	return &App{
		gitService: gitCoreService,
	}
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
	return a.gitService.Init(path)
}

// GitStatus 获取Git状态
func (a *App) GitStatus(path string) (string, error) {
	return a.gitService.Status(path)
}

// GitClone 克隆Git仓库
func (a *App) GitClone(repoURL, targetPath string) (string, error) {
	return a.gitService.Clone(repoURL, targetPath)
}

// GitAdd 添加文件到暂存区
func (a *App) GitAdd(path, files string) (string, error) {
	return a.gitService.Add(path, files)
}

// GitCommit 提交更改
func (a *App) GitCommit(path, message string) (string, error) {
	return a.gitService.Commit(path, message)
}

// GitLocalBranches 获取本地分支列表
func (a *App) GitLocalBranches(path string) (string, error) {
	result, err := a.gitService.GetLocalBranches(path)
	if err != nil {
		return "", err
	}
	return utils.ToJsonString(result), nil
}

// GitRemoteBranches 获取远程分支列表
func (a *App) GitRemoteBranches(path string) (string, error) {
	result, err := a.gitService.GetRemoteBranches(path)
	if err != nil {
		return "", err
	}
	return utils.ToJsonString(result), nil
}

// GitBranchLog 获取特定分支的提交日志（oneline 格式）
func (a *App) GitBranchLog(path, branch string, limit int) (string, error) {
	result, err := a.gitService.GetBranchLog(path, branch, limit)
	if err != nil {
		return "", err
	}
	return utils.ToJsonString(result), nil
}

// GitDiffFiles 获取工作区和暂存区的变更文件列表
func (a *App) GitDiffFiles(path string) (string, error) {
	result, err := a.gitService.GetDiffFiles(path)
	if err != nil {
		return "", err
	}
	return utils.ToJsonString(result), nil
}

// GitFileDiff 查看指定文件的 diff
func (a *App) GitFileDiff(path, filename string, staged bool) (string, error) {
	result, err := a.gitService.GetFileDiff(path, filename, staged)
	if err != nil {
		return "", err
	}
	return utils.ToJsonString(result), nil
}

// GitLog 获取Git提交历史
func (a *App) GitLog(path string, limit int) (string, error) {
	result, err := a.gitService.GetLog(path, limit)
	if err != nil {
		return "", err
	}
	return utils.ToJsonString(result), nil
}

// GitShowBranchTree 获取分支树结构
func (a *App) GitShowBranchTree(path string) (string, error) {
	result, err := a.gitService.GetGraphHistory(path, 0)
	if err != nil {
		return "", err
	}
	return result, nil
}

// GitCheckout 切换分支
func (a *App) GitCheckout(path, branch string) (string, error) {
	return a.gitService.Checkout(path, branch)
}

// GitCreateBranch 创建新分支
func (a *App) GitCreateBranch(path, branch string) (string, error) {
	return a.gitService.CreateBranch(path, branch)
}

// GitGetCurrentBranch 获取当前分支
func (a *App) GitGetCurrentBranch(path string) (string, error) {
	result, err := a.gitService.GetCurrentBranch(path)
	if err != nil {
		return "", err
	}
	return result, nil
}

// GitGetStatusStructured 获取 Git 状态（结构化返回）
func (a *App) GitGetStatusStructured(path string) ([]models.GitFileStatus, error) {
	return a.gitService.GetStatusStructured(path)
}

// GitGetStagedFiles 获取已暂存文件列表
func (a *App) GitGetStagedFiles(path string) ([]models.GitFileStatus, error) {
	return a.gitService.GetStagedFiles(path)
}

// GitStageFile 暂存单个文件
func (a *App) GitStageFile(path, filename string) (string, error) {
	return a.gitService.StageFile(path, filename)
}

// GitUnstageFile 取消暂存
func (a *App) GitUnstageFile(path, filename string) (string, error) {
	return a.gitService.UnstageFile(path, filename)
}

// GitStageAll 暂存所有变更
func (a *App) GitStageAll(path string) (string, error) {
	return a.gitService.StageAll(path)
}

// GitResetFile 重置文件
func (a *App) GitResetFile(path, filename string) (string, error) {
	return a.gitService.ResetFile(path, filename)
}

// GitFetch 获取远程更新
func (a *App) GitFetch(path string) (string, error) {
	return a.gitService.Fetch(path)
}

// GitPull 拉取分支
func (a *App) GitPull(path, branch string) (string, error) {
	return a.gitService.Pull(path, branch)
}

// GitPush 推送分支
func (a *App) GitPush(path, branch string, force bool) (string, error) {
	return a.gitService.Push(path, branch, force)
}

// GitGetRemoteInfo 获取远程仓库信息
func (a *App) GitGetRemoteInfo(path string) ([]models.GitRemoteInfo, error) {
	return a.gitService.GetRemoteInfo(path)
}

// GitMerge 合并分支
func (a *App) GitMerge(path, branch string) (string, error) {
	return a.gitService.Merge(path, branch)
}

// GitRebase 变基操作
func (a *App) GitRebase(path, branch string) (string, error) {
	return a.gitService.Rebase(path, branch)
}

// GitGetMergeConflicts 获取合并冲突列表
func (a *App) GitGetMergeConflicts(path string) ([]string, error) {
	return a.gitService.GetMergeConflicts(path)
}

// GitResolveConflict 解决冲突
func (a *App) GitResolveConflict(path, filename, strategy string) (string, error) {
	return a.gitService.ResolveConflict(path, filename, strategy)
}

// GitStashList 获取 stash 列表
func (a *App) GitStashList(path string) ([]models.GitStash, error) {
	return a.gitService.GetStashList(path)
}

// GitStashSave 保存 stash
func (a *App) GitStashSave(path, message string) (string, error) {
	return a.gitService.StashSave(path, message)
}

// GitStashApply 应用 stash
func (a *App) GitStashApply(path, stashId string) (string, error) {
	return a.gitService.StashApply(path, stashId)
}

// GitStashPop 弹出 stash
func (a *App) GitStashPop(path, stashId string) (string, error) {
	return a.gitService.StashPop(path, stashId)
}

// GitStashDrop 删除 stash
func (a *App) GitStashDrop(path, stashId string) (string, error) {
	return a.gitService.StashDrop(path, stashId)
}

// GitBlame 获取文件 blame 信息
func (a *App) GitBlame(path, filename string) ([]models.GitBlameLine, error) {
	return a.gitService.GetBlame(path, filename)
}

// GitGetGraphHistory 获取图形化历史数据
func (a *App) GitGetGraphHistory(path string, limit int) (string, error) {
	result, err := a.gitService.GetGraphHistory(path, limit)
	if err != nil {
		return "", err
	}
	return result, nil
}

// GitCommitAmend 修改最后一次提交
func (a *App) GitCommitAmend(path, message string) (string, error) {
	return a.gitService.AmendCommit(path, message)
}
func (a *App) GitBranch(path string) (string, error) {
	result, err := a.gitService.GetBranches(path)
	if err != nil {
		return "", err
	}
	return utils.ToJsonString(result), nil
}
func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:             "Go Git Client",
		Width:             1000,
		Height:            600,
		DisableResize:     false,
		Fullscreen:        false,
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
