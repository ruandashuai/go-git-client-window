<template>
  <div class="app-container">
    <div class="toolbar">
      <input
          type="text"
          v-model="repoPath"
          placeholder="Enter Git repository path"
          @keyup.enter="loadRepo"
      />
      <button @click="loadRepo">加载仓库</button>
      <button @click="refreshData">Refresh</button>
    </div>

    <div class="main-content">
      <div class="sidebar">
        <div class="sidebar-header">
          <span>🌳 Branch Tree</span>
          <div class="action-buttons">
            <button @click="createBranch">New Branch</button>
            <button class="refresh-btn" @click="refreshBranches">↻</button>
          </div>
        </div>
        <div class="branch-tree" id="branchTree">
          <div v-if="branchesLoading" class="loading">加载中...</div>
          <div v-else-if="branches.length === 0" class="error-message">请先加载仓库...</div>
          <div v-else>
            <!-- 本地分支 -->
            <div>
              <h4>本地分支</h4>
              <div v-if="localBranches.length > 0" class="branch-group">
                <div
                    v-for="branch in localBranches"
                    :key="branch.name"
                    :class="['tree-item', { 'active': branch.current }]"
                    @click="switchBranch(branch.name)"
                >
                  <span class="icon">🌿</span>
                  <span class="branch-name">{{ branch.name }}</span>
                  <span v-if="branch.current" class="branch-tag">Current</span>
                </div>
              </div>
              <div v-else>
                <div class="error-message">没有本地分支</div>
              </div>

            </div>
            <div>
              <h4>远程分支</h4>
              <div v-if="remoteBranches.length > 0" class="branch-group">
                <div
                    v-for="branch in remoteBranches"
                    :key="branch.name"
                    :class="['tree-item', { 'active': branch.current }]"
                    @click="switchBranch(branch.name)"
                >
                  <span class="icon">📡</span>
                  <span class="branch-name">{{ branch.name }}</span>
                  <span v-if="branch.current" class="branch-tag">Current</span>
                </div>
              </div>
              <div v-else>
                <div class="error-message">没有远程分支</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="commit-history">
        <div class="sidebar-header"
             style="background: transparent; border-bottom: 1px solid #4e5254; margin-bottom: 16px;">
          <span>📝 Commit History</span>
          <div class="action-buttons">
            <button @click="showStatus">View Status</button>
            <button class="refresh-btn" @click="refreshCommits">↻</button>
          </div>
        </div>
        <div id="commitHistory">
          <div v-if="commitsLoading" class="loading">加载中...</div>
          <div v-else-if="commits.length === 0" class="error-message">请先加载仓库...</div>
          <div v-else>
            <div
                v-for="commit in commits"
                :key="commit.hash"
                class="commit-item"
            >
              <div class="commit-header">
                <span class="commit-hash">{{ commit.hash }}</span>
                <span class="commit-author">{{ commit.author }}</span>
                <span class="commit-date">{{ commit.date }}</span>
              </div>
              <div class="commit-message">{{ commit.message }}</div>
              <div class="commit-branches">
                <span
                    v-for="branch in commit.branches"
                    :key="branch"
                    class="commit-branch"
                >
                  {{ branch }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="status-bar">
      <div>
        <span
            :class="['status-indicator', {
            'connected': repoPath, 
            'disconnected': !repoPath 
          }]"
        ></span>
        <span id="currentBranch">{{ currentBranch || '未加载仓库' }}</span>
      </div>
      <div id="repoInfo">{{ repoPath }}</div>
    </div>

    <!-- Notification component -->
    <div
        :class="['notification', notification.type, { show: notification.visible }]"
        v-if="notification.visible"
    >
      {{ notification.message }}
    </div>
  </div>
</template>

<script>
import {computed, reactive, ref} from 'vue'

// 导入Wails运行时和Go模块
// 注意：在生产环境中，Wails会在运行时注入这些对象，所以不需要显式导入
// 我们将在运行时通过window.go访问这些函数

export default {
  name: 'App',
  setup() {
    // State variables
    const repoPath = ref('D:/workspace/go-git-client-window')
    const currentBranch = ref('')
    // {
    //   "name": "master",
    //     "current": true,
    //     "remote": false,
    //     "tracked": ""
    // }
    const allBranches = ref([])
    const commits = ref([])
    const branchesLoading = ref(false)
    const commitsLoading = ref(false)

    // Notification state
    const notification = reactive({
      visible: false,
      message: '',
      type: 'info' // 'success', 'error', 'info'
    })

    // 计算属性：将分支分为本地分支和远程分支
    const localBranches = computed(() => {
      return allBranches.value.filter(branch => !branch.remote)
    })

    const remoteBranches = computed(() => {
      return allBranches.value.filter(branch => branch.remote)
    })

    // Methods
    const showNotification = (message, type = 'info') => {
      notification.message = message
      notification.type = type
      notification.visible = true

      // Auto-hide after 3 seconds
      setTimeout(() => {
        notification.visible = false
      }, 3000)
    }

    const loadRepo = async () => {
      if (!repoPath.value.trim()) {
        showNotification('Please enter a repository path', 'error')
        return
      }

      await refreshData()
    }

    const refreshData = async () => {
      if (!repoPath.value) return

      try {
        await Promise.all([
          loadBranches(),
          loadCommits(),
          loadCurrentBranch()
        ])
      } catch (error) {
        console.error('Error refreshing data:', error)
      }
    }

    const loadBranches = async () => {
      if (!repoPath.value) return
      branchesLoading.value = true
      try {
        const result = await window.go.main.App.GitBranch(repoPath.value)
        allBranches.value = JSON.parse(result);
      } catch (error) {
        showNotification(`Failed to load branches: ${error}`, 'error')
      } finally {
        branchesLoading.value = false
      }
    }

    const loadCommits = async () => {
      if (!repoPath.value) return

      commitsLoading.value = true

      try {
        const result = await window.go.main.App.GitLog(repoPath.value, 100)
        const commitList = result.split('\n').filter(c => c.trim())

        // Parse commits and create commit objects
        const parsedCommits = commitList.map(commit => {
          const parts = commit.split('|')
          if (parts.length >= 5) {
            const [hash, refs, message, author, date] = parts
            const branchTags = refs.split(',').filter(r => r.trim())

            // Clean up branch tags
            const branches = branchTags.map(tag => tag.trim().replace(/[()]/g, '')).filter(tag => tag)

            return {
              hash,
              message,
              author,
              date,
              branches
            }
          }
          return null
        }).filter(Boolean) // Remove any null values

        commits.value = parsedCommits
      } catch (error) {
        showNotification(`Failed to load commit history: ${error}`, 'error')
      } finally {
        commitsLoading.value = false
      }
    }

    const loadCurrentBranch = async () => {
      if (!repoPath.value) return

      try {
        const result = await window.go.main.App.GitGetCurrentBranch(repoPath.value)
        currentBranch.value = result.trim()
      } catch (error) {
        console.error('Failed to get current branch:', error)
      }
    }

    const switchBranch = async (branchName) => {
      if (!repoPath.value) return
      if (!confirm(`Are you sure to switch to branch "${branchName}"?`)) return

      try {
        await window.go.main.App.GitCheckout(repoPath.value, branchName)
        await refreshData()
        showNotification(`Switched to branch ${branchName}`, 'success')
      } catch (error) {
        showNotification(`Failed to switch branch: ${error}`, 'error')
      }
    }

    const createBranch = async () => {
      if (!repoPath.value) {
        showNotification('Please load a repository first', 'error')
        return
      }

      const branchName = prompt('Please enter new branch name:')
      if (!branchName) return

      try {
        await window.go.main.App.GitCreateBranch(repoPath.value, branchName)
        await refreshData()
        showNotification(`Created and switch to branch ${branchName}`, 'success')
      } catch (error) {
        showNotification(`Failed to create branch: ${error}`, 'error')
      }
    }

    const showStatus = async () => {
      if (!repoPath.value) {
        showNotification('Please load a repository first', 'error')
        return
      }

      try {
        const result = await window.go.main.App.GitStatus(repoPath.value)
        alert('Git Status:\n\n' + result)
      } catch (error) {
        showNotification(`Failed to get status: ${error}`, 'error')
      }
    }

    const refreshBranches = async () => {
      await loadBranches()
    }

    const refreshCommits = async () => {
      await loadCommits()
    }

    return {
      repoPath,
      currentBranch,
      branches: allBranches,
      commits,
      branchesLoading,
      commitsLoading,
      localBranches,
      remoteBranches,
      notification,
      loadRepo,
      refreshData,
      loadBranches,
      loadCommits,
      loadCurrentBranch,
      switchBranch,
      createBranch,
      showStatus,
      refreshBranches,
      refreshCommits,
      showNotification
    }
  }
}
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
  background-color: #2b2b2b;
  color: #a9b7c6;
  height: 100vh;
  overflow: hidden;
}

.app-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
}

.toolbar {
  background-color: #3c3f41;
  padding: 8px 16px;
  border-bottom: 1px solid #4e5254;
  display: flex;
  align-items: center;
  gap: 10px;
}

.toolbar input {
  background-color: #3c3f41;
  border: 1px solid #5e6366;
  color: #a9b7c6;
  padding: 6px 12px;
  border-radius: 4px;
  font-size: 13px;
  width: 300px;
}

.toolbar input:focus {
  outline: none;
  border-color: #4a6d8c;
}

.toolbar button {
  background-color: #4a6d8c;
  border: none;
  color: #ffffff;
  padding: 6px 16px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 13px;
  transition: background-color 0.2s;
}

.toolbar button:hover {
  background-color: #5c7d9e;
}

.main-content {
  display: flex;
  flex: 1;
  overflow: hidden;
}

.sidebar {
  width: 350px;
  background-color: #313335;
  border-right: 1px solid #4e5254;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.sidebar-header {
  background-color: #3c3f41;
  padding: 10px 16px;
  border-bottom: 1px solid #4e5254;
  font-weight: 600;
  font-size: 14px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.branch-tree {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}

.branch-tree .branch-group {
  margin-bottom: 20px;
}

.branch-tree .branch-group h4 {
  margin: 0 0 8px 0;
  padding: 4px 8px;
  color: #6a8759;
  font-size: 12px;
  font-weight: bold;
  border-bottom: 1px solid #4e5254;
}

.tree-item {
  padding: 6px 12px;
  cursor: pointer;
  border-radius: 4px;
  margin-bottom: 2px;
  font-size: 13px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.tree-item:hover {
  background-color: #4e5254;
}

.tree-item.active {
  background-color: #4a6d8c;
  color: #ffffff;
}

.tree-item .icon {
  font-size: 14px;
}

.tree-item .branch-name {
  flex: 1;
}

.tree-item .branch-tag {
  background-color: #6a8759;
  color: #ffffff;
  padding: 2px 6px;
  border-radius: 3px;
  font-size: 11px;
}

.tree-item .remote-tag {
  background-color: #cc7832;
  color: #ffffff;
  padding: 2px 6px;
  border-radius: 3px;
  font-size: 11px;
}

.commit-history {
  flex: 1;
  background-color: #2b2b2b;
  overflow-y: auto;
  padding: 16px;
}

.commit-item {
  background-color: #313335;
  border: 1px solid #4e5254;
  border-radius: 6px;
  padding: 12px 16px;
  margin-bottom: 12px;
  cursor: pointer;
  transition: border-color 0.2s;
}

.commit-item:hover {
  border-color: #4a6d8c;
}

.commit-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.commit-hash {
  color: #9876aa;
  font-family: 'Courier New', monospace;
  font-size: 12px;
}

.commit-author {
  color: #cc7832;
  font-size: 12px;
}

.commit-date {
  color: #808080;
  font-size: 12px;
}

.commit-message {
  color: #a9b7c6;
  font-size: 14px;
  margin-bottom: 8px;
  font-weight: 500;
}

.commit-branches {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.commit-branch {
  background-color: #4a6d8c;
  color: #ffffff;
  padding: 2px 8px;
  border-radius: 3px;
  font-size: 11px;
}

.status-bar {
  background-color: #3c3f41;
  padding: 6px 16px;
  border-top: 1px solid #4e5254;
  font-size: 12px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.status-indicator {
  display: inline-block;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  margin-right: 6px;
}

.status-indicator.connected {
  background-color: #6a8759;
}

.status-indicator.disconnected {
  background-color: #cc7832;
}

.loading {
  text-align: center;
  padding: 40px;
  color: #808080;
}

.error-message {
  color: #cc7832;
  padding: 16px;
  text-align: center;
}

.action-buttons {
  display: flex;
  gap: 6px;
}

.action-buttons button {
  background-color: #4a6d8c;
  border: none;
  color: #ffffff;
  padding: 4px 10px;
  border-radius: 3px;
  cursor: pointer;
  font-size: 11px;
}

.action-buttons button:hover {
  background-color: #5c7d9e;
}

.refresh-btn {
  background-color: transparent;
  border: 1px solid #5e6366;
  color: #a9b7c6;
  padding: 4px 8px;
  border-radius: 3px;
  cursor: pointer;
  font-size: 12px;
}

.refresh-btn:hover {
  background-color: #4e5254;
}

.notification {
  position: fixed;
  top: 20px;
  right: 20px;
  padding: 12px 20px;
  border-radius: 4px;
  color: white;
  font-size: 14px;
  z-index: 1000;
  opacity: 0;
  transform: translateX(100%);
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.notification.show {
  opacity: 1;
  transform: translateX(0);
}

.notification.success {
  background-color: #6a8759;
}

.notification.error {
  background-color: #cc7832;
}

.notification.info {
  background-color: #4a6d8c;
}
</style>