package core

import (
	"fmt"
	"os"
	"path/filepath"
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

func TestCreateBranchFrom(t *testing.T) {
	tempDir := t.TempDir()
	fmt.Println("Created a temp dir:")
	fmt.Println(tempDir)

	//初始化仓库测试
	t.Run("初始化仓库", func(t *testing.T) {
		initResult, err := gitService.Init(tempDir)
		fmt.Println("Init repository result:")
		fmt.Println(initResult)
		assert.NoError(t, err)
		assert.Contains(t, initResult, "Initialized empty Git repository")
	})

	t.Run("创建master", func(t *testing.T) {
		newBranchName := "master"
		result, err := gitService.CreateBranch(tempDir, newBranchName)
		fmt.Println("[CreateBranch] result:")
		fmt.Println(result)
		assert.NoError(t, err)
		assert.Contains(t, result, "Switched to a new branch")
		assert.Contains(t, result, newBranchName)
	})

	t.Run("Git add ", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "test.txt")
		file, err := os.Create(filePath)
		_ = file.Close()
		assert.NoError(t, err)
		assert.NotNil(t, file)
		fmt.Println("Created a file:", file.Name())

		addResult, err := gitService.Add(tempDir, file.Name())
		assert.NoError(t, err)
		assert.NotNil(t, addResult)
		fmt.Println("Git add result:")
		fmt.Println(addResult)
	})

	t.Run("Commit", func(t *testing.T) {
		commitResult, err := gitService.Commit(tempDir, "Initial commit")
		assert.NoError(t, err)
		fmt.Println("[Commit] result:")
		fmt.Println(commitResult)
	})

	t.Run("获取当前分支", func(t *testing.T) {
		currentBranch, err := gitService.GetCurrentBranch(tempDir)
		fmt.Println("[GetCurrentBranch] result:")
		fmt.Println(currentBranch)
		assert.NoError(t, err)
		assert.Equal(t, "master", currentBranch)
	})
}
