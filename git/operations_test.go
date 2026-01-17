package git

import (
	"fmt"
	"os"
	"testing"
)

// TestGitOperations 测试 Git 操作函数
func TestGitOperations(t *testing.T) {
	t.Run("TestGitStatus", func(t *testing.T) {
		rootDir := t.TempDir()
		gitInitResult, err := GitInit(rootDir)
		if err != nil {
			t.Fatalf("GitInit failed: %v", err)
		}
		fmt.Println("GitInit Success")
		fmt.Println(gitInitResult)

		file, err := os.Create(rootDir + "/test.txt")
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		println("Create test file Success: ", file.Name())

		defer file.Close()

		_, err = file.WriteString("hello world")
		if err != nil {
			t.Fatalf("Failed to write to test file: %v", err)
		}
		statusResult, err := GitStatus(rootDir)
		if err != nil {
			t.Fatalf("GitStatus failed: %v", err)
		}
		fmt.Println(statusResult)

	})
}
