package utils

import (
	"os"
	"os/exec"
	"runtime"
)

// OpenInBrowser 在默认浏览器（Chrome）中打开指定URL
func OpenInBrowser(url string) error {
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
