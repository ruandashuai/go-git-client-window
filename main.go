package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

// 数据结构定义
// GitFileStatus 文件状态
type GitFileStatus struct {
	Filename string `json:"filename"`
	Status   string `json:"status"` // M/A/D/R/?
	Staged   bool   `json:"staged"`
}

// GitCommit 提交信息
type GitCommit struct {
	Hash     string   `json:"hash"`
	Message  string   `json:"message"`
	Author   string   `json:"author"`
	Date     string   `json:"date"`
	Branches []string `json:"branches"`
}

// GitBranch 分支信息
type GitBranch struct {
	Name    string `json:"name"`
	Current bool   `json:"current"`
	Remote  bool   `json:"remote"`
	Tracked string `json:"tracked"` // 跟踪的远程分支
}

// GitStash Stash 信息
type GitStash struct {
	ID      string `json:"id"`
	Message string `json:"message"`
	Branch  string `json:"branch"`
	Date    string `json:"date"`
}

// GitBlameLine Blame 行信息
type GitBlameLine struct {
	Line    int    `json:"line"`
	Hash    string `json:"hash"`
	Author  string `json:"author"`
	Date    string `json:"date"`
	Message string `json:"message"`
	Content string `json:"content"`
}

// GitRemoteInfo 远程仓库信息
type GitRemoteInfo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

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

//func GetAllRemoteBranchList(path string) (string, error) {
//
//	return ""
//}

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

// GitLocalBranches 获取本地分支列表

func (a *App) GitLocalBranches(path string) (string, error) {

	log.Printf("[GitLocalBranches] 开始获取本地分支，路径: %s", path)

	if path == "" {

		return "", fmt.Errorf("路径不能为空")

	}

	cmd := exec.Command("git", "branch")

	cmd.Dir = path

	output, err := cmd.CombinedOutput()

	if err != nil {

		return "", fmt.Errorf("获取本地分支失败: %s", string(output))

	}

	return string(output), nil

}

// GitRemoteBranches 获取远程分支列表

func (a *App) GitRemoteBranches(path string) (string, error) {

	log.Printf("[GitRemoteBranches] 开始获取远程分支，路径: %s", path)

	if path == "" {

		return "", fmt.Errorf("路径不能为空")

	}

	cmd := exec.Command("git", "branch", "-r")

	cmd.Dir = path

	output, err := cmd.CombinedOutput()

	if err != nil {

		return "", fmt.Errorf("获取远程分支失败: %s", string(output))

	}

	return string(output), nil

}

// GitBranchLog 获取特定分支的提交日志（oneline 格式）

func (a *App) GitBranchLog(path, branch string, limit int) (string, error) {

	log.Printf("[GitBranchLog] 获取分支日志，路径: %s, 分支: %s, 限制: %d", path, branch, limit)

	if path == "" {

		return "", fmt.Errorf("路径不能为空")

	}

	if branch == "" {

		return "", fmt.Errorf("分支名不能为空")

	}

	if limit <= 0 {

		limit = 50

	}

	cmd := exec.Command("git", "log", branch, "--oneline", "-n", fmt.Sprintf("%d", limit))

	cmd.Dir = path

	output, err := cmd.CombinedOutput()

	if err != nil {

		return "", fmt.Errorf("获取分支日志失败: %s", string(output))

	}

	return string(output), nil

}

// GitDiffFiles 获取工作区和暂存区的变更文件列表

func (a *App) GitDiffFiles(path string) (string, error) {

	log.Printf("[GitDiffFiles] 获取变更文件列表，路径: %s", path)

	if path == "" {

		return "", fmt.Errorf("路径不能为空")

	}

	// 获取工作区变更文件

	cmd := exec.Command("git", "status", "--porcelain")

	cmd.Dir = path

	output, err := cmd.CombinedOutput()

	if err != nil {

		return "", fmt.Errorf("获取变更文件失败: %s", string(output))

	}

	return string(output), nil

}

// GitFileDiff 查看指定文件的 diff

func (a *App) GitFileDiff(path, filename string, staged bool) (string, error) {

	log.Printf("[GitFileDiff] 查看文件 diff，路径: %s, 文件: %s, 暂存: %v", path, filename, staged)

	if path == "" {

		return "", fmt.Errorf("路径不能为空")

	}

	if filename == "" {

		return "", fmt.Errorf("文件名不能为空")

	}

	args := []string{"diff"}

	if staged {

		args = append(args, "--staged")

	}

	args = append(args, filename)

	cmd := exec.Command("git", args...)

	cmd.Dir = path

	output, err := cmd.CombinedOutput()

	if err != nil {

		return "", fmt.Errorf("查看文件 diff 失败: %s", string(output))

	}

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
	return strings.TrimSpace(string(output)), nil
}

// GitGetStatusStructured 获取 Git 状态（结构化返回）
func (a *App) GitGetStatusStructured(path string) ([]GitFileStatus, error) {
	if path == "" {
		return nil, fmt.Errorf("路径不能为空")
	}
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("获取状态失败: %s", string(output))
	}

	var files []GitFileStatus
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if len(line) < 3 {
			continue
		}

		status := line[:2]
		filename := strings.TrimSpace(line[2:])
		staged := status[0] != ' ' && status[0] != '?'
		unstaged := len(status) > 1 && status[1] != ' ' && status[1] != '?'

		// 如果有暂存区的状态
		if staged {
			files = append(files, GitFileStatus{
				Filename: filename,
				Status:   string(status[0]),
				Staged:   true,
			})
		}
		// 如果有工作区的状态（且不同于暂存区）
		if unstaged && (!staged || status[1] != status[0]) {
			files = append(files, GitFileStatus{
				Filename: filename,
				Status:   string(status[1]),
				Staged:   false,
			})
		}
		// 如果是未跟踪文件
		if status[0] == '?' && status[1] == '?' {
			files = append(files, GitFileStatus{
				Filename: filename,
				Status:   "?",
				Staged:   false,
			})
		}
	}

	return files, nil
}

// GitGetStagedFiles 获取已暂存文件列表
func (a *App) GitGetStagedFiles(path string) ([]GitFileStatus, error) {
	files, err := a.GitGetStatusStructured(path)
	if err != nil {
		return nil, err
	}

	var stagedFiles []GitFileStatus
	for _, file := range files {
		if file.Staged {
			stagedFiles = append(stagedFiles, file)
		}
	}
	return stagedFiles, nil
}

// GitStageFile 暂存单个文件
func (a *App) GitStageFile(path, filename string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("路径不能为空")
	}
	if filename == "" {
		return "", fmt.Errorf("文件名不能为空")
	}
	cmd := exec.Command("git", "add", filename)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("暂存文件失败: %s", string(output))
	}
	return string(output), nil
}

// GitUnstageFile 取消暂存
func (a *App) GitUnstageFile(path, filename string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("路径不能为空")
	}
	if filename == "" {
		return "", fmt.Errorf("文件名不能为空")
	}
	cmd := exec.Command("git", "reset", "HEAD", filename)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("取消暂存失败: %s", string(output))
	}
	return string(output), nil
}

// GitStageAll 暂存所有变更
func (a *App) GitStageAll(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("路径不能为空")
	}
	cmd := exec.Command("git", "add", ".")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("暂存所有文件失败: %s", string(output))
	}
	return string(output), nil
}

// GitResetFile 重置文件
func (a *App) GitResetFile(path, filename string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("路径不能为空")
	}
	if filename == "" {
		return "", fmt.Errorf("文件名不能为空")
	}
	cmd := exec.Command("git", "checkout", "--", filename)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("重置文件失败: %s", string(output))
	}
	return string(output), nil
}

// GitFetch 获取远程更新
func (a *App) GitFetch(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("路径不能为空")
	}
	cmd := exec.Command("git", "fetch")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("获取远程更新失败: %s", string(output))
	}
	return string(output), nil
}

// GitPull 拉取分支
func (a *App) GitPull(path, branch string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("路径不能为空")
	}
	args := []string{"pull"}
	if branch != "" {
		args = append(args, "origin", branch)
	}
	cmd := exec.Command("git", args...)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("拉取分支失败: %s", string(output))
	}
	return string(output), nil
}

// GitPush 推送分支
func (a *App) GitPush(path, branch string, force bool) (string, error) {
	if path == "" {
		return "", fmt.Errorf("路径不能为空")
	}
	if branch == "" {
		branch = "HEAD"
	}
	args := []string{"push"}
	if force {
		args = append(args, "--force")
	}
	args = append(args, "origin", branch)
	cmd := exec.Command("git", args...)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("推送分支失败: %s", string(output))
	}
	return string(output), nil
}

// GitGetRemoteInfo 获取远程仓库信息
func (a *App) GitGetRemoteInfo(path string) ([]GitRemoteInfo, error) {
	if path == "" {
		return nil, fmt.Errorf("路径不能为空")
	}
	cmd := exec.Command("git", "remote", "-v")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("获取远程信息失败: %s", string(output))
	}

	var remotes []GitRemoteInfo
	lines := strings.Split(string(output), "\n")
	remoteMap := make(map[string]string)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			name := parts[0]
			url := parts[1]
			remoteMap[name] = url
		}
	}

	for name, url := range remoteMap {
		remotes = append(remotes, GitRemoteInfo{
			Name: name,
			URL:  url,
		})
	}

	return remotes, nil
}

// GitMerge 合并分支
func (a *App) GitMerge(path, branch string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("路径不能为空")
	}
	if branch == "" {
		return "", fmt.Errorf("分支名不能为空")
	}
	cmd := exec.Command("git", "merge", branch)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("合并分支失败: %s", string(output))
	}
	return string(output), nil
}

// GitRebase 变基操作
func (a *App) GitRebase(path, branch string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("路径不能为空")
	}
	if branch == "" {
		return "", fmt.Errorf("分支名不能为空")
	}
	cmd := exec.Command("git", "rebase", branch)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("变基失败: %s", string(output))
	}
	return string(output), nil
}

// GitGetMergeConflicts 获取合并冲突列表
func (a *App) GitGetMergeConflicts(path string) ([]string, error) {
	if path == "" {
		return nil, fmt.Errorf("路径不能为空")
	}
	cmd := exec.Command("git", "diff", "--name-only", "--diff-filter=U")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("获取冲突列表失败: %s", string(output))
	}

	var conflicts []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			conflicts = append(conflicts, line)
		}
	}
	return conflicts, nil
}

// GitResolveConflict 解决冲突
func (a *App) GitResolveConflict(path, filename, strategy string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("路径不能为空")
	}
	if filename == "" {
		return "", fmt.Errorf("文件名不能为空")
	}

	var cmd *exec.Cmd
	if strategy == "ours" {
		cmd = exec.Command("git", "checkout", "--ours", filename)
	} else if strategy == "theirs" {
		cmd = exec.Command("git", "checkout", "--theirs", filename)
	} else {
		return "", fmt.Errorf("不支持的解决策略: %s", strategy)
	}

	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("解决冲突失败: %s", string(output))
	}

	// 标记为已解决
	addCmd := exec.Command("git", "add", filename)
	addCmd.Dir = path
	addOutput, err := addCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("标记冲突已解决失败: %s", string(addOutput))
	}

	return string(output), nil
}

// GitStashList 获取 stash 列表
func (a *App) GitStashList(path string) ([]GitStash, error) {
	if path == "" {
		return nil, fmt.Errorf("路径不能为空")
	}
	cmd := exec.Command("git", "stash", "list")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	outputStr := strings.TrimSpace(string(output))
	if err != nil || outputStr == "" {
		return []GitStash{}, nil
	}

	var stashes []GitStash
	lines := strings.Split(outputStr, "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// 解析格式: stash@{0}: WIP on branch: commit message
		// 或者: stash@{0}: commit message
		stashId := fmt.Sprintf("stash@{%d}", i)
		parts := strings.SplitN(line, ":", 2)
		if len(parts) > 0 {
			// 尝试提取 stash ID
			idParts := strings.Fields(parts[0])
			if len(idParts) > 0 && strings.HasPrefix(idParts[0], "stash@{") {
				stashId = idParts[0]
			}
		}

		message := ""
		branch := ""
		if len(parts) >= 2 {
			msgPart := strings.TrimSpace(parts[1])
			// 尝试解析分支信息
			if strings.HasPrefix(msgPart, "WIP on ") {
				branchParts := strings.SplitN(msgPart, ":", 2)
				if len(branchParts) >= 2 {
					branch = strings.TrimSpace(strings.TrimPrefix(branchParts[0], "WIP on "))
					message = strings.TrimSpace(branchParts[1])
				} else {
					branch = strings.TrimSpace(strings.TrimPrefix(msgPart, "WIP on "))
				}
			} else {
				message = msgPart
			}
		}

		stashes = append(stashes, GitStash{
			ID:      stashId,
			Message: message,
			Branch:  branch,
			Date:    "", // git stash list 默认不显示日期
		})
	}
	return stashes, nil
}

// GitStashSave 保存 stash
func (a *App) GitStashSave(path, message string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("路径不能为空")
	}
	args := []string{"stash", "save"}
	if message != "" {
		args = append(args, message)
	}
	cmd := exec.Command("git", args...)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("保存 stash 失败: %s", string(output))
	}
	return string(output), nil
}

// GitStashApply 应用 stash
func (a *App) GitStashApply(path, stashId string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("路径不能为空")
	}
	args := []string{"stash", "apply"}
	if stashId != "" {
		args = append(args, stashId)
	}
	cmd := exec.Command("git", args...)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("应用 stash 失败: %s", string(output))
	}
	return string(output), nil
}

// GitStashPop 弹出 stash
func (a *App) GitStashPop(path, stashId string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("路径不能为空")
	}
	args := []string{"stash", "pop"}
	if stashId != "" {
		args = append(args, stashId)
	}
	cmd := exec.Command("git", args...)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("弹出 stash 失败: %s", string(output))
	}
	return string(output), nil
}

// GitStashDrop 删除 stash
func (a *App) GitStashDrop(path, stashId string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("路径不能为空")
	}
	if stashId == "" {
		return "", fmt.Errorf("stash ID 不能为空")
	}
	cmd := exec.Command("git", "stash", "drop", stashId)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("删除 stash 失败: %s", string(output))
	}
	return string(output), nil
}

// GitBlame 获取文件 blame 信息
func (a *App) GitBlame(path, filename string) ([]GitBlameLine, error) {
	if path == "" {
		return nil, fmt.Errorf("路径不能为空")
	}
	if filename == "" {
		return nil, fmt.Errorf("文件名不能为空")
	}
	cmd := exec.Command("git", "blame", "-w", "-M", "--date=short", "--pretty=format:%H|%an|%ad|%s", filename)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("获取 blame 信息失败: %s", string(output))
	}

	// 读取文件内容
	filePath := path
	if !strings.HasSuffix(filePath, string(os.PathSeparator)) {
		filePath += string(os.PathSeparator)
	}
	filePath += filename
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %s", err.Error())
	}
	contentLines := strings.Split(string(fileContent), "\n")

	var blameLines []GitBlameLine
	outputLines := strings.Split(string(output), "\n")
	for i, line := range outputLines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 4)
		if len(parts) >= 4 {
			content := ""
			if i < len(contentLines) {
				content = contentLines[i]
			}
			blameLines = append(blameLines, GitBlameLine{
				Line:    i + 1,
				Hash:    strings.TrimSpace(parts[0]),
				Author:  strings.TrimSpace(parts[1]),
				Date:    strings.TrimSpace(parts[2]),
				Message: strings.TrimSpace(parts[3]),
				Content: content,
			})
		}
	}

	return blameLines, nil
}

// GitGetGraphHistory 获取图形化历史数据
func (a *App) GitGetGraphHistory(path string, limit int) (string, error) {
	if path == "" {
		return "", fmt.Errorf("路径不能为空")
	}
	if limit <= 0 {
		limit = 100
	}
	format := "%H|%P|%an|%ad|%s|%d"
	cmd := exec.Command("git", "log", "--all", "--graph", "--pretty=format:"+format, "--date=short", "-n", fmt.Sprintf("%d", limit))
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("获取历史数据失败: %s", string(output))
	}
	return string(output), nil
}

// GitCommitAmend 修改最后一次提交
func (a *App) GitCommitAmend(path, message string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("路径不能为空")
	}
	args := []string{"commit", "--amend"}
	if message != "" {
		args = append(args, "-m", message)
	} else {
		args = append(args, "--no-edit")
	}
	cmd := exec.Command("git", args...)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("修改提交失败: %s", string(output))
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
