<template>
  <div class="app-container">
    <!-- 顶部导航栏 - 高度压缩，显示仓库路径、分支和操作按钮 -->
    <header class="top-nav">
      <div class="nav-left">
        <h1 class="app-title">
          <span class="git-icon">📊</span>
          Git 客户端
        </h1>
        <div class="repo-info">
          <div class="repo-path">{{ repoPath || '未选择仓库' }}</div>
          <div class="branch-selector">
            <select v-model="currentBranch" @change="switchBranch(currentBranch)" class="branch-dropdown">
              <option value="" disabled>选择分支...</option>
              <option 
                v-for="branch in allBranches" 
                :key="branch.name" 
                :value="branch.name"
                :selected="branch.current"
              >
                {{ branch.name }} {{ branch.current ? '(当前)' : '' }}
              </option>
            </select>
          </div>
        </div>
      </div>
      <div class="nav-right">
        <div class="nav-actions">
          <button @click="$emit('fetch-changes')" class="action-btn fetch-btn" title="获取最新更改">
            <span class="btn-icon">📥</span>
            Fetch
          </button>
          <button @click="$emit('pull-changes')" class="action-btn pull-btn" title="拉取更改">
            <span class="btn-icon">⬇️</span>
            Pull
            <span v-if="pullCount > 0" class="badge">{{ pullCount }}</span>
          </button>
          <button @click="$emit('push-changes')" class="action-btn push-btn" title="推送更改">
            <span class="btn-icon">⬆️</span>
            Push
          </button>
        </div>
      </div>
    </header>

    <!-- 主内容区域 - 三栏布局 -->
    <div class="main-content">
      <!-- 左侧边栏 - 导航面板（分支、标签等）-->
      <BranchList 
        :all-branches="allBranches"
        :branch-filter="branchFilter"
        @switch-branch="switchBranch"
        @create-branch="createBranch"
        @delete-branch="deleteBranch"
        @refresh-branches="loadBranches"
        @show-branch-history="showBranchHistory"
        @open-branch-context-menu="openBranchContextMenu"
        @update:branch-filter="branchFilter = $event"
      />
      
      <!-- 中间主内容区 - 提交历史图表 -->
      <main class="main-content-area">
        <div class="panel-header">
          <h2 class="panel-title">
            <span class="history-icon">📅</span>
            提交历史
          </h2>
        </div>
        <div class="panel-content">
          <div class="commit-graph-container">
            <div 
              v-for="(commit, index) in commits" 
              :key="commit.hash"
              :class="['commit-item', { 'selected': selectedCommit && selectedCommit.hash === commit.hash }]"
              @click="selectCommit(commit)"
            >
              <div class="commit-graph">
                <div class="commit-dot" :style="getCommitColor(index)"></div>
                <div class="commit-line" :style="getCommitLineColor(index)"></div>
              </div>
              <div class="commit-details">
                <div class="commit-main">
                  <div class="commit-message">{{ commit.message }}</div>
                  <div class="commit-hash">{{ commit.hash.substring(0, 7) }}</div>
                </div>
                <div class="commit-meta">
                  <div class="commit-author">{{ commit.author }}</div>
                  <div class="commit-date">{{ formatDate(commit.date) }}</div>
                </div>
              </div>
            </div>
            <div v-if="commits.length === 0" class="empty-state">
              <div class="empty-icon">📦</div>
              <p>暂无提交记录</p>
            </div>
          </div>
        </div>
      </main>

      <!-- 右侧操作面板 - 暂存和提交 -->
      <aside class="sidebar right-sidebar">
        <div class="panel-header">
          <h2 class="panel-title">
            <span class="staging-icon">📋</span>
            <span v-if="selectedCommit">提交详情</span>
            <span v-else>暂存 & 提交</span>
          </h2>
        </div>
        <div class="panel-content">
          <!-- 当选中提交时显示详情 -->
          <div v-if="selectedCommit" class="commit-detail-view">
            <div class="detail-header">
              <div class="detail-hash">{{ selectedCommit.hash }}</div>
              <div class="detail-message">{{ selectedCommit.message }}</div>
            </div>
            <div class="detail-meta">
              <div class="detail-author">作者: {{ selectedCommit.author }}</div>
              <div class="detail-date">日期: {{ formatDate(selectedCommit.date) }}</div>
            </div>
            <div class="diff-preview">
              <h4>文件变更预览</h4>
              <div class="diff-placeholder">
                <!-- 此处将显示文件差异预览 -->
                <p>变更文件列表将在后续版本中实现</p>
              </div>
            </div>
          </div>

          <!-- 当没有选中提交且有工作区更改时显示暂存/提交视图 -->
          <div v-else-if="workingFiles.length > 0 || stagedFiles.length > 0" class="staging-view">
            <!-- 文件变更列表 -->
            <div class="file-changes-section">
              <h4 class="section-subtitle">
                <span class="file-change-icon modified">📝</span>
                文件变更 ({{ workingFiles.length }})
              </h4>
              <div class="file-list">
                <div 
                  v-for="file in workingFiles" 
                  :key="'working-' + file.path"
                  class="file-item"
                >
                  <input 
                    type="checkbox" 
                    @change="toggleStageFile(file.path)"
                    class="file-checkbox"
                  >
                  <span class="file-status modified">{{ file.status }}</span>
                  <span class="file-path">{{ file.path }}</span>
                  <div class="file-actions">
                    <button @click="discardChanges(file.path)" class="small-btn danger">丢弃</button>
                  </div>
                </div>
              </div>
            </div>

            <div class="file-changes-section">
              <h4 class="section-subtitle">
                <span class="file-change-icon staged">✅</span>
                已暂存 ({{ stagedFiles.length }})
              </h4>
              <div class="file-list">
                <div 
                  v-for="file in stagedFiles" 
                  :key="'staged-' + file.path"
                  class="file-item"
                >
                  <input 
                    type="checkbox" 
                    checked
                    @change="toggleUnstageFile(file.path)"
                    class="file-checkbox"
                  >
                  <span class="file-status staged">{{ file.status }}</span>
                  <span class="file-path">{{ file.path }}</span>
                  <div class="file-actions">
                    <button @click="unstageFile(file.path)" class="small-btn secondary">取消</button>
                  </div>
                </div>
              </div>
            </div>

            <!-- 提交区域 -->
            <div class="commit-section">
              <div class="commit-input-group">
                <input 
                  v-model="commitMessage" 
                  type="text" 
                  placeholder="提交摘要 (必填)" 
                  class="commit-summary-input"
                >
              </div>
              <textarea 
                v-model="commitDescription" 
                placeholder="详细描述 (可选)" 
                class="commit-description-input"
              ></textarea>
              <button 
                @click="commitChanges" 
                :disabled="!canCommit" 
                class="commit-action-btn"
              >
                <span class="btn-icon">📤</span>
                提交 ({{ stagedFiles.length }})
              </button>
            </div>
          </div>

          <!-- 当没有更改时的空状态 -->
          <div v-else class="empty-staging-view">
            <div class="empty-icon">✨</div>
            <p>工作区干净，无需提交</p>
          </div>
        </div>
      </aside>
    </div>

    <!-- 底部状态栏 -->
    <footer class="status-bar">
      <div class="status-left">
        <span class="current-branch">分支: {{ currentBranch || '无' }}</span>
        <span class="repo-path">{{ repoPath || '未选择仓库' }}</span>
      </div>
      <div class="status-center">
        <span class="change-summary">
          <span class="working-changes">修改: {{ workingFiles.length }}</span>
          <span class="staged-changes">暂存: {{ stagedFiles.length }}</span>
        </span>
      </div>
      <div class="status-right">
        <span class="connection-status">
          <span :class="['status-indicator', { 'connected': repoPath, 'disconnected': !repoPath }]" ></span>
          <span>{{ repoPath ? '已连接' : '未连接' }}</span>
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
import {computed, onMounted, reactive, ref, onUpdated} from 'vue'
import BranchList from './components/BranchList.vue'

// 导入Wails运行时和Go模块
// 注意：在生产环境中，Wails会在运行时注入这些对象，所以不需要显式导入
// 我们将在运行时通过window.go访问这些函数

export default {
  name: 'App',
  components: {
    BranchList
  },
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
    const commitMessage = ref('') // 提交信息摘要
    const commitDescription = ref('') // 提交信息详细描述
    const branchFilter = ref('')  // 分支过滤器
    const selectedCommit = ref(null) // 当前选中的提交
    const pullCount = ref(0) // 落后提交数
    
    // 控制折叠面板展开状态
    const expandedSections = reactive({
      localBranches: true,
      remoteBranches: true,
      tags: false
    })

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
        // 返回相对时间
        const now = new Date()
        const diffMs = now - date
        const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))
        
        if (diffDays === 0) {
          const diffHours = Math.floor(diffMs / (1000 * 60 * 60))
          if (diffHours === 0) {
            const diffMinutes = Math.floor(diffMs / (1000 * 60))
            return `${diffMinutes}分钟前`
          }
          return `${diffHours}小时前`
        } else if (diffDays === 1) {
          return '昨天'
        } else if (diffDays < 7) {
          return `${diffDays}天前`
        } else {
          return date.toLocaleDateString('zh-CN', {
            year: 'numeric',
            month: 'short',
            day: 'numeric'
          })
        }
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
    
    // 切换文件暂存状态
    const toggleStageFile = async (filePath) => {
      await stageFile(filePath)
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
    
    // 切换文件取消暂存状态
    const toggleUnstageFile = async (filePath) => {
      await unstageFile(filePath)
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
        const fullMessage = commitDescription.value ? 
          `${commitMessage.value}\n\n${commitDescription.value}` : 
          commitMessage.value
        
        await window.go.main.App.GitCommit(repoPath.value, fullMessage)
        commitMessage.value = ''
        commitDescription.value = ''
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
    
    // 选择分支（不切换，只高亮显示）
    const selectBranch = (branchName) => {
      console.log(`选中分支: ${branchName}`)
    }
    
    // 选择远程分支
    const selectRemoteBranch = (branchName) => {
      console.log(`选中远程分支: ${branchName}`)
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
    
    // 获取最新更改
    const fetchChanges = async () => {
      if (!repoPath.value) {
        showNotification('请先加载仓库', 'error')
        return
      }

      try {
        const result = await window.go.main.App.GitFetch(repoPath.value)
        showNotification(`获取成功: ${result || '无新更改'}`, 'success')
        // 更新落后的提交数
        updatePullCount()
      } catch (error) {
        showNotification(`获取失败: ${error}`, 'error')
      }
    }
    
    // 更新落后的提交数
    const updatePullCount = async () => {
      if (!repoPath.value) return
      
      try {
        // 这里应该调用一个计算落后提交数的API
        // 暂时设置为模拟值
        pullCount.value = Math.floor(Math.random() * 5) // 模拟随机落后数
      } catch (error) {
        console.error('更新落后提交数失败:', error)
      }
    }
    
    // 切换折叠面板
    const toggleSection = (section) => {
      expandedSections[section] = !expandedSections[section]
    }

    // 页面加载时初始化
    onMounted(async () => {
      // 尝试加载默认仓库
      if (repoPath.value) {
        await loadRepo()
      }
    })
    
    // 组件更新后重新计算落后提交数
    onUpdated(() => {
      if (repoPath.value) {
        updatePullCount()
      }
    })

    return {
      repoPath,
      currentBranch,
      allBranches,
      commits,
      branchesLoading,
      commitsLoading,
      statusLoading,
      localBranches,
      remoteBranches,
      workingFiles,
      stagedFiles,
      commitMessage,
      commitDescription,
      canCommit,
      selectedCommit,
      pullCount,
      expandedSections,
      notification,
      branchFilter,
      loadRepo,
      refreshData,
      loadBranches,
      loadCommits,
      loadCurrentBranch,
      loadStatus,
      refreshStatus,
      switchBranch,
      selectBranch,
      selectRemoteBranch,
      createBranch,
      createBranchFromRemote,
      showStatus,
      refreshBranches,
      refreshCommits,
      stageFile,
      unstageFile,
      toggleStageFile,
      toggleUnstageFile,
      stageAll,
      discardChanges,
      commitChanges,
      browseRepo,
      pullChanges,
      pushChanges,
      fetchChanges,
      selectCommit,
      showBranchHistory,
      openBranchContextMenu,
      getCommitColor,
      getCommitLineColor,
      formatDate,
      toggleSection,
      showNotification
    }
  }
}
</script>

<style src="./styles/app.css"></style>