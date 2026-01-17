<template>
  <div class="app-container">
    <!-- 顶部导航栏 -->
    <header class="top-nav">
      <div class="nav-left">
        <h1 class="app-title">
          <span class="git-icon">📊</span>
          Git 客户端
        </h1>
      </div>
      <div class="nav-center">
        <div class="repo-input-container">
          <input
              type="text"
              v-model="repoPath"
              placeholder="输入 Git 仓库路径"
              @keyup.enter="loadRepo"
              class="repo-path-input"
          />
          <button @click="browseRepo" class="browse-btn" title="浏览目录">
            📁
          </button>
        </div>
        <div class="repo-actions">
          <button @click="loadRepo" class="primary-btn">
            <span class="btn-icon">📂</span>
            加载仓库
          </button>
          <button @click="refreshData" class="secondary-btn">
            <span class="btn-icon">🔄</span>
            刷新
          </button>
          <button @click="pullChanges" class="secondary-btn">
            <span class="btn-icon">⬇️</span>
            拉取
          </button>
          <button @click="pushChanges" class="secondary-btn">
            <span class="btn-icon">⬆️</span>
            推送
          </button>
        </div>
      </div>
      <div class="nav-right">
        <div class="connection-status">
          <span
              :class="['status-indicator', {
              'connected': repoPath,
              'disconnected': !repoPath
            }]"
          ></span>
          <span class="status-text">{{ repoPath ? '已连接' : '未连接' }}</span>
        </div>
      </div>
    </header>

    <!-- 主内容区域 -->
    <div class="main-content">
      <!-- 左侧边栏 - 分支管理 -->
      <aside class="sidebar left-sidebar">
        <div class="panel-header">
          <h2 class="panel-title">
            <span class="branch-icon">🌱</span>
            分支管理
          </h2>
          <div class="panel-actions">
            <button @click="createBranch" class="icon-btn" title="创建分支">
              ➕
            </button>
            <button @click="refreshBranches" class="icon-btn" title="刷新">
              🔄
            </button>
          </div>
        </div>

        <div class="panel-content">
          <!-- 分支搜索 -->
          <div class="search-box">
            <input
                type="text"
                v-model="branchFilter"
                placeholder="搜索分支..."
                class="search-input"
            />
            <span class="search-icon">🔍</span>
          </div>

          <div class="branch-section">
            <div class="section-header">
              <h3 class="section-title">本地分支</h3>
              <span class="item-count">({{ (localBranches.length || 0) }})</span>
            </div>
            <div class="branch-list">
              <div v-for="branch in localBranches"
                  :key="'local-' + branch.name"
                  :class="[
                  'branch-item',
                  { 'active': branch.current, 'current': branch.current }
                ]"
                  @dblclick="showBranchHistory(branch.name)"
                  @contextmenu.prevent="openBranchContextMenu($event, branch, 'local')"
              >
                <div class="branch-info">
                  <span class="branch-type">🌿</span>
                  <span class="branch-name" :title="branch.name">{{ branch.name }}</span>
                  <span v-if="branch.current" class="branch-current-badge" title="当前分支">●</span>
                </div>
                <div class="branch-actions">
                  <button
                      v-if="!branch.current"
                      @click.stop="switchBranch(branch.name)"
                      class="action-btn switch-btn"
                      title="切换到此分支"
                  >
                    ↔️
                  </button>
                  <button
                      v-if="!branch.current"
                      @click.stop="deleteBranch(branch.name)"
                      class="action-btn delete-btn"
                      title="删除分支"
                  >
                    ❌
                  </button>
                </div>
              </div>
              <div v-if="filteredLocalBranches?.length === 0 && branchFilter === ''" class="no-branches">
                <span class="empty-message">无本地分支</span>
              </div>
            </div>
            <div v-if="(filteredLocalBranches?.value?.length || 0) === 0 && branchFilter?.value === ''"
                 class="empty-branches">
              <span class="empty-message">无本地分支</span>
            </div>
          </div>

          <div class="branch-section">
            <div class="section-header">
              <h3 class="section-title">远程分支</h3>
              <span class="item-count">({{ (filteredRemoteBranches?.value?.length || 0) }})</span>
            </div>
            <div class="branch-list">
              <div
                  v-for="branch in filteredRemoteBranches"
                  :key="'remote-' + branch.name"
                  :class="['branch-item', { 'active': branch.current }]"
                  @dblclick="showBranchHistory(branch.name.replace('origin/', ''))"
                  @contextmenu.prevent="openBranchContextMenu($event, branch, 'remote')"
              >
                <div class="branch-info">
                  <span class="branch-type">📡</span>
                  <span class="branch-name" :title="branch.name">{{ branch.name }}</span>
                  <span v-if="branch.current" class="branch-current-badge" title="当前分支">●</span>
                </div>
                <div class="branch-actions">
                  <button
                      v-if="!branch.current"
                      @click.stop="switchBranch(branch.name)"
                      class="action-btn switch-btn"
                      title="切换到此分支"
                  >
                    ↔️
                  </button>
                </div>
              </div>
              <div v-if="(filteredRemoteBranches?.value?.length || 0) === 0 && branchFilter?.value === ''"
                   class="empty-branches">
                <span class="empty-message">无远程分支</span>
              </div>
            </div>
          </div>
        </div>
      </aside>

      <!-- 中间主内容区 -->
      <main class="main-area">
        <!-- 文件状态面板 -->
        <section class="panel file-status-panel">
          <div class="panel-header">
            <h2 class="panel-title">
              <span class="status-icon">📋</span>
              工作区状态
            </h2>
            <div class="panel-actions">
              <button @click="showStatus" class="icon-btn" title="查看详细状态">
                👁️
              </button>
              <button @click="stageAll" class="icon-btn" title="暂存全部">
                📥
              </button>
              <button @click="refreshStatus" class="icon-btn" title="刷新状态">
                🔄
              </button>
            </div>
          </div>

          <div class="panel-content">
            <div v-if="statusLoading" class="loading-state">
              <div class="spinner"></div>
              <span>加载状态中...</span>
            </div>
            <div v-else-if="!repoPath" class="empty-state">
              <span class="empty-icon">📁</span>
              <p>请先加载仓库以查看工作区状态</p>
            </div>
            <div v-else-if="(workingFiles?.value?.length || 0) === 0 && (stagedFiles?.value?.length || 0) === 0"
                 class="empty-state">
              <span class="empty-icon">✅</span>
              <p>工作区干净，无待提交更改</p>
            </div>
            <div v-else class="status-content">
              <!-- 未暂存文件 -->
              <div v-if="(workingFiles?.length || 0) > 0" class="status-section">
                <h3 class="status-section-title">
                  <span class="file-change-icon modified">●</span>
                  修改的文件
                  <span class="item-count">({{ workingFiles?.length || 0 }})</span>
                </h3>
                <div class="file-list">
                  <div
                      v-for="file in workingFiles"
                      :key="'working-' + file.path"
                      class="file-item"
                  >
                    <span class="file-status modified">●</span>
                    <span class="file-path">{{ file.path }}</span>
                    <div class="file-actions">
                      <button @click="stageFile(file.path)" class="small-btn primary">
                        暂存
                      </button>
                      <button @click="discardChanges(file.path)" class="small-btn danger">
                        丢弃
                      </button>
                    </div>
                  </div>
                </div>
              </div>

              <!-- 已暂存文件 -->
              <div v-if="(stagedFiles?.length || 0) > 0" class="status-section">
                <h3 class="status-section-title">
                  <span class="file-change-icon staged">✓</span>
                  已暂存文件
                  <span class="item-count">({{ stagedFiles?.length || 0 }})</span>
                </h3>
                <div class="file-list">
                  <div
                      v-for="file in stagedFiles"
                      :key="'staged-' + file.path"
                      class="file-item"
                  >
                    <span class="file-status staged">✓</span>
                    <span class="file-path">{{ file.path }}</span>
                    <div class="file-actions">
                      <button @click="unstageFile(file.path)" class="small-btn secondary">
                        取消暂存
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>

        <!-- 提交面板 -->
        <section class="panel commit-panel">
          <div class="panel-header">
            <h2 class="panel-title">
              <span class="commit-icon">✍️</span>
              创建提交
            </h2>
          </div>

          <div class="panel-content">
            <div class="commit-form">
              <textarea
                  v-model="commitMessage"
                  placeholder="输入提交信息..."
                  class="commit-message-input"
                  rows="3"
              ></textarea>
              <div class="commit-actions">
                <button
                    @click="commitChanges"
                    :disabled="!canCommit"
                    class="primary-btn commit-btn"
                >
                  <span class="btn-icon">💾</span>
                  提交更改 ({{ (stagedFiles?.value?.length || 0) }} 个文件)
                </button>
              </div>
            </div>
          </div>
        </section>
      </main>

      <!-- 右侧边栏 - 提交历史 -->
      <aside class="sidebar right-sidebar">
        <div class="panel-header">
          <h2 class="panel-title">
            <span class="history-icon">📜</span>
            提交历史
          </h2>
        </div>

        <div class="panel-content">
          <div v-if="commitsLoading" class="loading-state">
            <div class="spinner"></div>
            <span>加载提交历史...</span>
          </div>
          <div v-else-if="!repoPath" class="empty-state">
            <span class="empty-icon">📁</span>
            <p>请先加载仓库以查看提交历史</p>
          </div>
          <div v-else-if="(commits?.length || 0) === 0" class="empty-state">
            <span class="empty-icon">📝</span>
            <p>暂无提交历史</p>
          </div>
          <div v-else-if="commits && Array.isArray(commits) && commits.length > 0" class="commits-container">
            <div
                v-for="(commit, index) in commits"
                :key="commit?.hash || index.toString()"
                class="commit-item"
                @click="selectCommit(commit)"
                :class="{ 'selected': selectedCommit && selectedCommit.hash === commit.hash }"
            >
              <div class="commit-overview">
                <div class="commit-graph">
                  <div class="commit-dot" :style="getCommitColor(index)"></div>
                  <div class="commit-line" :style="getCommitLineColor(index)"></div>
                </div>
                <div class="commit-main">
                  <div class="commit-hash" :title="commit?.hash">{{ commit?.hash?.substring(0, 8) || 'N/A' }}</div>
                  <div class="commit-message" :title="commit?.message">{{ truncateText(commit?.message, 60) }}</div>
                </div>
                <div class="commit-meta">
                  <div class="commit-author" :title="commit?.author">{{
                      (commit?.author || '').split('<')[0].trim()
                    }}
                  </div>
                  <div class="commit-date" :title="commit?.date">{{ formatDate(commit?.date) }}</div>
                </div>
              </div>
              <div class="commit-refs" v-if="commit?.branches && commit.branches.length > 0">
                <span
                    v-for="branch in commit.branches"
                    :key="branch"
                    class="commit-ref"
                >
                  {{ branch }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </aside>
    </div>

    <!-- 底部状态栏 -->
    <footer class="status-bar">
      <div class="status-left">
        <span class="current-branch">
          🌿 {{ currentBranch || '未加载仓库' }}
        </span>
      </div>
      <div class="status-center">
        <span class="repo-path" v-if="repoPath" :title="repoPath">
          {{ repoPath }}
        </span>
      </div>
      <div class="status-right">
        <span class="change-summary">
          <span v-if="(workingFiles?.length || 0) > 0" class="working-changes">
            🔴 {{ workingFiles?.length || 0 }} 个修改
          </span>
          <span v-if="(stagedFiles?.length || 0) > 0" class="staged-changes">
            🟢 {{ stagedFiles?.length || 0 }} 个暂存
          </span>
        </span>
      </div>
    </footer>

    <!-- 通知组件 -->
    <transition name="slide-fade">
      <div
          v-if="notification.visible"
          :class="['notification', notification.type]"
      >
        {{ notification.message }}
      </div>
    </transition>
  </div>
</template>

<script>
import {computed, onMounted, reactive, ref} from 'vue'

// 导入Wails运行时和Go模块
// 注意：在生产环境中，Wails会在运行时注入这些对象，所以不需要显式导入
// 我们将在运行时通过window.go访问这些函数

export default {
  name: 'App',
  setup() {
    // State variables
    const repoPath = ref('D:/workspace/go-git-client-window')
    const currentBranch = ref('')
    const allBranches = ref([])
    const commits = ref([])
    const branchesLoading = ref(false)
    const commitsLoading = ref(false)
    const statusLoading = ref(false)

    // 新增状态变量
    const workingFiles = ref([])  // 未暂存的文件
    const stagedFiles = ref([])   // 已暂存的文件
    const commitMessage = ref('') // 提交信息
    const branchFilter = ref('')  // 分支过滤器
    const selectedCommit = ref(null) // 当前选中的提交

    // 用于存储原始文件列表，以便进行暂存/取消暂存操作
    const originalStatus = ref('')

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

    // 计算属性：是否可以提交
    const canCommit = computed(() => {
      return stagedFiles.value.length > 0 && commitMessage.value.trim()
    })

    // 计算属性：过滤后的本地分支
    const filteredLocalBranches = computed(() => {
      if (!localBranches.value || !Array.isArray(localBranches.value)) {
        return []
      }
      if (!branchFilter.value) {
        return localBranches.value
      }
      return localBranches.value.filter(branch =>
          branch.name.toLowerCase().includes(branchFilter.value.toLowerCase())
      )
    })

    // 计算属性：过滤后的远程分支
    const filteredRemoteBranches = computed(() => {
      if (!remoteBranches.value || !Array.isArray(remoteBranches.value)) {
        return []
      }
      if (!branchFilter.value) {
        return remoteBranches.value
      }
      return remoteBranches.value.filter(branch =>
          branch.name.toLowerCase().includes(branchFilter.value.toLowerCase())
      )
    })

    // 截断文本以适应显示
    const truncateText = (text, maxLength) => {
      if (!text) return ''
      return text.length > maxLength ? text.substring(0, maxLength) + '...' : text
    }

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

    // 格式化日期显示
    const formatDate = (dateString) => {
      try {
        const date = new Date(dateString)
        return date.toLocaleDateString('zh-CN', {
          year: 'numeric',
          month: 'short',
          day: 'numeric',
          hour: '2-digit',
          minute: '2-digit'
        })
      } catch {
        return dateString
      }
    }

    // 浏览目录
    const browseRepo = async () => {
      try {
        // 尝试调用 Wails 的目录选择功能
        if (window.go && window.go.main && window.go.main.App && window.go.main.App.SelectDirectory) {
          const selectedPath = await window.go.main.App.SelectDirectory()
          if (selectedPath) {
            repoPath.value = selectedPath
            await loadRepo()
          }
        } else {
          // 如果没有 Wails 目录选择功能，则使用提示框
          const path = prompt('请输入仓库路径:', repoPath.value)
          if (path) {
            repoPath.value = path
            await loadRepo()
          }
        }
      } catch (error) {
        showNotification(`浏览目录失败: ${error}`, 'error')
      }
    }

    const loadRepo = async () => {
      if (!repoPath.value.trim()) {
        showNotification('请输入仓库路径', 'error')
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
          loadCurrentBranch(),
          loadStatus()
        ])
      } catch (error) {
        console.error('刷新数据错误:', error)
      }
    }

    const refreshStatus = async () => {
      if (!repoPath.value) return
      await loadStatus()
    }

    const loadBranches = async () => {
      if (!repoPath.value) return
      branchesLoading.value = true
      try {
        const result = await window.go.main.App.GitBranch(repoPath.value)
        allBranches.value = JSON.parse(result);
      } catch (error) {
        showNotification(`加载分支失败: ${error}`, 'error')
      } finally {
        branchesLoading.value = false
      }
    }

    const loadCommits = async () => {
      if (!repoPath.value) return

      commitsLoading.value = true

      try {
        const result = await window.go.main.App.GitLog(repoPath.value, 50) // 减少数量以提高性能
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
        showNotification(`加载提交历史失败: ${error}`, 'error')
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
        console.error('获取当前分支失败:', error)
      }
    }

    // 加载工作区状态
    const loadStatus = async () => {
      if (!repoPath.value) return

      statusLoading.value = true
      try {
        const result = await window.go.main.App.GitStatus(repoPath.value)
        originalStatus.value = result
        parseGitStatus(result)
      } catch (error) {
        showNotification(`加载状态失败: ${error}`, 'error')
      } finally {
        statusLoading.value = false
      }
    }

    // 解析 Git 状态输出
    const parseGitStatus = (statusOutput) => {
      const lines = statusOutput.split('\\n').filter(line => line.trim())

      // 重置数组
      workingFiles.value = []
      stagedFiles.value = []

      // 解析 Git 状态输出，处理各种状态标记
      lines.forEach(line => {
        // 去除前导空格并解析状态
        const trimmedLine = line.trim()
        if (trimmedLine.length < 2) return

        // 第一个字符表示暂存区状态，第二个字符表示工作区状态
        const stagedStatus = trimmedLine.charAt(0)
        const workingStatus = trimmedLine.charAt(1)

        // 获取文件路径（跳过状态字符）
        let filePath = trimmedLine.slice(3).trim() // 跳过状态字符和空格

        // 处理重命名文件格式: "R  old_file -> new_file"
        if (stagedStatus === 'R' || workingStatus === 'R') {
          const arrowIndex = trimmedLine.indexOf(' -> ')
          if (arrowIndex !== -1) {
            filePath = trimmedLine.substring(arrowIndex + 4).trim()
          }
        }

        if (filePath) {
          // 如果暂存区有变化（非空格）
          if (stagedStatus !== ' ' && stagedStatus !== '?') {
            stagedFiles.value.push({
              path: filePath,
              status: stagedStatus
            })
          }

          // 如果工作区有变化（非空格或非A/U等合并状态）
          if (workingStatus !== ' ' && workingStatus !== '?' &&
              workingStatus !== 'A' && workingStatus !== 'U') { // A/U通常表示已添加到暂存区
            workingFiles.value.push({
              path: filePath,
              status: workingStatus
            })
          }

          // 特殊处理：如果文件是新增且未被暂存
          if (stagedStatus === ' ' && workingStatus === '?') {
            workingFiles.value.push({
              path: filePath,
              status: '?'
            })
          }
        }
      })
    }

    // 暂存单个文件
    const stageFile = async (filePath) => {
      if (!repoPath.value) return

      try {
        await window.go.main.App.GitAdd(repoPath.value, filePath)
        await loadStatus()
        showNotification(`已暂存文件: ${filePath}`, 'success')
      } catch (error) {
        showNotification(`暂存文件失败: ${error}`, 'error')
      }
    }

    // 取消暂存单个文件
    const unstageFile = async (filePath) => {
      if (!repoPath.value) return

      try {
        await window.go.main.App.GitReset(repoPath.value, filePath)
        await loadStatus()
        showNotification(`已取消暂存: ${filePath}`, 'success')
      } catch (error) {
        showNotification(`取消暂存失败: ${error}`, 'error')
      }
    }

    // 暂存所有文件
    const stageAll = async () => {
      if (!repoPath.value) return

      try {
        await window.go.main.App.GitAddAll(repoPath.value)
        await loadStatus()
        showNotification('已暂存所有文件', 'success')
      } catch (error) {
        showNotification(`暂存所有文件失败: ${error}`, 'error')
      }
    }

    // 丢弃文件更改
    const discardChanges = async (filePath) => {
      if (!repoPath.value) return

      if (!confirm(`确定要丢弃文件 "${filePath}" 的更改吗？此操作不可撤销！`)) {
        return
      }

      try {
        await window.go.main.App.GitCheckoutFile(repoPath.value, filePath)
        await loadStatus()
        showNotification(`已丢弃文件更改: ${filePath}`, 'success')
      } catch (error) {
        showNotification(`丢弃更改失败: ${error}`, 'error')
      }
    }

    // 执行提交
    const commitChanges = async () => {
      if (!repoPath.value) return

      if (!canCommit.value) {
        showNotification('请先暂存文件并输入提交信息', 'error')
        return
      }

      try {
        await window.go.main.App.GitCommit(repoPath.value, commitMessage.value)
        commitMessage.value = ''
        await refreshData()
        showNotification('提交成功', 'success')
      } catch (error) {
        showNotification(`提交失败: ${error}`, 'error')
      }
    }

    const switchBranch = async (branchName) => {
      if (!repoPath.value) return
      if (!confirm(`确定切换到分支 "${branchName}" 吗？`)) return

      try {
        await window.go.main.App.GitCheckout(repoPath.value, branchName)
        await refreshData()
        showNotification(`已切换到分支 ${branchName}`, 'success')
      } catch (error) {
        showNotification(`切换分支失败: ${error}`, 'error')
      }
    }

    const createBranch = async () => {
      if (!repoPath.value) {
        showNotification('请先加载仓库', 'error')
        return
      }

      const branchName = prompt('请输入新分支名称:')
      if (!branchName) return

      try {
        await window.go.main.App.GitCreateBranch(repoPath.value, branchName)
        await refreshData()
        showNotification(`已创建并切换到分支 ${branchName}`, 'success')
      } catch (error) {
        showNotification(`创建分支失败: ${error}`, 'error')
      }
    }

    const showStatus = async () => {
      if (!repoPath.value) {
        showNotification('请先加载仓库', 'error')
        return
      }

      try {
        const result = await window.go.main.App.GitStatus(repoPath.value)
        alert('Git 状态:\n\n' + result)
      } catch (error) {
        showNotification(`获取状态失败: ${error}`, 'error')
      }
    }

    const refreshBranches = async () => {
      await loadBranches()
    }

    const refreshCommits = async () => {
      await loadCommits()
    }

    // 右键点击分支打开上下文菜单
    const openBranchContextMenu = (event, branch, type) => {
      // 创建右键菜单
      const menu = document.createElement('div');
      menu.className = 'context-menu';
      menu.style.position = 'fixed';
      menu.style.left = event.clientX + 'px';
      menu.style.top = event.clientY + 'px';
      menu.style.zIndex = '1000';
      menu.style.backgroundColor = '#3c3f41';
      menu.style.border = '1px solid #4e5254';
      menu.style.borderRadius = '4px';
      menu.style.padding = '4px 0';
      menu.style.boxShadow = '0 4px 6px rgba(0, 0, 0, 0.1)';
      menu.style.minWidth = '120px';
      
      // 清除任何现有的右键菜单
      document.querySelectorAll('.context-menu').forEach(el => el.remove());
      
      // 对于所有分支，都提供查看历史选项
      const historyItem = document.createElement('div');
      historyItem.className = 'context-menu-item';
      historyItem.innerHTML = '📜 查看提交历史';
      historyItem.style.padding = '8px 12px';
      historyItem.style.cursor = 'pointer';
      historyItem.style.color = '#a9b7c6';
      historyItem.style.fontSize = '12px';
      historyItem.onmouseover = () => historyItem.style.backgroundColor = '#4e5254';
      historyItem.onmouseout = () => historyItem.style.backgroundColor = 'transparent';
      historyItem.onclick = () => {
        showBranchHistory(branch.name);
        document.body.removeChild(menu);
      };
      menu.appendChild(historyItem);
      
      // 添加分隔线
      const separator = document.createElement('hr');
      separator.style.margin = '4px 0';
      separator.style.borderColor = '#4e5254';
      separator.style.borderStyle = 'solid';
      menu.appendChild(separator);
      
      // 非当前分支的操作
      if (!branch.current) {
        // 切换分支选项
        const switchItem = document.createElement('div');
        switchItem.className = 'context-menu-item';
        switchItem.innerHTML = '🔄 切换到此分支';
        switchItem.style.padding = '8px 12px';
        switchItem.style.cursor = 'pointer';
        switchItem.style.color = '#a9b7c6';
        switchItem.style.fontSize = '12px';
        switchItem.onmouseover = () => switchItem.style.backgroundColor = '#4e5254';
        switchItem.onmouseout = () => switchItem.style.backgroundColor = 'transparent';
        switchItem.onclick = () => {
          switchBranch(branch.name);
          document.body.removeChild(menu);
        };
        menu.appendChild(switchItem);
        
        // 删除分支选项
        const deleteItem = document.createElement('div');
        deleteItem.className = 'context-menu-item';
        deleteItem.innerHTML = '❌ 删除分支';
        deleteItem.style.padding = '8px 12px';
        deleteItem.style.cursor = 'pointer';
        deleteItem.style.color = '#a9b7c6';
        deleteItem.style.fontSize = '12px';
        deleteItem.onmouseover = () => deleteItem.style.backgroundColor = '#4e5254';
        deleteItem.onmouseout = () => deleteItem.style.backgroundColor = 'transparent';
        deleteItem.onclick = () => {
          deleteBranch(branch.name);
          document.body.removeChild(menu);
        };
        menu.appendChild(deleteItem);
      } else {
        // 当前分支的操作
        const mergeItem = document.createElement('div');
        mergeItem.className = 'context-menu-item';
        mergeItem.innerHTML = '🔀 合并其他分支';
        mergeItem.style.padding = '8px 12px';
        mergeItem.style.cursor = 'pointer';
        mergeItem.style.color = '#a9b7c6';
        mergeItem.style.fontSize = '12px';
        mergeItem.onmouseover = () => mergeItem.style.backgroundColor = '#4e5254';
        mergeItem.onmouseout = () => mergeItem.style.backgroundColor = 'transparent';
        mergeItem.onclick = () => {
          alert('合并功能将在后续版本中实现');
          document.body.removeChild(menu);
        };
        menu.appendChild(mergeItem);
      }
      
      // 如果是远程分支，添加拉取到本地选项
      if (type === 'remote') {
        const pullItem = document.createElement('div');
        pullItem.className = 'context-menu-item';
        pullItem.innerHTML = '📥 拉取到本地';
        pullItem.style.padding = '8px 12px';
        pullItem.style.cursor = 'pointer';
        pullItem.style.color = '#a9b7c6';
        pullItem.style.fontSize = '12px';
        pullItem.onmouseover = () => pullItem.style.backgroundColor = '#4e5254';
        pullItem.onmouseout = () => pullItem.style.backgroundColor = 'transparent';
        pullItem.onclick = () => {
          // 创建本地同名分支并跟踪远程分支
          createBranchFromRemote(branch.name);
          document.body.removeChild(menu);
        };
        menu.appendChild(pullItem);
      }
      
      // 添加到页面
      document.body.appendChild(menu);
      
      // 点击其他地方关闭菜单
      const closeMenu = (e) => {
        if (!menu.contains(e.target)) {
          if (document.body.contains(menu)) {
            document.body.removeChild(menu);
          }
          document.removeEventListener('click', closeMenu);
        }
      };
      setTimeout(() => {
        document.addEventListener('click', closeMenu);
      }, 100);
    }

    // 创建本地分支跟踪远程分支
    const createBranchFromRemote = async (remoteBranchName) => {
      if (!repoPath.value) {
        showNotification('请先加载仓库', 'error');
        return;
      }
      
      // 从远程分支名中提取本地分支名 (例如，从 'origin/main' 提取 'main')
      const localBranchName = remoteBranchName.replace('origin/', '').replace('remotes/', '');
      
      if (!confirm(`确定要从远程分支 "${remoteBranchName}" 创建本地分支 "${localBranchName}" 吗？`)) {
        return;
      }
      
      try {
        // 检查本地分支是否已存在
        const localExists = localBranches.value.some(branch => branch.name === localBranchName);
        if (localExists) {
          if (!confirm(`本地分支 "${localBranchName}" 已存在，是否切换到该分支？`)) {
            return;
          }
          await switchBranch(localBranchName);
          return;
        }
        
        // 创建本地分支并跟踪远程分支
        await window.go.main.App.GitCreateBranch(repoPath.value, localBranchName);
        await window.go.main.App.GitCheckout(repoPath.value, localBranchName);
        
        // 拉取远程分支的最新内容
        await window.go.main.App.GitPull(repoPath.value, localBranchName);
        
        await refreshData();
        showNotification(`已创建本地分支 "${localBranchName}" 并切换到该分支`, 'success');
      } catch (error) {
        showNotification(`创建分支失败: ${error}`, 'error');
      }
    };
    
    // 删除分支
    const deleteBranch = async (branchName) => {
      if (!repoPath.value) return

      if (!confirm(`确定要删除分支 "${branchName}" 吗？此操作不可撤销！`)) {
        return
      }

      try {
        await window.go.main.App.GitDeleteBranch(repoPath.value, branchName)
        await refreshData()
        showNotification(`已删除分支: ${branchName}`, 'success')
      } catch (error) {
        showNotification(`删除分支失败: ${error}`, 'error')
      }
    }

    // 选择提交
    const selectCommit = (commit) => {
      selectedCommit.value = commit
      console.log('Selected commit:', commit)
    }

    // 获取提交颜色
    const getCommitColor = (index) => {
      // 基于索引生成不同的颜色
      const hue = (index * 137.5) % 360 // 使用黄金角度生成颜色差异
      return {
        backgroundColor: `hsl(${hue}, 70%, 60%)`
      }
    }

    // 获取提交线条颜色
    const getCommitLineColor = (index) => {
      const hue = (index * 137.5) % 360
      return {
        borderColor: `hsl(${hue}, 70%, 60%)`
      }
    }
    
    // 显示指定分支的提交历史
    const showBranchHistory = async (branchName) => {
      if (!repoPath.value) {
        showNotification('请先加载仓库', 'error');
        return;
      }
      
      try {
        commitsLoading.value = true;
        // 获取指定分支的提交历史
        const result = await window.go.main.App.GitBranchLog(repoPath.value, branchName, 50);
        commits.value = JSON.parse(result);
        showNotification(`已加载分支 "${branchName}" 的提交历史`, 'info');
      } catch (error) {
        showNotification(`加载分支 "${branchName}" 历史失败: ${error}`, 'error');
      } finally {
        commitsLoading.value = false;
      }
    }

    // 拉取远程更改
    const pullChanges = async () => {
      if (!repoPath.value) {
        showNotification('请先加载仓库', 'error')
        return
      }

      try {
        const result = await window.go.main.App.GitPull(repoPath.value)
        await refreshData()
        showNotification(`拉取成功: ${result || '无新更改'}`, 'success')
      } catch (error) {
        showNotification(`拉取失败: ${error}`, 'error')
      }
    }

    // 推送到远程仓库
    const pushChanges = async () => {
      if (!repoPath.value) {
        showNotification('请先加载仓库', 'error')
        return
      }

      try {
        const result = await window.go.main.App.GitPush(repoPath.value)
        showNotification(`推送成功: ${result || '已同步'}`, 'success')
      } catch (error) {
        showNotification(`推送失败: ${error}`, 'error')
      }
    }

    // 页面加载时初始化
    onMounted(async () => {
      // 尝试加载默认仓库
      if (repoPath.value) {
        await loadRepo()
      }
    })

    return {
      repoPath,
      currentBranch,
      branches: allBranches,
      commits,
      branchesLoading,
      commitsLoading,
      statusLoading,
      localBranches,
      remoteBranches,
      workingFiles,
      stagedFiles,
      commitMessage,
      canCommit,
      selectedCommit,
      notification,
      loadRepo,
      refreshData,
      loadBranches,
      loadCommits,
      loadCurrentBranch,
      loadStatus,
      refreshStatus,
      switchBranch,
      createBranch,
      createBranchFromRemote,
      showStatus,
      refreshBranches,
      refreshCommits,
      stageFile,
      unstageFile,
      stageAll,
      discardChanges,
      commitChanges,
      browseRepo,
      pullChanges,
      pushChanges,
      selectCommit,
      showBranchHistory,
      openBranchContextMenu,
      getCommitColor,
      getCommitLineColor,
      truncateText,
      formatDate,
      showNotification
    }
  }
}
</script>

<style src="./styles/app.css"></style>