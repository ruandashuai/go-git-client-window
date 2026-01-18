<template>
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
</template>

<script>
export default {
  name: 'StatusBar',
  props: {
    currentBranch: {
      type: String,
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
  }
}
</script>