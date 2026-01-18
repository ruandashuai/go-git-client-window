<template>
  <section class="panel file-status-panel">
    <div class="panel-header">
      <h2 class="panel-title">
        <span class="status-icon">📋</span>
        工作区状态
      </h2>
      <div class="panel-actions">
        <button @click="$emit('show-status')" class="icon-btn" title="查看详细状态">
          👁️
        </button>
        <button @click="$emit('stage-all')" class="icon-btn" title="暂存全部">
          📥
        </button>
        <button @click="$emit('refresh-status')" class="icon-btn" title="刷新状态">
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
                <button @click="$emit('stage-file', file.path)" class="small-btn primary">
                  暂存
                </button>
                <button @click="$emit('discard-changes', file.path)" class="small-btn danger">
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
                <button @click="$emit('unstage-file', file.path)" class="small-btn secondary">
                  取消暂存
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
export default {
  name: 'FileStatusPanel',
  props: {
    statusLoading: {
      type: Boolean,
      required: true
    },
    repoPath: {
      type: String,
      required: true
    },
    workingFiles: {
      type: Array,
      required: true
    },
    stagedFiles: {
      type: Array,
      required: true
    }
  },
  emits: [
    'show-status',
    'stage-all',
    'refresh-status',
    'stage-file',
    'discard-changes',
    'unstage-file'
  ]
}
</script>