<template>
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
            :value="repoPath"
            @input="$emit('update:repo-path', $event.target.value)"
            placeholder="输入 Git 仓库路径"
            @keyup.enter="$emit('load-repo')"
            class="repo-path-input"
        />
        <button @click="$emit('browse-repo')" class="browse-btn" title="浏览目录">
          📁
        </button>
      </div>
      <div class="repo-actions">
        <button @click="$emit('load-repo')" class="primary-btn">
          <span class="btn-icon">📂</span>
          加载仓库
        </button>
        <button @click="$emit('refresh-data')" class="secondary-btn">
          <span class="btn-icon">🔄</span>
          刷新
        </button>
        <button @click="$emit('pull-changes')" class="secondary-btn">
          <span class="btn-icon">⬇️</span>
          拉取
        </button>
        <button @click="$emit('push-changes')" class="secondary-btn">
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
</template>

<script>
export default {
  name: 'TopNavBar',
  props: {
    repoPath: {
      type: String,
      required: true
    }
  },
  emits: [
    'browse-repo',
    'load-repo',
    'refresh-data',
    'pull-changes',
    'push-changes'
  ]
}
</script>