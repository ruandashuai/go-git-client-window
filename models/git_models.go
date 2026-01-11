package models

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

// FileDiff 文件差异信息
type FileDiff struct {
	Filename   string   `json:"filename"`
	OldContent string   `json:"oldContent"`
	NewContent string   `json:"newContent"`
	Changes    []string `json:"changes"`
}
