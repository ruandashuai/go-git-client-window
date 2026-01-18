<template>
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
            @click="$emit('select-commit', commit)"
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
</template>

<script>
export default {
  name: 'CommitHistory',
  props: {
    commits: {
      type: Array,
      required: true
    },
    commitsLoading: {
      type: Boolean,
      required: true
    },
    repoPath: {
      type: String,
      required: true
    },
    selectedCommit: {
      type: Object,
      default: null
    }
  },
  emits: ['select-commit'],
  methods: {
    // 截断文本以适应显示
    truncateText(text, maxLength) {
      if (!text) return ''
      return text.length > maxLength ? text.substring(0, maxLength) + '...' : text
    },
    
    // 格式化日期显示
    formatDate(dateString) {
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
    },
    
    // 获取提交颜色
    getCommitColor(index) {
      // 基于索引生成不同的颜色
      const hue = (index * 137.5) % 360 // 使用黄金角度生成颜色差异
      return {
        backgroundColor: `hsl(${hue}, 70%, 60%)`
      }
    },
    
    // 获取提交线条颜色
    getCommitLineColor(index) {
      const hue = (index * 137.5) % 360
      return {
        borderColor: `hsl(${hue}, 70%, 60%)`
      }
    }
  }
}
</script>