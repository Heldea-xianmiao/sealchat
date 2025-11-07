<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import dayjs from 'dayjs'
import { useDraggable } from '@vueuse/core'
import { SearchOutline, CloseOutline, ChevronDownOutline, ChevronUpOutline } from '@vicons/ionicons5'
import { useChannelSearchStore } from '@/stores/channelSearch'
import { useChatStore } from '@/stores/chat'

interface JumpPayload {
  messageId: string
  displayOrder?: number
}

const emit = defineEmits<{
  (event: 'jump-to-message', payload: JumpPayload): void
}>()

const PANEL_WIDTH = 420

const chat = useChatStore()
const channelSearch = useChannelSearchStore()

const {
  panelVisible,
  keyword,
  matchMode,
  filters,
  loading,
  results,
  total,
  page,
  pageSize,
  error,
  lastKeyword,
} = storeToRefs(channelSearch)

const searchInputRef = ref<HTMLInputElement | null>(null)
const panelRef = ref<HTMLElement | null>(null)
const dragHandleRef = ref<HTMLElement | null>(null)
const advancedFiltersVisible = ref(false)

channelSearch.bindChannel(chat.curChannel?.id || null)

const { x, y } = useDraggable(panelRef, {
  handle: dragHandleRef,
  initialValue: channelSearch.panelPosition,
})

const panelStyle = computed(() => ({
  left: `${x.value}px`,
  top: `${y.value}px`,
  width: `${PANEL_WIDTH}px`,
  zIndex: 6000,
}))

const clampPosition = (nx: number, ny: number) => {
  if (typeof window === 'undefined') {
    return { x: nx, y: ny }
  }
  const maxX = Math.max(16, window.innerWidth - PANEL_WIDTH - 16)
  const maxY = Math.max(100, window.innerHeight - 160)
  return {
    x: Math.min(Math.max(16, nx), maxX),
    y: Math.min(Math.max(80, ny), maxY),
  }
}

watch(
  () => panelVisible.value,
  (visible) => {
    if (visible) {
      channelSearch.bindChannel(chat.curChannel?.id || null)
      nextTick(() => searchInputRef.value?.focus())
      const { x: px, y: py } = channelSearch.panelPosition
      x.value = px
      y.value = py
    }
  },
  { immediate: true },
)

watch(
  () => chat.curChannel?.id,
  (channelId) => {
    channelSearch.bindChannel(channelId || null)
    if (panelVisible.value) {
      channelSearch.search(channelId || undefined)
    }
  },
)

watch(
  [x, y],
  ([nx, ny]) => {
    if (!panelVisible.value) return
    const clamped = clampPosition(nx, ny)
    if (clamped.x !== nx) {
      x.value = clamped.x
    }
    if (clamped.y !== ny) {
      y.value = clamped.y
    }
    channelSearch.setPanelPosition({ x: clamped.x, y: clamped.y })
  },
  { flush: 'post' },
)

const memberOptions = computed(() =>
  chat.curChannelUsers.map((member) => ({
    label: member.nick || member.name || '未知成员',
    value: member.id || '',
  })),
)

const hasSearched = computed(() => !!lastKeyword.value)
const showEmptyState = computed(() => hasSearched.value && !loading.value && results.value.length === 0)
const filterActive = computed(() => channelSearch.isFilterActive)
const activeFilterCount = computed(() => {
  let count = 0
  const current = filters.value
  if (current.speakerIds.length) count++
  if (current.archived !== 'all') count++
  if (current.icMode !== 'all') count++
  if (current.includeOutside === false) count++
  if (current.timeRange) count++
  return count
})

const speakerFilter = computed({
  get: () => filters.value.speakerIds,
  set: (val: string[]) => channelSearch.updateFilters({ speakerIds: val }),
})

const archivedFilter = computed({
  get: () => filters.value.archived,
  set: (val: 'all' | 'only' | 'exclude') => channelSearch.updateFilters({ archived: val }),
})

const icModeFilter = computed({
  get: () => filters.value.icMode,
  set: (val: 'all' | 'ic' | 'ooc') => channelSearch.updateFilters({ icMode: val }),
})

const includeOutsideFilter = computed({
  get: () => filters.value.includeOutside,
  set: (val: boolean) => channelSearch.updateFilters({ includeOutside: val }),
})

const timeRangeFilter = computed({
  get: () => filters.value.timeRange,
  set: (val: [number, number] | null) => channelSearch.updateFilters({ timeRange: val }),
})

const matchModeValue = computed({
  get: () => matchMode.value,
  set: (val: 'fuzzy' | 'exact') => channelSearch.setMatchMode(val),
})

const handleClose = () => {
  channelSearch.closePanel()
}

const toggleAdvancedFilters = () => {
  advancedFiltersVisible.value = !advancedFiltersVisible.value
}

const runSearch = () => {
  channelSearch.setPage(1)
  channelSearch.search(chat.curChannel?.id || undefined)
}

const handleEnter = (event: KeyboardEvent) => {
  if (event.key === 'Enter') {
    runSearch()
  }
}

const handlePageChange = (nextPage: number) => {
  channelSearch.setPage(nextPage)
  channelSearch.search(chat.curChannel?.id || undefined)
}

const handleResultClick = (item: JumpPayload) => {
  emit('jump-to-message', item)
}

const formatTime = (timestamp: number) => {
  if (!timestamp) {
    return '未知时间'
  }
  return dayjs(timestamp).format('YYYY-MM-DD HH:mm')
}

const escapeRegExp = (value: string) => value.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')

const renderSnippetFragments = (content: string, ranges?: Array<[number, number]>) => {
  if (!content) {
    return [{ text: '（无内容）', highlighted: false }]
  }
  if (ranges && ranges.length > 0) {
    const fragments: Array<{ text: string; highlighted: boolean }> = []
    let cursor = 0
    ranges
      .slice()
      .sort((a, b) => a[0] - b[0])
      .forEach(([start, end]) => {
        const safeStart = Math.max(0, start)
        const safeEnd = Math.min(content.length, end)
        if (safeStart > cursor) {
          fragments.push({ text: content.slice(cursor, safeStart), highlighted: false })
        }
        fragments.push({ text: content.slice(safeStart, safeEnd), highlighted: true })
        cursor = safeEnd
      })
    if (cursor < content.length) {
      fragments.push({ text: content.slice(cursor), highlighted: false })
    }
    return fragments
  }
  const keywordValue = lastKeyword.value.trim()
  if (!keywordValue) {
    return [{ text: content, highlighted: false }]
  }
  const fragments: Array<{ text: string; highlighted: boolean }> = []
  const pattern = new RegExp(escapeRegExp(keywordValue), 'gi')
  let lastIndex = 0
  const matches = content.matchAll(pattern)
  for (const match of matches) {
    if (!match[0]) continue
    const start = match.index ?? 0
    const end = start + match[0].length
    if (start > lastIndex) {
      fragments.push({ text: content.slice(lastIndex, start), highlighted: false })
    }
    fragments.push({ text: content.slice(start, end), highlighted: true })
    lastIndex = end
  }
  if (lastIndex < content.length) {
    fragments.push({ text: content.slice(lastIndex), highlighted: false })
  }
  return fragments.length ? fragments : [{ text: content, highlighted: false }]
}

const shortContent = (text: string) => {
  if (!text) return ''
  return text.length > 200 ? `${text.slice(0, 200)}...` : text
}
</script>

<template>
  <Teleport to="body">
    <transition name="fade">
      <div v-if="panelVisible" ref="panelRef" class="chat-search-panel" :style="panelStyle">
        <div ref="dragHandleRef" class="chat-search-panel__header">
          <div>
            <div class="chat-search-panel__title">频道搜索</div>
            <div class="chat-search-panel__subtitle">
              {{ chat.curChannel?.name || '未选择频道' }}
            </div>
          </div>
          <div class="chat-search-panel__header-actions">
            <n-tag size="small" type="info" v-if="hasSearched">
              共 {{ total }} 条
            </n-tag>
            <button class="chat-search-panel__close" type="button" @click="handleClose" aria-label="关闭搜索面板">
              <n-icon size="16">
                <CloseOutline />
              </n-icon>
            </button>
          </div>
        </div>

        <div class="chat-search-panel__body">
          <div class="chat-search-panel__input-group">
            <n-input
              ref="searchInputRef"
              v-model:value="keyword"
              placeholder="输入关键字，按回车搜索"
              size="large"
              clearable
              @keyup="handleEnter"
            >
              <template #prefix>
                <n-icon size="16">
                  <SearchOutline />
                </n-icon>
              </template>
              <template #suffix>
                <n-button type="primary" ghost size="small" @click="runSearch">
                  搜索
                </n-button>
              </template>
            </n-input>
          </div>

          <div class="chat-search-panel__filter-toggle">
            <n-button
              type="primary"
              ghost
              strong
              size="small"
              class="filter-toggle-button"
              @click="toggleAdvancedFilters"
              :aria-expanded="advancedFiltersVisible"
            >
              <n-icon size="16" class="mr-1">
                <component :is="advancedFiltersVisible ? ChevronUpOutline : ChevronDownOutline" />
              </n-icon>
              {{ advancedFiltersVisible ? '收起筛选' : '展开筛选' }}
            </n-button>
            <n-tag v-if="filterActive" size="small" type="warning" round>
              {{ activeFilterCount }} 项筛选
            </n-tag>
          </div>

          <transition name="expand">
            <div v-if="advancedFiltersVisible" class="chat-search-panel__filter-bar">
              <div class="filter-group">
                <span class="filter-label">模式</span>
                <n-radio-group v-model:value="matchModeValue" size="small">
                  <n-radio-button value="fuzzy">模糊</n-radio-button>
                  <n-radio-button value="exact">精准</n-radio-button>
                </n-radio-group>
              </div>

              <div class="filter-group">
                <span class="filter-label">场内/场外</span>
                <n-radio-group v-model:value="icModeFilter" size="small">
                  <n-radio-button value="all">全部</n-radio-button>
                  <n-radio-button value="ic">场内</n-radio-button>
                  <n-radio-button value="ooc">场外</n-radio-button>
                </n-radio-group>
              </div>

              <div class="filter-group">
                <span class="filter-label">归档</span>
                <n-radio-group v-model:value="archivedFilter" size="small">
                  <n-radio-button value="all">全部</n-radio-button>
                  <n-radio-button value="exclude">未归档</n-radio-button>
                  <n-radio-button value="only">仅归档</n-radio-button>
                </n-radio-group>
              </div>

              <div class="filter-group filter-group--inline">
                <span class="filter-label">发言人</span>
                <n-select
                  placeholder="筛选成员"
                  :options="memberOptions"
                  v-model:value="speakerFilter"
                  multiple
                  clearable
                  size="small"
                  style="width: 180px"
                />
              </div>

              <div class="filter-group filter-group--inline">
                <span class="filter-label">时间范围</span>
                <n-date-picker
                  v-model:value="timeRangeFilter"
                  type="datetimerange"
                  :update-value-on-close="true"
                  clearable
                  size="small"
                  :value-format="'timestamp'"
                />
              </div>

              <div class="filter-group filter-group--inline">
                <span class="filter-label">场外消息</span>
                <n-switch v-model:value="includeOutsideFilter" size="small">
                  <template #checked>包含</template>
                  <template #unchecked>忽略</template>
                </n-switch>
              </div>

              <n-button v-if="filterActive" size="tiny" tertiary @click="channelSearch.resetFilters">
                重置筛选
              </n-button>
            </div>
          </transition>

          <n-alert v-if="error" type="error" class="mt-2" :bordered="false">
            {{ error }}
          </n-alert>

          <div class="chat-search-panel__results">
            <n-spin :show="loading">
              <n-empty
                v-if="showEmptyState"
                description="没有匹配的结果，尝试更换关键词或放宽筛选条件"
                size="small"
              />
              <div v-else>
                <div v-for="item in results" :key="item.id" class="search-result">
                  <div class="search-result__meta">
                    <div class="search-result__author">{{ item.senderName }}</div>
                    <div class="search-result__time">{{ formatTime(item.createdAt) }}</div>
                  </div>
                  <div class="search-result__badges">
                    <n-tag size="small" :type="item.icMode === 'ic' ? 'success' : 'default'" round>
                      {{ item.icMode === 'ic' ? '场内' : '场外' }}
                    </n-tag>
                    <n-tag v-if="item.isArchived" size="small" round type="warning">
                      已归档
                    </n-tag>
                  </div>
                  <div class="search-result__content">
                    <template v-for="(fragment, idx) in renderSnippetFragments(shortContent(item.contentSnippet), item.highlightRanges)" :key="idx">
                      <mark v-if="fragment.highlighted">{{ fragment.text }}</mark>
                      <span v-else>{{ fragment.text }}</span>
                    </template>
                  </div>
                  <div class="search-result__actions">
                    <span class="search-result__hint">点击下方按钮跳转到该消息</span>
                    <n-button size="tiny" type="primary" ghost @click="handleResultClick({ messageId: item.id, displayOrder: item.displayOrder })">
                      跳转
                    </n-button>
                  </div>
                </div>
              </div>
            </n-spin>
          </div>

          <div class="chat-search-panel__footer" v-if="results.length || total > pageSize">
            <n-pagination
              :page="page"
              :page-size="pageSize"
              :item-count="total"
              :show-size-picker="false"
              @update:page="handlePageChange"
            />
          </div>
        </div>
      </div>
    </transition>
  </Teleport>
</template>

<style scoped lang="scss">
.chat-search-panel {
  position: fixed;
  top: 120px;
  right: 40px;
  background: rgba(255, 255, 255, 0.98);
  border-radius: 1rem;
  box-shadow: 0 20px 60px rgba(15, 23, 42, 0.25);
  border: 1px solid rgba(148, 163, 184, 0.3);
  padding: 1rem 1.25rem 1.5rem;
  backdrop-filter: blur(12px);
  display: flex;
  flex-direction: column;
  max-height: calc(100vh - 160px);
}

.chat-search-panel__header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  cursor: grab;
  gap: 0.75rem;
}

.chat-search-panel__title {
  font-size: 1.1rem;
  font-weight: 600;
  color: #0f172a;
}

.chat-search-panel__subtitle {
  font-size: 0.85rem;
  color: #64748b;
  margin-top: 0.15rem;
}

.chat-search-panel__header-actions {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.chat-search-panel__close {
  width: 2rem;
  height: 2rem;
  border-radius: 999px;
  border: none;
  background: rgba(15, 23, 42, 0.05);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: #475569;
  cursor: pointer;
  transition: background 0.2s ease, color 0.2s ease;
}

.chat-search-panel__close:hover {
  background: rgba(15, 23, 42, 0.12);
  color: #0f172a;
}

.chat-search-panel__body {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  margin-top: 1rem;
  overflow: visible;
}

.chat-search-panel__input-group {
  width: 100%;
}

.chat-search-panel__filter-bar {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  border: 1px dashed rgba(148, 163, 184, 0.6);
  border-radius: 0.75rem;
  padding: 0.75rem;
  background: rgba(248, 250, 252, 0.8);
}

.filter-group {
  display: flex;
  align-items: center;
  gap: 0.65rem;
  flex-wrap: wrap;
}

.filter-group--inline {
  flex-wrap: nowrap;
}

.filter-label {
  font-size: 0.8rem;
  color: #475569;
  min-width: 60px;
}

.chat-search-panel__results {
  flex: 1;
  min-height: 220px;
  overflow-y: auto;
  padding-right: 0.25rem;
}

.chat-search-panel__filter-toggle {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.75rem;
  margin-top: 0.25rem;
}

.filter-toggle-button {
  letter-spacing: 0.02em;
}

.search-result {
  border-radius: 0.75rem;
  border: 1px solid rgba(226, 232, 240, 0.9);
  padding: 0.85rem;
  margin-bottom: 0.75rem;
  background: rgba(255, 255, 255, 0.9);
  box-shadow: 0 10px 20px rgba(15, 23, 42, 0.05);
}

.search-result__meta {
  display: flex;
  justify-content: space-between;
  gap: 0.5rem;
  font-size: 0.85rem;
  color: #0f172a;
  font-weight: 500;
}

.search-result__time {
  color: #94a3b8;
  font-weight: 400;
}

.search-result__badges {
  display: flex;
  gap: 0.5rem;
  margin-top: 0.35rem;
}

.search-result__content {
  margin-top: 0.5rem;
  line-height: 1.45;
  color: #1f2937;
  font-size: 0.9rem;
}

.search-result__content mark {
  background: rgba(14, 165, 233, 0.2);
  color: #0f172a;
  border-radius: 0.15rem;
  padding: 0 0.1rem;
}

.search-result__actions {
  margin-top: 0.65rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.search-result__hint {
  font-size: 0.75rem;
  color: #94a3b8;
}

.chat-search-panel__footer {
  display: flex;
  justify-content: center;
  margin-top: 0.5rem;
  padding-top: 0.5rem;
  border-top: 1px solid rgba(226, 232, 240, 0.6);
}

.expand-enter-active,
.expand-leave-active {
  transition: all 0.25s ease;
}

.expand-enter-from,
.expand-leave-to {
  opacity: 0;
  transform: translateY(-6px);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.25s ease, transform 0.25s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}
</style>
