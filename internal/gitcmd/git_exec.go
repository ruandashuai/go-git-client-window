package gitcmd

import (
	"fmt"
	"go-git-client-window/models"
	"os/exec"
	"strings"
)

// ExecuteGitCommand 执行 Git 命令
func ExecuteGitCommand(dir string, args ...string) (string, error) {
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

// ParseCommitLine 解析提交行 (格式: hash|refs|message|author|date)
func ParseCommitLine(commitLine string) (*models.GitCommit, error) {
	parts := strings.Split(commitLine, "|")
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
		branches = strings.Split(refs, ",")
		for i, branch := range branches {
			branches[i] = strings.TrimSpace(strings.Trim(branch, "()"))
		}
	}

	return &models.GitCommit{
		Hash:     hash,
		Message:  message,
		Author:   author,
		Date:     date,
		Branches: branches,
	}, nil
}
