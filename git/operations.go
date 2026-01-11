package git

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"go-git-client-window/models"
)

// GitInit 初始化Git仓库
func GitInit(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
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
func GitStatus(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	cmd := exec.Command("git", "status")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get status: %s", string(output))
	}
	return string(output), nil
}

// GitClone 克隆Git仓库
func GitClone(repoURL, targetPath string) (string, error) {
	if repoURL == "" {
		return "", fmt.Errorf("repository URL cannot be empty")
	}
	if targetPath == "" {
		return "", fmt.Errorf("target path cannot be empty")
	}
	cmd := exec.Command("git", "clone", repoURL, targetPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("克隆失败: %s", string(output))
	}
	return string(output), nil
}

// GitAdd 添加文件到暂存区
func GitAdd(path, files string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
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
func GitCommit(path, message string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	if message == "" {
		return "", fmt.Errorf("commit message cannot be empty")
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
func GitBranch(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	cmd := exec.Command("git", "branch", "-a")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get branches: %s", string(output))
	}
	return string(output), nil
}

// GitLocalBranches 获取本地分支列表
func GitLocalBranches(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	cmd := exec.Command("git", "branch")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get local branches: %s", string(output))
	}
	return string(output), nil
}

// GitRemoteBranches 获取远程分支列表
func GitRemoteBranches(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	cmd := exec.Command("git", "branch", "-r")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get remote branches: %s", string(output))
	}
	return string(output), nil
}

// GitBranchLog 获取特定分支的提交日志（oneline 格式）
func GitBranchLog(path, branch string, limit int) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	if branch == "" {
		return "", fmt.Errorf("branch name cannot be empty")
	}
	if limit <= 0 {
		limit = 50
	}
	cmd := exec.Command("git", "log", branch, "--oneline", "-n", fmt.Sprintf("%d", limit))
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get branch log: %s", string(output))
	}
	return string(output), nil
}

// GitDiffFiles 获取工作区和暂存区的变更文件列表
func GitDiffFiles(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	// 获取工作区变更文件
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get changed files: %s", string(output))
	}
	return string(output), nil
}

// GitFileDiff 查看指定文件的 diff
func GitFileDiff(path, filename string, staged bool) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	if filename == "" {
		return "", fmt.Errorf("filename cannot be empty")
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
		return "", fmt.Errorf("failed to get file diff: %s", string(output))
	}

	return string(output), nil
}

// GitLog 获取Git提交历史
func GitLog(path string, limit int) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	if limit <= 0 {
		limit = 50
	}
	cmd := exec.Command("git", "log", "--all", "--graph", "--decorate", "--oneline", "--pretty=format:%h|%d|%s|%an|%ad", "--date=short", "-n", fmt.Sprintf("%d", limit))
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get logs: %s", string(output))
	}
	return string(output), nil
}

// GitShowBranchTree 获取分支树结构
func GitShowBranchTree(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	cmd := exec.Command("git", "log", "--all", "--graph", "--decorate", "--oneline")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get branch tree: %s", string(output))
	}
	return string(output), nil
}

// GitCheckout 切换分支
func GitCheckout(path, branch string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	if branch == "" {
		return "", fmt.Errorf("branch name cannot be empty")
	}
	cmd := exec.Command("git", "checkout", branch)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to switch branch: %s", string(output))
	}
	return string(output), nil
}

// GitCreateBranch 创建新分支
func GitCreateBranch(path, branch string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	if branch == "" {
		return "", fmt.Errorf("branch name cannot be empty")
	}
	cmd := exec.Command("git", "checkout", "-b", branch)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to create branch: %s", string(output))
	}
	return string(output), nil
}

// GitGetCurrentBranch 获取当前分支
func GitGetCurrentBranch(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %s", string(output))
	}
	return strings.TrimSpace(string(output)), nil
}

// GitGetStatusStructured 获取 Git 状态（结构化返回）
func GitGetStatusStructured(path string) ([]models.GitFileStatus, error) {
	if path == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("获取状态失败: %s", string(output))
	}

	var files []models.GitFileStatus
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

// GitGetStagedFiles 获取已暂存文件列表
func GitGetStagedFiles(path string) ([]models.GitFileStatus, error) {
	files, err := GitGetStatusStructured(path)
	if err != nil {
		return nil, err
	}

	var stagedFiles []models.GitFileStatus
	for _, file := range files {
		if file.Staged {
			stagedFiles = append(stagedFiles, file)
		}
	}
	return stagedFiles, nil
}

// GitStageFile 暂存单个文件
func GitStageFile(path, filename string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	if filename == "" {
		return "", fmt.Errorf("filename cannot be empty")
	}
	cmd := exec.Command("git", "add", filename)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to stage file: %s", string(output))
	}
	return string(output), nil
}

// GitUnstageFile 取消暂存
func GitUnstageFile(path, filename string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	if filename == "" {
		return "", fmt.Errorf("filename cannot be empty")
	}
	cmd := exec.Command("git", "reset", "HEAD", filename)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to unstage file: %s", string(output))
	}
	return string(output), nil
}

// GitStageAll 暂存所有变更
func GitStageAll(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	cmd := exec.Command("git", "add", ".")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to stage all files: %s", string(output))
	}
	return string(output), nil
}

// GitResetFile 重置文件
func GitResetFile(path, filename string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	if filename == "" {
		return "", fmt.Errorf("filename cannot be empty")
	}
	cmd := exec.Command("git", "checkout", "--", filename)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to reset file: %s", string(output))
	}
	return string(output), nil
}

// GitFetch 获取远程更新
func GitFetch(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	cmd := exec.Command("git", "fetch")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to fetch remote updates: %s", string(output))
	}
	return string(output), nil
}

// GitPull 拉取分支
func GitPull(path, branch string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	args := []string{"pull"}
	if branch != "" {
		args = append(args, "origin", branch)
	}
	cmd := exec.Command("git", args...)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to pull branch: %s", string(output))
	}
	return string(output), nil
}

// GitPush 推送分支
func GitPush(path, branch string, force bool) (string, error) {
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
	cmd := exec.Command("git", args...)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to push branch: %s", string(output))
	}
	return string(output), nil
}

// GitGetRemoteInfo 获取远程仓库信息
func GitGetRemoteInfo(path string) (map[string]string, error) {
	if path == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}
	cmd := exec.Command("git", "remote", "-v")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get remote info: %s", string(output))
	}

	remotes := make(map[string]string)
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			name := parts[0]
			url := parts[1]
			remotes[name] = url
		}
	}

	return remotes, nil
}

// GitGetRemoteInfoStruct 获取远程仓库信息（结构化）
func GitGetRemoteInfoStruct(path string) ([]models.GitRemoteInfo, error) {
	remotes, err := GitGetRemoteInfo(path)
	if err != nil {
		return nil, err
	}

	var result []models.GitRemoteInfo
	for name, url := range remotes {
		result = append(result, models.GitRemoteInfo{
			Name: name,
			URL:  url,
		})
	}

	return result, nil
}

// GitMerge 合并分支
func GitMerge(path, branch string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	if branch == "" {
		return "", fmt.Errorf("branch name cannot be empty")
	}
	cmd := exec.Command("git", "merge", branch)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to merge branch: %s", string(output))
	}
	return string(output), nil
}

// GitRebase 变基操作
func GitRebase(path, branch string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	if branch == "" {
		return "", fmt.Errorf("branch name cannot be empty")
	}
	cmd := exec.Command("git", "rebase", branch)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to rebase: %s", string(output))
	}
	return string(output), nil
}

// GitGetMergeConflicts 获取合并冲突列表
func GitGetMergeConflicts(path string) ([]string, error) {
	if path == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}
	cmd := exec.Command("git", "diff", "--name-only", "--diff-filter=U")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get conflict list: %s", string(output))
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
func GitResolveConflict(path, filename, strategy string) (string, error) {
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
		return "", fmt.Errorf("failed to resolve conflict: %s", string(output))
	}

	// 标记为已解决
	addCmd := exec.Command("git", "add", filename)
	addCmd.Dir = path
	addOutput, err := addCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to mark conflict as resolved: %s", string(addOutput))
	}

	return string(output), nil
}

// GitStashList 获取 stash 列表
func GitStashList(path string) ([]map[string]string, error) {
	if path == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}
	cmd := exec.Command("git", "stash", "list")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	outputStr := strings.TrimSpace(string(output))
	if err != nil || outputStr == "" {
		return []map[string]string{}, nil
	}

	var stashes []map[string]string
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

		stash := map[string]string{
			"ID":      stashId,
			"Message": message,
			"Branch":  branch,
			"Date":    "", // git stash list 默认不显示日期
		}
		stashes = append(stashes, stash)
	}
	return stashes, nil
}

// GitStashListStruct 获取 stash 列表（结构化）
func GitStashListStruct(path string) ([]models.GitStash, error) {
	stashes, err := GitStashList(path)
	if err != nil {
		return nil, err
	}

	var result []models.GitStash
	for _, stash := range stashes {
		result = append(result, models.GitStash{
			ID:      stash["ID"],
			Message: stash["Message"],
			Branch:  stash["Branch"],
			Date:    stash["Date"],
		})
	}

	return result, nil
}

// GitStashSave 保存 stash
func GitStashSave(path, message string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	args := []string{"stash", "save"}
	if message != "" {
		args = append(args, message)
	}
	cmd := exec.Command("git", args...)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to save stash: %s", string(output))
	}
	return string(output), nil
}

// GitStashApply 应用 stash
func GitStashApply(path, stashId string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	args := []string{"stash", "apply"}
	if stashId != "" {
		args = append(args, stashId)
	}
	cmd := exec.Command("git", args...)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to apply stash: %s", string(output))
	}
	return string(output), nil
}

// GitStashPop 弹出 stash
func GitStashPop(path, stashId string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	args := []string{"stash", "pop"}
	if stashId != "" {
		args = append(args, stashId)
	}
	cmd := exec.Command("git", args...)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to pop stash: %s", string(output))
	}
	return string(output), nil
}

// GitStashDrop 删除 stash
func GitStashDrop(path, stashId string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	if stashId == "" {
		return "", fmt.Errorf("stash ID cannot be empty")
	}
	cmd := exec.Command("git", "stash", "drop", stashId)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to drop stash: %s", string(output))
	}
	return string(output), nil
}

// GitBlame 获取文件 blame 信息
func GitBlame(path, filename string) ([]models.GitBlameLine, error) {
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
		return nil, fmt.Errorf("failed to get blame info: %s", string(output))
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
		if line == "" {
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

// GitGetGraphHistory 获取图形化历史数据
func GitGetGraphHistory(path string, limit int) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	if limit <= 0 {
		limit = 100
	}
	format := "%H|%P|%an|%ad|%s|%d"
	cmd := exec.Command("git", "log", "--all", "--graph", "--pretty=format:"+format, "--date=short", "-n", fmt.Sprintf("%d", limit))
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get history data: %s", string(output))
	}
	return string(output), nil
}

// GitCommitAmend 修改最后一次提交
func GitCommitAmend(path, message string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
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
		return "", fmt.Errorf("failed to amend commit: %s", string(output))
	}
	return string(output), nil
}
