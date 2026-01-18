package core

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var repoPath = "D:\\workspace\\go-git-client-window"
var gitService = NewGitCoreService()

// TestGetBranchLog 测试 GetBranchLog 方法
func TestGetBranchLog(t *testing.T) {
	branch, _ := gitService.GetCurrentBranch(repoPath)
	logs, _ := gitService.GetBranchLog(repoPath, branch, 100)
	if logs != nil {
		for _, item := range logs {
			fmt.Println(item)
		}
	}
	assert.NotEmpty(t, logs)
}

func TestGetCurrentBranch(t *testing.T) {
	branch, err := gitService.GetCurrentBranch(repoPath)
	assert.NoError(t, err)
	logger.Info("当前分支", "branch", branch)
	assert.True(t, branch != "")
}

func TestGetGraphHistory(t *testing.T) {
	history, _ := gitService.GetGraphHistory(repoPath, 100)
	fmt.Println(history)
}
