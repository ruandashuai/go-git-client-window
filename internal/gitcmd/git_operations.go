package gitcmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"go-git-client-window/models"
)

// GitService Git服务接口
type GitService struct{}

// NewGitService 创建新的Git服务实例
func NewGitService() *GitService {
	return &GitService{}
}

// Init 初始化Git仓库
func (g *GitService) Init(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	return ExecuteGitCommand(path, "init")
}

// Status 获取Git状态
func (g *GitService) Status(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	return ExecuteGitCommand(path, "status")
}

// Clone 克隆仓库
func (g *GitService) Clone(repoURL, targetPath string) (string, error) {
	if repoURL == "" {
		return "", fmt.Errorf("repository URL cannot be empty")
	}
	if targetPath == "" {
		return "", fmt.Errorf("target path cannot be empty")
	}
	return ExecuteGitCommand(targetPath, "clone", repoURL, targetPath)
}

// Add 添加文件到暂存区
func (g *GitService) Add(path, files string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	if files == "" {
		files = "."
	}
	return ExecuteGitCommand(path, "add", files)
}

// Commit 提交更改
func (g *GitService) Commit(path, message string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	if message == "" {
		return "", fmt.Errorf("commit message cannot be empty")
	}
	return ExecuteGitCommand(path, "commit", "-m", message)
}

// GetBranches 获取所有分支（本地和远程）
func (g *GitService) GetBranches(path string) ([]models.GitBranch, error) {
	if path == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}

	output, err := ExecuteGitCommand(path, "branch", "-a")
	if err != nil {
		return nil, err
	}

	var branches []models.GitBranch
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		name, isCurrent, isRemote := ParseBranchLine(line)
		if name != "" {
			branch := models.GitBranch{
				Name:    name,
				Current: isCurrent,
				Remote:  isRemote,
			}
			branches = append(branches, branch)
		}
	}

	return branches, nil
}

// GetLocalBranches 获取本地分支
func (g *GitService) GetLocalBranches(path string) ([]models.GitBranch, error) {
	if path == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}

	output, err := ExecuteGitCommand(path, "branch")
	if err != nil {
		return nil, err
	}

	var branches []models.GitBranch
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		name, isCurrent, _ := ParseBranchLine(line)
		if name != "" {
			branch := models.GitBranch{
				Name:    name,
				Current: isCurrent,
				Remote:  false,
			}
			branches = append(branches, branch)
		}
	}

	return branches, nil
}

// GetRemoteBranches 获取远程分支
func (g *GitService) GetRemoteBranches(path string) ([]models.GitBranch, error) {
	if path == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}

	output, err := ExecuteGitCommand(path, "branch", "-r")
	if err != nil {
		return nil, err
	}

	var branches []models.GitBranch
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		name, isCurrent, isRemote := ParseBranchLine(line)
		if name != "" {
			branch := models.GitBranch{
				Name:    name,
				Current: isCurrent,
				Remote:  isRemote,
			}
			branches = append(branches, branch)
		}
	}

	return branches, nil
}

// GetBranchLog 获取分支提交日志
func (g *GitService) GetBranchLog(path, branch string, limit int) ([]models.GitCommit, error) {
	if path == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}
	if branch == "" {
		return nil, fmt.Errorf("branch name cannot be empty")
	}
	if limit <= 0 {
		limit = 50
	}

	output, err := ExecuteGitCommand(path, "log", branch, "--oneline", "-n", fmt.Sprintf("%d", limit), "--pretty=format:%h|%d|%s|%an|%ad", "--date=short")
	if err != nil {
		return nil, err
	}

	var commits []models.GitCommit
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		commit, err := ParseCommitLine(line)
		if err != nil {
			continue // 跳过无法解析的行
		}
		commits = append(commits, *commit)
	}

	return commits, nil
}

// GetDiffFiles 获取变更文件
func (g *GitService) GetDiffFiles(path string) ([]models.GitFileStatus, error) {
	if path == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}

	output, err := ExecuteGitCommand(path, "status", "--porcelain")
	if err != nil {
		return nil, err
	}

	var files []models.GitFileStatus
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
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
			files = append(files, models.GitFileStatus{
				Filename: filename,
				Status:   string(status[0]),
				Staged:   true,
			})
		}
		// 如果有工作区的状态（且不同于暂存区）
		if unstaged && (!staged || status[1] != status[0]) {
			files = append(files, models.GitFileStatus{
				Filename: filename,
				Status:   string(status[1]),
				Staged:   false,
			})
		}
		// 如果是未跟踪文件
		if status[0] == '?' && status[1] == '?' {
			files = append(files, models.GitFileStatus{
				Filename: filename,
				Status:   "?",
				Staged:   false,
			})
		}
	}

	return files, nil
}

// GetFileDiff 获取文件差异
func (g *GitService) GetFileDiff(path, filename string, staged bool) ([]models.FileDiff, error) {
	if path == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}
	if filename == "" {
		return nil, fmt.Errorf("filename cannot be empty")
	}

	args := []string{"diff"}
	if staged {
		args = append(args, "--staged")
	}
	args = append(args, filename)

	diffOutput, err := ExecuteGitCommand(path, args...)
	if err != nil {
		return nil, err
	}

	// 解析diff输出并返回结构化数据
	return []models.FileDiff{{
		Filename:   filename,
		NewContent: diffOutput,
	}}, nil
}

// GetLog 获取提交历史
func (g *GitService) GetLog(path string, limit int) ([]models.GitCommit, error) {
	if path == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}
	if limit <= 0 {
		limit = 50
	}

	output, err := ExecuteGitCommand(path, "log", "--all", "--graph", "--decorate", "--oneline", "--pretty=format:%h|%d|%s|%an|%ad", "--date=short", "-n", fmt.Sprintf("%d", limit))
	if err != nil {
		return nil, err
	}

	var commits []models.GitCommit
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		commit, err := ParseCommitLine(line)
		if err != nil {
			continue // 跳过无法解析的行
		}
		commits = append(commits, *commit)
	}

	return commits, nil
}

// Checkout 切换分支
func (g *GitService) Checkout(path, branch string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	if branch == "" {
		return "", fmt.Errorf("branch name cannot be empty")
	}
	return ExecuteGitCommand(path, "checkout", branch)
}

// CreateBranch 创建分支
func (g *GitService) CreateBranch(path, branch string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	if branch == "" {
		return "", fmt.Errorf("branch name cannot be empty")
	}
	return ExecuteGitCommand(path, "checkout", "-b", branch)
}

// GetCurrentBranch 获取当前分支
func (g *GitService) GetCurrentBranch(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	output, err := ExecuteGitCommand(path, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(output), nil
}

// GetStatusStructured 获取结构化状态
func (g *GitService) GetStatusStructured(path string) ([]models.GitFileStatus, error) {
	return g.GetDiffFiles(path)
}

// GetStagedFiles 获取已暂存文件
func (g *GitService) GetStagedFiles(path string) ([]models.GitFileStatus, error) {
	if path == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}

	output, err := ExecuteGitCommand(path, "diff", "--cached", "--name-status")
	if err != nil {
		return nil, err
	}

	var files []models.GitFileStatus
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.Split(line, "\t")
		if len(parts) == 2 {
			files = append(files, models.GitFileStatus{
				Filename: strings.TrimSpace(parts[1]),
				Status:   strings.TrimSpace(parts[0]),
				Staged:   true,
			})
		}
	}

	return files, nil
}

// StageFile 暂存文件
func (g *GitService) StageFile(path, filename string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	if filename == "" {
		return "", fmt.Errorf("filename cannot be empty")
	}
	return ExecuteGitCommand(path, "add", filename)
}

// UnstageFile 取消暂存文件
func (g *GitService) UnstageFile(path, filename string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	if filename == "" {
		return "", fmt.Errorf("filename cannot be empty")
	}
	return ExecuteGitCommand(path, "reset", "HEAD", filename)
}

// StageAll 暂存所有变更
func (g *GitService) StageAll(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	return ExecuteGitCommand(path, "add", ".")
}

// ResetFile 重置文件
func (g *GitService) ResetFile(path, filename string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	if filename == "" {
		return "", fmt.Errorf("filename cannot be empty")
	}
	return ExecuteGitCommand(path, "checkout", "--", filename)
}

// Fetch 获取远程更新
func (g *GitService) Fetch(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	return ExecuteGitCommand(path, "fetch")
}

// Pull 拉取更新
func (g *GitService) Pull(path, branch string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	args := []string{"pull"}
	if branch != "" {
		args = append(args, "origin", branch)
	}
	return ExecuteGitCommand(path, args...)
}

// Push 推送更新
func (g *GitService) Push(path, branch string, force bool) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	if branch == "" {
		branch = "HEAD"
	}
	args := []string{"push"}
	if force {
		args = append(args, "--force")
	}
	args = append(args, "origin", branch)
	return ExecuteGitCommand(path, args...)
}

// GetRemoteInfo 获取远程信息
func (g *GitService) GetRemoteInfo(path string) ([]models.GitRemoteInfo, error) {
	if path == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}

	output, err := ExecuteGitCommand(path, "remote", "-v")
	if err != nil {
		return nil, err
	}

	var remotes []models.GitRemoteInfo
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) >= 2 {
			// 提取远程名称和URL
			name := parts[0]
			url := strings.Trim(parts[1], "()")

			// 只添加唯一的远程仓库（避免重复的fetch/push条目）
			exists := false
			for _, remote := range remotes {
				if remote.Name == name && remote.URL == url {
					exists = true
					break
				}
			}

			if !exists {
				remotes = append(remotes, models.GitRemoteInfo{
					Name: name,
					URL:  url,
				})
			}
		}
	}

	return remotes, nil
}

// Merge 合并分支
func (g *GitService) Merge(path, branch string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	if branch == "" {
		return "", fmt.Errorf("branch name cannot be empty")
	}
	return ExecuteGitCommand(path, "merge", branch)
}

// Rebase 变基操作
func (g *GitService) Rebase(path, branch string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	if branch == "" {
		return "", fmt.Errorf("branch name cannot be empty")
	}
	return ExecuteGitCommand(path, "rebase", branch)
}

// GetMergeConflicts 获取合并冲突
func (g *GitService) GetMergeConflicts(path string) ([]string, error) {
	if path == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}

	output, err := ExecuteGitCommand(path, "diff", "--name-only", "--diff-filter=U")
	if err != nil {
		return nil, err
	}

	var conflicts []string
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			conflicts = append(conflicts, line)
		}
	}

	return conflicts, nil
}

// ResolveConflict 解决冲突
func (g *GitService) ResolveConflict(path, filename, strategy string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	if filename == "" {
		return "", fmt.Errorf("filename cannot be empty")
	}

	var cmd *exec.Cmd
	if strategy == "ours" {
		cmd = exec.Command("git", "checkout", "--ours", filename)
	} else if strategy == "theirs" {
		cmd = exec.Command("git", "checkout", "--theirs", filename)
	} else {
		return "", fmt.Errorf("unsupported conflict resolution strategy: %s", strategy)
	}

	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to resolve conflict: %s, output: %s", err.Error(), string(output))
	}

	// 标记为已解决
	addCmd := exec.Command("git", "add", filename)
	addCmd.Dir = path
	addOutput, err := addCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to mark conflict as resolved: %s, output: %s", err.Error(), string(addOutput))
	}

	return string(output), nil
}

// GetStashList 获取Stash列表
func (g *GitService) GetStashList(path string) ([]models.GitStash, error) {
	if path == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}

	output, err := ExecuteGitCommand(path, "stash", "list")
	if err != nil {
		return nil, err
	}

	var stashes []models.GitStash
	lines := strings.Split(strings.TrimSpace(output), "\n")

	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		// 解析格式: stash@{0}: WIP on branch: commit message
		stashId := fmt.Sprintf("stash@{%d}", i)
		parts := strings.SplitN(line, ": ", 2)

		if len(parts) >= 2 {
			stashId = strings.TrimSpace(parts[0])
			message := strings.TrimSpace(parts[1])

			// 尝试提取分支信息
			branch := ""
			if strings.Contains(message, "on ") {
				onIndex := strings.Index(message, "on ")
				if onIndex != -1 {
					remaining := message[onIndex+3:]
					endIndex := strings.Index(remaining, ":")
					if endIndex != -1 {
						branch = strings.TrimSpace(remaining[:endIndex])
					}
				}
			}

			stash := models.GitStash{
				ID:      stashId,
				Message: message,
				Branch:  branch,
			}
			stashes = append(stashes, stash)
		}
	}

	return stashes, nil
}

// StashSave 保存Stash
func (g *GitService) StashSave(path, message string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	args := []string{"stash", "save"}
	if message != "" {
		args = append(args, message)
	}

	return ExecuteGitCommand(path, args...)
}

// StashApply 应用Stash
func (g *GitService) StashApply(path, stashId string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	if stashId == "" {
		return "", fmt.Errorf("stash ID cannot be empty")
	}

	return ExecuteGitCommand(path, "stash", "apply", stashId)
}

// StashPop 弹出Stash
func (g *GitService) StashPop(path, stashId string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	if stashId == "" {
		return "", fmt.Errorf("stash ID cannot be empty")
	}

	return ExecuteGitCommand(path, "stash", "pop", stashId)
}

// StashDrop 删除Stash
func (g *GitService) StashDrop(path, stashId string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	if stashId == "" {
		return "", fmt.Errorf("stash ID cannot be empty")
	}

	return ExecuteGitCommand(path, "stash", "drop", stashId)
}

// GetBlame 获取文件Blame信息
func (g *GitService) GetBlame(path, filename string) ([]models.GitBlameLine, error) {
	if path == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}
	if filename == "" {
		return nil, fmt.Errorf("filename cannot be empty")
	}

	cmd := exec.Command("git", "blame", "-w", "-M", "--date=short", "--pretty=format:%H|%an|%ad|%s", filename)
	cmd.Dir = path

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get blame info: %s, output: %s", err.Error(), string(output))
	}

	// 读取文件内容
	var filePath string
	if runtime.GOOS == "windows" {
		if !strings.HasSuffix(path, "\\") {
			filePath = path + "\\" + filename
		} else {
			filePath = path + filename
		}
	} else {
		if !strings.HasSuffix(path, "/") {
			filePath = path + "/" + filename
		} else {
			filePath = path + filename
		}
	}

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %s", err.Error())
	}
	contentLines := strings.Split(string(fileContent), "\n")

	var blameLines []models.GitBlameLine
	outputLines := strings.Split(string(output), "\n")

	for i, line := range outputLines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.SplitN(line, "|", 4)
		if len(parts) >= 4 {
			content := ""
			if i < len(contentLines) {
				content = contentLines[i]
			}

			blameLine := models.GitBlameLine{
				Line:    i + 1,
				Hash:    strings.TrimSpace(parts[0]),
				Author:  strings.TrimSpace(parts[1]),
				Date:    strings.TrimSpace(parts[2]),
				Message: strings.TrimSpace(parts[3]),
				Content: content,
			}
			blameLines = append(blameLines, blameLine)
		}
	}

	return blameLines, nil
}

// GetGraphHistory 获取图形化历史
func (g *GitService) GetGraphHistory(path string, limit int) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	if limit <= 0 {
		limit = 100
	}

	format := "%H|%P|%an|%ad|%s|%d"
	return ExecuteGitCommand(path, "log", "--all", "--graph", "--pretty=format:"+format, "--date=short", "-n", fmt.Sprintf("%d", limit))
}

// AmendCommit 修改最后一次提交
func (g *GitService) AmendCommit(path, message string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	args := []string{"commit", "--amend"}
	if message != "" {
		args = append(args, "-m", message)
	} else {
		args = append(args, "--no-edit")
	}

	return ExecuteGitCommand(path, args...)
}
