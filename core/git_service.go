package core

import (
	"go-git-client-window/internal/gitcmd"
	"go-git-client-window/models"
)

// GitCoreService Git核心服务
type GitCoreService struct {
	gitService *gitcmd.GitService
}

// NewGitCoreService 创建新的Git核心服务
func NewGitCoreService(gitService *gitcmd.GitService) *GitCoreService {
	return &GitCoreService{
		gitService: gitService,
	}
}

// GetBranches 获取所有分支（本地和远程）
func (s *GitCoreService) GetBranches(repoPath string) ([]models.GitBranch, error) {
	return s.gitService.GetBranches(repoPath)
}

// GetLocalBranches 获取本地分支
func (s *GitCoreService) GetLocalBranches(repoPath string) ([]models.GitBranch, error) {
	return s.gitService.GetLocalBranches(repoPath)
}

// GetRemoteBranches 获取远程分支
func (s *GitCoreService) GetRemoteBranches(repoPath string) ([]models.GitBranch, error) {
	return s.gitService.GetRemoteBranches(repoPath)
}

// GetBranchLog 获取分支提交日志
func (s *GitCoreService) GetBranchLog(repoPath, branch string, limit int) ([]models.GitCommit, error) {
	return s.gitService.GetBranchLog(repoPath, branch, limit)
}

// GetCommits 获取提交历史
func (s *GitCoreService) GetCommits(repoPath string, limit int) ([]models.GitCommit, error) {
	return s.gitService.GetLog(repoPath, limit)
}

// GetStatus 获取仓库状态
func (s *GitCoreService) GetStatus(repoPath string) ([]models.GitFileStatus, error) {
	return s.gitService.GetStatusStructured(repoPath)
}

// GetChangedFiles 获取变更文件
func (s *GitCoreService) GetChangedFiles(repoPath string) ([]models.GitFileStatus, error) {
	return s.gitService.GetDiffFiles(repoPath)
}

// GetStagedFiles 获取已暂存文件
func (s *GitCoreService) GetStagedFiles(repoPath string) ([]models.GitFileStatus, error) {
	return s.gitService.GetStagedFiles(repoPath)
}

// GetRemotes 获取远程仓库信息
func (s *GitCoreService) GetRemotes(repoPath string) ([]models.GitRemoteInfo, error) {
	return s.gitService.GetRemoteInfo(repoPath)
}

// GetStashes 获取Stash列表
func (s *GitCoreService) GetStashes(repoPath string) ([]models.GitStash, error) {
	return s.gitService.GetStashList(repoPath)
}

// GetCurrentBranch 获取当前分支
func (s *GitCoreService) GetCurrentBranch(repoPath string) (string, error) {
	return s.gitService.GetCurrentBranch(repoPath)
}

// CheckoutBranch 切换分支
func (s *GitCoreService) CheckoutBranch(repoPath, branchName string) (string, error) {
	return s.gitService.Checkout(repoPath, branchName)
}

// CreateBranch 创建分支
func (s *GitCoreService) CreateBranch(repoPath, branchName string) (string, error) {
	return s.gitService.CreateBranch(repoPath, branchName)
}

// CommitFiles 提交文件
func (s *GitCoreService) CommitFiles(repoPath, message string) (string, error) {
	return s.gitService.Commit(repoPath, message)
}

// AddFiles 添加文件到暂存区
func (s *GitCoreService) AddFiles(repoPath, files string) (string, error) {
	return s.gitService.Add(repoPath, files)
}

// StageFile 暂存文件
func (s *GitCoreService) StageFile(repoPath, filename string) (string, error) {
	return s.gitService.StageFile(repoPath, filename)
}

// UnstageFile 取消暂存文件
func (s *GitCoreService) UnstageFile(repoPath, filename string) (string, error) {
	return s.gitService.UnstageFile(repoPath, filename)
}

// Push 推送到远程
func (s *GitCoreService) Push(repoPath, branch string, force bool) (string, error) {
	return s.gitService.Push(repoPath, branch, force)
}

// Pull 从远程拉取
func (s *GitCoreService) Pull(repoPath, branch string) (string, error) {
	return s.gitService.Pull(repoPath, branch)
}

// Fetch 获取远程更新
func (s *GitCoreService) Fetch(repoPath string) (string, error) {
	return s.gitService.Fetch(repoPath)
}

// Merge 合并分支
func (s *GitCoreService) Merge(repoPath, branch string) (string, error) {
	return s.gitService.Merge(repoPath, branch)
}

// Init 初始化仓库
func (s *GitCoreService) Init(repoPath string) (string, error) {
	return s.gitService.Init(repoPath)
}

// Status 获取仓库状态
func (s *GitCoreService) Status(repoPath string) (string, error) {
	return s.gitService.Status(repoPath)
}

// Clone 克隆仓库
func (s *GitCoreService) Clone(repoURL, targetPath string) (string, error) {
	return s.gitService.Clone(repoURL, targetPath)
}

// Add 添加文件到暂存区
func (s *GitCoreService) Add(repoPath, files string) (string, error) {
	return s.gitService.Add(repoPath, files)
}

// Commit 提交更改
func (s *GitCoreService) Commit(repoPath, message string) (string, error) {
	return s.gitService.Commit(repoPath, message)
}

// GetDiffFiles 获取工作区和暂存区的变更文件列表
func (s *GitCoreService) GetDiffFiles(repoPath string) ([]models.GitFileStatus, error) {
	return s.gitService.GetDiffFiles(repoPath)
}

// GetFileDiff 查看指定文件的 diff
func (s *GitCoreService) GetFileDiff(repoPath, filename string, staged bool) ([]models.FileDiff, error) {
	return s.gitService.GetFileDiff(repoPath, filename, staged)
}

// GetLog 获取Git提交历史
func (s *GitCoreService) GetLog(repoPath string, limit int) ([]models.GitCommit, error) {
	return s.gitService.GetLog(repoPath, limit)
}

// GetGraphHistory 获取图形化历史数据
func (s *GitCoreService) GetGraphHistory(repoPath string, limit int) (string, error) {
	return s.gitService.GetGraphHistory(repoPath, limit)
}

// StageAll 暂存所有变更
func (s *GitCoreService) StageAll(repoPath string) (string, error) {
	return s.gitService.StageAll(repoPath)
}

// ResetFile 重置文件
func (s *GitCoreService) ResetFile(repoPath, filename string) (string, error) {
	return s.gitService.ResetFile(repoPath, filename)
}

// GetMergeConflicts 获取合并冲突列表
func (s *GitCoreService) GetMergeConflicts(repoPath string) ([]string, error) {
	return s.gitService.GetMergeConflicts(repoPath)
}

// ResolveConflict 解决冲突
func (s *GitCoreService) ResolveConflict(repoPath, filename, strategy string) (string, error) {
	return s.gitService.ResolveConflict(repoPath, filename, strategy)
}

// StashSave 保存 stash
func (s *GitCoreService) StashSave(repoPath, message string) (string, error) {
	return s.gitService.StashSave(repoPath, message)
}

// StashApply 应用 stash
func (s *GitCoreService) StashApply(repoPath, stashId string) (string, error) {
	return s.gitService.StashApply(repoPath, stashId)
}

// StashPop 弹出 stash
func (s *GitCoreService) StashPop(repoPath, stashId string) (string, error) {
	return s.gitService.StashPop(repoPath, stashId)
}

// StashDrop 删除 stash
func (s *GitCoreService) StashDrop(repoPath, stashId string) (string, error) {
	return s.gitService.StashDrop(repoPath, stashId)
}

// GetBlame 获取文件 blame 信息
func (s *GitCoreService) GetBlame(repoPath, filename string) ([]models.GitBlameLine, error) {
	return s.gitService.GetBlame(repoPath, filename)
}

// AmendCommit 修改最后一次提交
func (s *GitCoreService) AmendCommit(repoPath, message string) (string, error) {
	return s.gitService.AmendCommit(repoPath, message)
}

// Checkout 切换分支
func (s *GitCoreService) Checkout(repoPath, branch string) (string, error) {
	return s.gitService.Checkout(repoPath, branch)
}

// GetStatusStructured 获取 Git 状态（结构化返回）
func (s *GitCoreService) GetStatusStructured(repoPath string) ([]models.GitFileStatus, error) {
	return s.gitService.GetStatusStructured(repoPath)
}

// GetRemoteInfo 获取远程仓库信息
func (s *GitCoreService) GetRemoteInfo(repoPath string) ([]models.GitRemoteInfo, error) {
	return s.gitService.GetRemoteInfo(repoPath)
}

// Rebase 变基操作
func (s *GitCoreService) Rebase(repoPath, branch string) (string, error) {
	return s.gitService.Rebase(repoPath, branch)
}

// GetStashList 获取 stash 列表
func (s *GitCoreService) GetStashList(repoPath string) ([]models.GitStash, error) {
	return s.gitService.GetStashList(repoPath)
}
