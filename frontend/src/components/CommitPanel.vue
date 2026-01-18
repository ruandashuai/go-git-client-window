<template>
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
            :value="commitMessage"
            @input="$emit('update:commit-message', $event.target.value)"
            placeholder="输入提交信息..."
            class="commit-message-input"
            rows="3"
        ></textarea>
        <div class="commit-actions">
          <button
              @click="$emit('commit-changes')"
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
</template>

<script>
export default {
  name: 'CommitPanel',
  props: {
    commitMessage: {
      type: String,
      required: true
    },
    stagedFiles: {
      type: Array,
      required: true
    },
    canCommit: {
      type: Boolean,
      required: true
    }
  },
  emits: ['commit-changes']
}
</script>