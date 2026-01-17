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
              <span class="item-count">({{ localBranchesCount.value }})</span>
            </div>
            <div class="branch-list">
              <div
                v-for="branch in filteredLocalBranches"
                :key="'local-' + branch.name"
                :class="[
                  'branch-item',
                  { 'active': branch.current, 'current': branch.current }
                ]"
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
            </div>
          </div>
          
          <div class="branch-section">
            <div class="section-header">
              <h3 class="section-title">远程分支</h3>
              <span class="item-count">({{ remoteBranchesCount.value }})</span>
            </div>
            <div class="branch-list">
              <div
                v-for="branch in filteredRemoteBranches"
                :key="'remote-' + branch.name"
                :class="['branch-item', { 'active': branch.current }]"
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
            <div v-else-if="(workingFiles?.length || 0) === 0 && (stagedFiles?.length || 0) === 0" class="empty-state">
              <span class="empty-icon">✅</span>
              <p>工作区干净，无待提交更改</p>
            </div>
            <div v-else class="status-content">
              <!-- 未暂存文件 -->
              <div v-if="workingFiles.length > 0" class="status-section">
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
              <div v-if="stagedFiles.length > 0" class="status-section">
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
                  提交更改 ({{ stagedFiles?.length || 0 }} 个文件)
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
          <div class="panel-actions">
            <button @click="refreshCommits" class="icon-btn" title="刷新">
              🔄
            </button>
          </div>
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
          <div v-else-if="commits.length === 0" class="empty-state">
            <span class="empty-icon">📝</span>
            <p>暂无提交历史</p>
          </div>
          <div v-else class="commits-container">
            <div 
              v-for="(commit, index) in commits" 
              :key="commit.hash" 
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
                  <div class="commit-hash" :title="commit.hash">{{ commit.hash.substring(0, 8) }}</div>
                  <div class="commit-message" :title="commit.message">{{ truncateText(commit.message, 60) }}</div>
                </div>
                <div class="commit-meta">
                  <div class="commit-author" :title="commit.author">{{ commit.author.split('<')[0].trim() }}</div>
                  <div class="commit-date" :title="commit.date">{{ formatDate(commit.date) }}</div>
                </div>
              </div>
              <div class="commit-refs" v-if="commit.branches && commit.branches.length > 0">
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
          <span v-if="workingFiles.length > 0" class="working-changes">
            🔴 {{ workingFiles.length }} 个修改
          </span>
          <span v-if="stagedFiles.length > 0" class="staged-changes">
            🟢 {{ stagedFiles.length }} 个暂存
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
import {computed, reactive, ref, onMounted} from 'vue'

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
    
    // 安全的长度计算
    const localBranchesCount = computed(() => {
      return filteredLocalBranches.value?.length || 0
    })
    
    const remoteBranchesCount = computed(() => {
      return filteredRemoteBranches.value?.length || 0
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
      // 在实际应用中，这里会打开一个上下文菜单
      console.log(`右键点击分支 ${branch.name} (${type})`)
    }
    
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
      localBranchesCount,
      remoteBranchesCount,
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
      getCommitColor,
      getCommitLineColor,
      truncateText,
      formatDate,
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

/* 搜索框样式 */
.search-box {
  position: relative;
  margin-bottom: 12px;
  padding: 0 8px;
}

.search-input {
  width: 100%;
  padding: 8px 30px 8px 12px;
  background-color: #3c3f41;
  border: 1px solid #5e6366;
  border-radius: 4px;
  color: #a9b7c6;
  font-size: 13px;
}

.search-input:focus {
  outline: none;
  border-color: #4a6d8c;
}

.search-icon {
  position: absolute;
  right: 10px;
  top: 50%;
  transform: translateY(-50%);
  color: #808080;
}

/* 面板通用样式 */
.panel {
  background-color: #313335;
  border-radius: 6px;
  overflow: hidden;
  margin-bottom: 16px;
  border: 1px solid #4e5254;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background-color: #3c3f41;
  border-bottom: 1px solid #4e5254;
}

.panel-title {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: #a9b7c6;
}

.panel-actions {
  display: flex;
  gap: 6px;
}

.icon-btn {
  background-color: transparent;
  border: 1px solid #5e6366;
  color: #a9b7c6;
  width: 28px;
  height: 28px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.icon-btn:hover {
  background-color: #4e5254;
  border-color: #4a6d8c;
  color: #4a6d8c;
}

.panel-content {
  padding: 12px;
  max-height: 400px;
  overflow-y: auto;
}

/* 加载状态 */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 30px;
  color: #808080;
  gap: 10px;
}

.spinner {
  width: 24px;
  height: 24px;
  border: 2px solid #4e5254;
  border-top: 2px solid #4a6d8c;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* 空状态 */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 30px;
  color: #808080;
  text-align: center;
}

.empty-icon {
  font-size: 24px;
  margin-bottom: 10px;
}

/* 分支项样式 */
.branch-section {
  margin-bottom: 20px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 8px 8px 8px;
  border-bottom: 1px solid #4e5254;
  margin-bottom: 8px;
}

.section-title {
  font-size: 13px;
  font-weight: 600;
  color: #6a8759;
  margin: 0;
}

.item-count {
  background-color: #4e5254;
  color: #a9b7c6;
  padding: 2px 6px;
  border-radius: 10px;
  font-size: 11px;
}

.branch-list {
  padding: 0 4px;
}

.branch-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 6px 8px;
  margin-bottom: 2px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 13px;
  transition: background-color 0.2s;
}

.branch-item:hover {
  background-color: #4e5254;
}

.branch-item.active {
  background-color: #4a6d8c;
}

.branch-info {
  display: flex;
  align-items: center;
  flex: 1;
  min-width: 0;
}

.branch-type {
  margin-right: 6px;
  font-size: 12px;
}

.branch-name {
  flex: 1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.branch-current-badge {
  background-color: #6a8759;
  color: #ffffff;
  padding: 1px 5px;
  border-radius: 3px;
  font-size: 10px;
  margin-left: 6px;
}

.branch-actions {
  display: flex;
  gap: 4px;
  margin-left: 8px;
}

.action-btn {
  background: none;
  border: 1px solid #5e6366;
  color: #a9b7c6;
  width: 22px;
  height: 22px;
  border-radius: 3px;
  cursor: pointer;
  font-size: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.action-btn:hover {
  background-color: #4e5254;
}

.switch-btn:hover {
  border-color: #4a6d8c;
  color: #4a6d8c;
}

.delete-btn:hover {
  border-color: #cc7832;
  color: #cc7832;
}

/* 文件状态样式 */
.status-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.status-section {
  margin-bottom: 16px;
}

.status-section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  font-weight: 600;
  color: #a9b7c6;
  margin: 0 0 8px 0;
  padding-bottom: 4px;
  border-bottom: 1px solid #4e5254;
}

.file-change-icon {
  font-size: 14px;
}

.file-change-icon.modified {
  color: #cc7832;
}

.file-change-icon.staged {
  color: #6a8759;
}

.file-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.file-item {
  display: flex;
  align-items: center;
  padding: 6px 8px;
  border-radius: 4px;
  background-color: #3c3f41;
  transition: background-color 0.2s;
}

.file-item:hover {
  background-color: #4e5254;
}

.file-status {
  width: 20px;
  text-align: center;
  font-weight: bold;
  margin-right: 6px;
}

.file-status.modified {
  color: #cc7832;
}

.file-status.staged {
  color: #6a8759;
}

.file-path {
  flex: 1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-right: 8px;
  font-family: monospace;
  font-size: 12px;
}

.file-actions {
  display: flex;
  gap: 4px;
}

.small-btn {
  padding: 4px 8px;
  border-radius: 3px;
  border: none;
  cursor: pointer;
  font-size: 11px;
  transition: background-color 0.2s;
}

.small-btn.primary {
  background-color: #4a6d8c;
  color: white;
}

.small-btn.primary:hover {
  background-color: #5c7d9e;
}

.small-btn.secondary {
  background-color: #4e5254;
  color: #a9b7c6;
}

.small-btn.secondary:hover {
  background-color: #5e6366;
}

.small-btn.danger {
  background-color: #cc7832;
  color: white;
}

.small-btn.danger:hover {
  background-color: #da8b45;
}

.commit-btn {
  background-color: #6a8759;
  color: white;
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 13px;
  display: flex;
  align-items: center;
  gap: 6px;
  transition: background-color 0.2s;
}

.commit-btn:disabled {
  background-color: #4e5254;
  cursor: not-allowed;
}

.commit-btn:not(:disabled):hover {
  background-color: #7baa6d;
}

/* 提交历史样式 */
.commits-container {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.commit-item {
  padding: 12px;
  border: 1px solid #4e5254;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
  background-color: #3c3f41;
}

.commit-item:hover {
  border-color: #4a6d8c;
  background-color: #474a4d;
}

.commit-item.selected {
  border-color: #6a8759;
  background-color: #424547;
  box-shadow: 0 0 0 2px rgba(106, 135, 89, 0.3);
}

.commit-overview {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.commit-main {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
}

.commit-hash {
  font-family: monospace;
  font-size: 12px;
  color: #9876aa;
  background-color: #2b2b2b;
  padding: 2px 6px;
  border-radius: 3px;
  min-width: 60px;
  text-align: center;
}

.commit-message {
  flex: 1;
  font-weight: 500;
  color: #a9b7c6;
  word-break: break-word;
}

.commit-meta {
  display: flex;
  justify-content: space-between;
  gap: 10px;
  font-size: 12px;
  color: #808080;
  margin-top: 4px;
}

.commit-author {
  color: #cc7832;
}

.commit-date {
  color: #808080;
}

.commit-refs {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
  margin-top: 6px;
}

.commit-ref {
  background-color: #4a6d8c;
  color: #ffffff;
  padding: 2px 6px;
  border-radius: 10px;
  font-size: 11px;
}

.commit-graph {
  display: flex;
  align-items: center;
  gap: 8px;
}

.commit-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  flex-shrink: 0;
}

.commit-line {
  height: 20px;
  width: 2px;
  border-left: 2px solid;
  margin-left: 5px;
}

/* 按钮样式 */
.btn-icon {
  margin-right: 6px;
}

.primary-btn {
  background-color: #4a6d8c;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 13px;
  display: flex;
  align-items: center;
  transition: background-color 0.2s;
}

.primary-btn:hover {
  background-color: #5c7d9e;
}

.secondary-btn {
  background-color: #3c3f41;
  color: #a9b7c6;
  border: 1px solid #5e6366;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 13px;
  display: flex;
  align-items: center;
  transition: all 0.2s;
}

.secondary-btn:hover {
  background-color: #4e5254;
  border-color: #4a6d8c;
}

/* 顶部导航栏样式 */
.top-nav {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background-color: #3c3f41;
  border-bottom: 1px solid #4e5254;
}

.nav-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.app-title {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #a9b7c6;
}

.git-icon {
  font-size: 20px;
}

.nav-center {
  display: flex;
  align-items: center;
  gap: 12px;
}

.repo-input-container {
  display: flex;
  align-items: center;
  background-color: #3c3f41;
  border: 1px solid #5e6366;
  border-radius: 4px;
  overflow: hidden;
}

.repo-path-input {
  padding: 8px 12px;
  background-color: #3c3f41;
  border: none;
  color: #a9b7c6;
  font-size: 13px;
  width: 300px;
}

.repo-path-input:focus {
  outline: none;
}

.browse-btn {
  padding: 8px 12px;
  background-color: #4a6d8c;
  border: none;
  color: #ffffff;
  cursor: pointer;
  font-size: 13px;
  transition: background-color 0.2s;
}

.browse-btn:hover {
  background-color: #5c7d9e;
}

.repo-actions {
  display: flex;
  gap: 8px;
}

.nav-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.connection-status {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #a9b7c6;
}

/* 主内容区域 */
.main-content {
  display: flex;
  flex: 1;
  overflow: hidden;
  gap: 16px;
  padding: 16px;
}

.sidebar {
  width: 320px;
  background-color: #313335;
  border-radius: 6px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.left-sidebar {
  border: 1px solid #4e5254;
}

.right-sidebar {
  border: 1px solid #4e5254;
}

.main-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 16px;
  overflow-y: auto;
}

.file-status-panel {
  flex: 1;
}

.commit-panel {
  margin-bottom: 0;
}

/* 状态栏 */
.status-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 16px;
  background-color: #3c3f41;
  border-top: 1px solid #4e5254;
  font-size: 12px;
}

.status-left, .status-center, .status-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.current-branch {
  color: #6a8759;
}

.repo-path {
  color: #a9b7c6;
  max-width: 400px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.change-summary {
  display: flex;
  gap: 12px;
}

.working-changes {
  color: #cc7832;
}

.staged-changes {
  color: #6a8759;
}

/* 动画效果 */
.slide-fade-enter-active {
  transition: all 0.3s ease;
}

.slide-fade-leave-active {
  transition: all 0.3s cubic-bezier(1, 0.5, 0.8, 1);
}

.slide-fade-enter-from, .slide-fade-leave-to {
  transform: translateX(20px);
  opacity: 0;
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