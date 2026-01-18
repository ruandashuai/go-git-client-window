package core

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"

	"go-git-client-window/models"
)

var logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

// GitCoreService Git核心服务
type GitCoreService struct{}

// NewGitCoreService 创建新的Git核心服务
func NewGitCoreService() *GitCoreService {
	return &GitCoreService{}
}

// ExecuteGitCommand 执行 Git 命令
func ExecuteGitCommand(dir string, args ...string) (string, error) {
	//在这里判断dir不能是nil或者空字符串
	if strings.TrimSpace(dir) == "" {
		return "", fmt.Errorf("dir cannot be empty")
	}
	cmd := exec.Command("git", args...)
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git command failed: %s, output: %s", err.Error(), string(output))
	}

	return string(output), nil
}

// ParseBranchLine 解析分支行
func ParseBranchLine(branchLine string) (name string, isCurrent, isRemote bool) {
	line := strings.TrimSpace(branchLine)
	isCurrent = strings.HasPrefix(line, "* ")

	// 移除当前分支标记
	cleanLine := strings.TrimPrefix(line, "* ")

	// 检查是否为远程分支
	isRemote = strings.Contains(cleanLine, "remotes/")

	// 提取分支名称
	name = strings.TrimSpace(cleanLine)

	// 如果是远程分支，进一步清理
	if isRemote {
		name = strings.TrimPrefix(name, "remotes/")
	}

	return name, isCurrent, isRemote
}

// ParseCommitLine 解析提交行 (格式: hash\x1frefs\x1fmessage\x1fauthor\x1fdate)
func ParseCommitLine(commitLine string) (*models.GitCommitRecord, error) {
	parts := strings.Split(commitLine, "\x1f")
	if len(parts) < 5 {
		return nil, fmt.Errorf("invalid commit line format: %s", commitLine)
	}

	hash := strings.TrimSpace(parts[0])
	refs := strings.TrimSpace(parts[1])
	message := strings.TrimSpace(parts[2])
	author := strings.TrimSpace(parts[3])
	date := strings.TrimSpace(parts[4])

	var branches []string
	if refs != "" {
		refsParts := strings.Split(refs, ",")
		branches = make([]string, 0, len(refsParts)) // 预分配容量提高性能
		for _, ref := range refsParts {
			branch := strings.TrimSpace(strings.Trim(ref, "()"))
			if branch != "" { // 过滤空字符串
				branches = append(branches, branch)
			}
		}
	}

	commit := models.GitCommitRecord{
		Hash:     hash,
		Message:  message,
		Author:   author,
		Date:     date,
		Branches: branches,
	}
	return &commit, nil
}

// GetBranches 获取所有分支（本地和远程）
func (s *GitCoreService) GetBranches(repoPath string) ([]models.GitBranch, error) {
	if strings.TrimSpace(repoPath) == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}

	output, err := ExecuteGitCommand(repoPath, "branch", "-a")
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
func (s *GitCoreService) GetLocalBranches(repoPath string) ([]models.GitBranch, error) {
	if strings.TrimSpace(repoPath) == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}

	output, err := ExecuteGitCommand(repoPath, "branch", "--list")
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
func (s *GitCoreService) GetRemoteBranches(repoPath string) ([]models.GitBranch, error) {
	if strings.TrimSpace(repoPath) == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}

	output, err := ExecuteGitCommand(repoPath, "branch", "-r")
	if err != nil {
		return nil, err
	}

	var branches []models.GitBranch
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		name, _, isRemote := ParseBranchLine(line)
		if name != "" && isRemote {
			branch := models.GitBranch{
				Name:    name,
				Current: false,
				Remote:  isRemote,
			}
			branches = append(branches, branch)
		}
	}

	return branches, nil
}

// GetBranchLog 获取分支提交日志
func (s *GitCoreService) GetBranchLog(repoPath, branch string, limit int) ([]models.GitCommitRecord, error) {
	limitStr := fmt.Sprintf("-%d", limit)
	output, err := ExecuteGitCommand(repoPath, "log", "--pretty=format:%H\x1f%D\x1f%s\x1f%an\x1f%ad", "--date=iso", limitStr, branch)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(output, "\n")
	var commits []models.GitCommitRecord

	for _, line := range lines {
		//trim 一下每行
		if strings.TrimSpace(line) == "" {
			continue
		}
		//解析每行
		commit, err := ParseCommitLine(line)
		if err != nil {
			logger.Error("parse commit line error", "error", err)
			continue
		}
		commits = append(commits, *commit)
	}

	return commits, nil
}

// GetStatus 获取仓库状态
func (s *GitCoreService) GetStatus(repoPath string) ([]models.GitFileStatus, error) {
	return s.GetStatusStructured(repoPath)
}

// GetChangedFiles 获取变更文件
func (s *GitCoreService) GetChangedFiles(repoPath string) ([]models.GitFileStatus, error) {
	return s.GetDiffFiles(repoPath)
}

// GetStagedFiles 获取已暂存文件
func (s *GitCoreService) GetStagedFiles(repoPath string) ([]models.GitFileStatus, error) {
	return s.GetStagedFiles(repoPath)
}

// GetRemotes 获取远程仓库信息
func (s *GitCoreService) GetRemotes(repoPath string) ([]models.GitRemoteInfo, error) {
	return s.GetRemoteInfo(repoPath)
}

// GetStashes 获取Stash列表
func (s *GitCoreService) GetStashes(repoPath string) ([]models.GitStash, error) {
	return s.GetStashList(repoPath)
}

// GetCurrentBranch 获取当前分支
func (s *GitCoreService) GetCurrentBranch(repoPath string) (string, error) {
	output, err := ExecuteGitCommand(repoPath, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", err
	}

	branch := strings.TrimSpace(output)
	if branch == "HEAD" {
		return "", fmt.Errorf("no branch is currently checked out (detached HEAD)")
	}
	return branch, nil
}

// CheckoutBranch 切换分支
func (s *GitCoreService) CheckoutBranch(repoPath, branchName string) (string, error) {
	return s.Checkout(repoPath, branchName)
}

// CreateBranch 创建分支
func (s *GitCoreService) CreateBranch(repoPath, branchName string) (string, error) {
	output, err := ExecuteGitCommand(repoPath, "checkout", "-b", branchName)
	if err != nil {
		return "", err
	}

	return output, nil
}

// CommitFiles 提交文件
func (s *GitCoreService) CommitFiles(repoPath, message string) (string, error) {
	return s.Commit(repoPath, message)
}

// AddFiles 添加文件到暂存区
func (s *GitCoreService) AddFiles(repoPath, files string) (string, error) {
	return s.Add(repoPath, files)
}

// StageFile 暂存文件
func (s *GitCoreService) StageFile(repoPath, filename string) (string, error) {
	if strings.TrimSpace(repoPath) == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	output, err := ExecuteGitCommand(repoPath, "add", filename)
	if err != nil {
		return "", err
	}

	return output, nil
}

// UnstageFile 取消暂存文件
func (s *GitCoreService) UnstageFile(repoPath, filename string) (string, error) {
	if strings.TrimSpace(repoPath) == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	output, err := ExecuteGitCommand(repoPath, "reset", "HEAD", "--", filename)
	if err != nil {
		return "", err
	}

	return output, nil
}

// Push 推送到远程
func (s *GitCoreService) Push(repoPath, branch string, force bool) (string, error) {
	if strings.TrimSpace(repoPath) == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	args := []string{"push"}
	if force {
		args = append(args, "--force")
	}
	args = append(args, "origin", branch)

	output, err := ExecuteGitCommand(repoPath, args...)
	if err != nil {
		return "", err
	}

	return output, nil
}

// Pull 从远程拉取
func (s *GitCoreService) Pull(repoPath, branch string) (string, error) {
	if strings.TrimSpace(repoPath) == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	output, err := ExecuteGitCommand(repoPath, "pull", "origin", branch)
	if err != nil {
		return "", err
	}

	return output, nil
}

// Fetch 获取远程更新
func (s *GitCoreService) Fetch(repoPath string) (string, error) {
	if strings.TrimSpace(repoPath) == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	output, err := ExecuteGitCommand(repoPath, "fetch")
	if err != nil {
		return "", err
	}

	return output, nil
}

// Merge 合并分支
func (s *GitCoreService) Merge(repoPath, branch string) (string, error) {
	if strings.TrimSpace(repoPath) == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	output, err := ExecuteGitCommand(repoPath, "merge", branch)
	if err != nil {
		return "", err
	}

	return output, nil
}

// Init 初始化仓库
func (s *GitCoreService) Init(repoPath string) (string, error) {
	if strings.TrimSpace(repoPath) == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	output, err := ExecuteGitCommand(repoPath, "init")
	if err != nil {
		return "", err
	}

	return output, nil
}

// Status 获取仓库状态
func (s *GitCoreService) Status(repoPath string) (string, error) {
	if strings.TrimSpace(repoPath) == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	output, err := ExecuteGitCommand(repoPath, "status")
	if err != nil {
		return "", err
	}

	return output, nil
}

// Clone 克隆仓库
func (s *GitCoreService) Clone(repoURL, targetPath string) (string, error) {
	if strings.TrimSpace(repoURL) == "" || strings.TrimSpace(targetPath) == "" {
		return "", fmt.Errorf("repo URL and target path cannot be empty")
	}

	cmd := exec.Command("git", "clone", repoURL, targetPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git clone failed: %s, output: %s", err.Error(), string(output))
	}

	return string(output), nil
}

// Add 添加文件到暂存区
func (s *GitCoreService) Add(repoPath, files string) (string, error) {
	if strings.TrimSpace(repoPath) == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	output, err := ExecuteGitCommand(repoPath, "add", files)
	if err != nil {
		return "", err
	}

	return output, nil
}

// Commit 提交更改
func (s *GitCoreService) Commit(repoPath, message string) (string, error) {
	if strings.TrimSpace(repoPath) == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	output, err := ExecuteGitCommand(repoPath, "commit", "-m", message)
	if err != nil {
		return "", err
	}

	return output, nil
}

// GetDiffFiles 获取工作区和暂存区的变更文件列表
func (s *GitCoreService) GetDiffFiles(repoPath string) ([]models.GitFileStatus, error) {
	if strings.TrimSpace(repoPath) == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}

	output, err := ExecuteGitCommand(repoPath, "diff", "--name-status")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(output, "\n")
	var files []models.GitFileStatus

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 2 {
			status := fields[0]
			filename := strings.Join(fields[1:], " ")
			files = append(files, models.GitFileStatus{
				Filename: filename,
				Status:   status,
			})
		}
	}

	return files, nil
}

// GetFileDiff 查看指定文件的 diff
func (s *GitCoreService) GetFileDiff(repoPath, filename string, staged bool) ([]models.FileDiff, error) {
	if strings.TrimSpace(repoPath) == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}

	var args []string
	if staged {
		args = []string{"diff", "--cached", "--", filename}
	} else {
		args = []string{"diff", "--", filename}
	}

	output, err := ExecuteGitCommand(repoPath, args...)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(output, "\n")
	var diffs []models.FileDiff

	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			diffs = append(diffs, models.FileDiff{Changes: []string{line}})
		}
	}

	return diffs, nil
}

// GetGraphHistory 获取图形化历史数据
func (s *GitCoreService) GetGraphHistory(repoPath string, limit int) (string, error) {
	if strings.TrimSpace(repoPath) == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	limitStr := fmt.Sprintf("-%d", limit)
	output, err := ExecuteGitCommand(repoPath, "log", "--graph", "--oneline", "--all", limitStr)
	if err != nil {
		return "", err
	}

	return output, nil
}

// StageAll 暂存所有变更
func (s *GitCoreService) StageAll(repoPath string) (string, error) {
	if strings.TrimSpace(repoPath) == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	output, err := ExecuteGitCommand(repoPath, "add", "--all")
	if err != nil {
		return "", err
	}

	return output, nil
}

// ResetFile 重置文件
func (s *GitCoreService) ResetFile(repoPath, filename string) (string, error) {
	if strings.TrimSpace(repoPath) == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	output, err := ExecuteGitCommand(repoPath, "checkout", "HEAD", "--", filename)
	if err != nil {
		return "", err
	}

	return output, nil
}

// GetMergeConflicts 获取合并冲突列表
func (s *GitCoreService) GetMergeConflicts(repoPath string) ([]string, error) {
	if strings.TrimSpace(repoPath) == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}

	output, err := ExecuteGitCommand(repoPath, "diff", "--name-only", "--diff-filter=U")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(output, "\n")
	var conflicts []string

	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			conflicts = append(conflicts, strings.TrimSpace(line))
		}
	}

	return conflicts, nil
}

// ResolveConflict 解决冲突
func (s *GitCoreService) ResolveConflict(repoPath, filename, strategy string) (string, error) {
	if strings.TrimSpace(repoPath) == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	var args []string
	switch strategy {
	case "ours":
		args = []string{"checkout", "--ours", "--", filename}
	case "theirs":
		args = []string{"checkout", "--theirs", "--", filename}
	default:
		args = []string{"add", filename} // 默认为标记为已解决
	}

	output, err := ExecuteGitCommand(repoPath, args...)
	if err != nil {
		return "", err
	}

	return output, nil
}

// StashSave 保存 stash
func (s *GitCoreService) StashSave(repoPath, message string) (string, error) {
	if strings.TrimSpace(repoPath) == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	args := []string{"stash", "save"}
	if message != "" {
		args = append(args, message)
	}

	output, err := ExecuteGitCommand(repoPath, args...)
	if err != nil {
		return "", err
	}

	return output, nil
}

// StashApply 应用 stash
func (s *GitCoreService) StashApply(repoPath, stashId string) (string, error) {
	if strings.TrimSpace(repoPath) == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	output, err := ExecuteGitCommand(repoPath, "stash", "apply", stashId)
	if err != nil {
		return "", err
	}

	return output, nil
}

// StashPop 弹出 stash
func (s *GitCoreService) StashPop(repoPath, stashId string) (string, error) {
	if strings.TrimSpace(repoPath) == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	output, err := ExecuteGitCommand(repoPath, "stash", "pop", stashId)
	if err != nil {
		return "", err
	}

	return output, nil
}

// StashDrop 删除 stash
func (s *GitCoreService) StashDrop(repoPath, stashId string) (string, error) {
	if strings.TrimSpace(repoPath) == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	output, err := ExecuteGitCommand(repoPath, "stash", "drop", stashId)
	if err != nil {
		return "", err
	}

	return output, nil
}

// GetBlame 获取文件 blame 信息
func (s *GitCoreService) GetBlame(repoPath, filename string) ([]models.GitBlameLine, error) {
	if strings.TrimSpace(repoPath) == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}

	output, err := ExecuteGitCommand(repoPath, "blame", "-p", filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(output, "\n")
	var blameLines []models.GitBlameLine

	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			blameLines = append(blameLines, models.GitBlameLine{Content: line})
		}
	}

	return blameLines, nil
}

// AmendCommit 修改最后一次提交
func (s *GitCoreService) AmendCommit(repoPath, message string) (string, error) {
	if strings.TrimSpace(repoPath) == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	output, err := ExecuteGitCommand(repoPath, "commit", "--amend", "-m", message)
	if err != nil {
		return "", err
	}

	return output, nil
}

// Checkout 切换分支
func (s *GitCoreService) Checkout(repoPath, branch string) (string, error) {
	if strings.TrimSpace(repoPath) == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	output, err := ExecuteGitCommand(repoPath, "checkout", branch)
	if err != nil {
		return "", err
	}

	return output, nil
}

// GetStatusStructured 获取 Git 状态（结构化返回）
func (s *GitCoreService) GetStatusStructured(repoPath string) ([]models.GitFileStatus, error) {
	if strings.TrimSpace(repoPath) == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}

	output, err := ExecuteGitCommand(repoPath, "status", "--porcelain", "-uall")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(output, "\n")
	var files []models.GitFileStatus

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 2 {
			status := fields[0]
			filename := strings.Join(fields[1:], " ")
			files = append(files, models.GitFileStatus{
				Filename: filename,
				Status:   status,
			})
		}
	}

	return files, nil
}

// GetRemoteInfo 获取远程仓库信息
func (s *GitCoreService) GetRemoteInfo(repoPath string) ([]models.GitRemoteInfo, error) {
	if strings.TrimSpace(repoPath) == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}

	output, err := ExecuteGitCommand(repoPath, "remote", "-v")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(output, "\n")
	var remotes []models.GitRemoteInfo

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 2 {
			name := fields[0]
			url := strings.Trim(fields[1], "()")
			remote := models.GitRemoteInfo{
				Name: name,
				URL:  url,
			}
			remotes = append(remotes, remote)
		}
	}

	return remotes, nil
}

// Rebase 变基操作
func (s *GitCoreService) Rebase(repoPath, branch string) (string, error) {
	if strings.TrimSpace(repoPath) == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	output, err := ExecuteGitCommand(repoPath, "rebase", branch)
	if err != nil {
		return "", err
	}

	return output, nil
}

// GetStashList 获取 stash 列表
func (s *GitCoreService) GetStashList(repoPath string) ([]models.GitStash, error) {
	if strings.TrimSpace(repoPath) == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}

	output, err := ExecuteGitCommand(repoPath, "stash", "list")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(output, "\n")
	var stashes []models.GitStash

	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		stashes = append(stashes, models.GitStash{
			ID:      fmt.Sprintf("stash@{%d}", i),
			Message: line,
		})
	}

	return stashes, nil
}

// GetGraphHistoryWithFormat 获取图形化历史数据
func (s *GitCoreService) GetGraphHistoryWithFormat(repoPath string, limit int) (string, error) {
	if strings.TrimSpace(repoPath) == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	limitStr := fmt.Sprintf("-%d", limit)
	output, err := ExecuteGitCommand(repoPath, "log", "--graph", "--oneline", "--all", limitStr)
	if err != nil {
		return "", err
	}

	return output, nil
}

// ShowBranchTree 获取分支树结构
func (s *GitCoreService) ShowBranchTree(repoPath string) (string, error) {
	if strings.TrimSpace(repoPath) == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	output, err := ExecuteGitCommand(repoPath, "log", "--graph", "--oneline", "--all", "--decorate")
	if err != nil {
		return "", err
	}

	return output, nil
}
