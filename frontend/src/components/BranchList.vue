<template>
  <aside class="sidebar left-sidebar">
    <div class="panel-header">
      <h2 class="panel-title">
        <span class="branch-icon">🌱</span>
        分支管理
      </h2>
      <div class="panel-actions">
        <button @click="$emit('create-branch')" class="icon-btn" title="创建分支">
          ➕
        </button>
        <button @click="$emit('refresh-branches')" class="icon-btn" title="刷新">
          🔄
        </button>
      </div>
    </div>

    <div class="panel-content">
      <!-- 分支搜索 -->
      <div class="search-box">
        <input
            type="text"
            :value="branchFilter"
            @input="$emit('update:branch-filter', $event.target.value)"
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
               @dblclick="$emit('show-branch-history', branch.name)"
               @contextmenu.prevent="$emit('open-branch-context-menu', $event, branch, 'local')"
          >
            <div class="branch-info">
              <span class="branch-name" :title="branch.name">{{ branch.name }}</span>
              <span v-if="branch.current" class="branch-current-badge" >（当前分支）</span>
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
          <span class="item-count">({{ (remoteBranches?.value?.length || 0) }})</span>
        </div>
        <div class="branch-list">
          <div
              v-for="branch in remoteBranches"
              :key="'remote-' + branch.name"
              :class="['branch-item', { 'active': branch.current }]"
              @dblclick="$emit('show-branch-history', branch.name.replace('origin/', ''))"
              @contextmenu.prevent="$emit('open-branch-context-menu', $event, branch, 'remote')"
          >
            <div class="branch-info">
              <span class="branch-type">📡</span>
              <span class="branch-name" :title="branch.name">{{ branch.name }}</span>
              <span v-if="branch.current" class="branch-current-badge" title="当前分支">●</span>
            </div>
            <div class="branch-actions">
              <button
                  v-if="!branch.current"
                  @click.stop="$emit('switch-branch', branch.name)"
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
</template>

<script>
import {computed} from 'vue'

export default {
  name: 'BranchList',
  props: {
    allBranches: {
      type: Array,
      required: true
    },
    branchFilter: {
      type: String,
      required: true
    }
  },
  emits: [
    'switch-branch',
    'create-branch',
    'delete-branch',
    'refresh-branches',
    'show-branch-history',
    'open-branch-context-menu'
  ],
  setup(props) {
    // 计算属性：将分支分为本地分支和远程分支
    const localBranches = computed(() => {
      if (!props.allBranches || !Array.isArray(props.allBranches)) {
        return []
      }
      let resultList = props.allBranches.filter(branch => !branch.remote);
      if (props.branchFilter && props.branchFilter.trim() !== '') {
        let searchLower = props.branchFilter.toLowerCase().trim();
        resultList = resultList.filter(branch =>
            branch.name.toLowerCase().includes(searchLower)
        )
      }
      return resultList
    })

    const remoteBranches = computed(() => {
      if (!props.allBranches || !Array.isArray(props.allBranches)) {
        return []
      }
      let resultList = props.allBranches.filter(branch => branch.remote);
      if (props.branchFilter && props.branchFilter.trim() !== '') {
        let searchLower = props.branchFilter.toLowerCase().trim();
        resultList = resultList.filter(branch =>
            branch.name.toLowerCase().includes(searchLower)
        )
      }
      return resultList
    })

    return {
      localBranches,
      remoteBranches
    }
  }
}
</script>