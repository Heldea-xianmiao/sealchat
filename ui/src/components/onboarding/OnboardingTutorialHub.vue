<script setup lang="ts">
/**
 * OnboardingTutorialHub - 教程中心组件
 * 展示所有教程模块，支持分类选择和批量学习
 */
import { computed, onMounted } from 'vue'
import { NButton, NIcon, NCheckbox, NTag, NProgress } from 'naive-ui'
import { X as CloseIcon } from '@vicons/tabler'
import { useOnboardingStore } from '@/stores/onboarding'
import { TUTORIAL_CATEGORIES, RECOMMENDED_MODULES, formatDuration } from '@/config/tutorialModules'

const onboarding = useOnboardingStore()

// 打开时自动选中推荐模块（仅当没有选中任何模块时）
onMounted(() => {
  if (onboarding.selectedModules.length === 0) {
    onboarding.selectRecommendedModules(RECOMMENDED_MODULES)
  }
})

const groupedModules = computed(() => onboarding.modulesByCategory)

const isModuleSelected = (id: string) => onboarding.selectedModules.includes(id)
const isModuleCompleted = (id: string) => onboarding.isModuleCompleted(id)

const handleToggleModule = (id: string) => {
  onboarding.toggleModuleSelection(id)
}

const handleSelectRecommended = () => {
  onboarding.selectRecommendedModules(RECOMMENDED_MODULES)
}

const handleSelectAll = () => {
  onboarding.selectAllModules()
}

const handleClearSelection = () => {
  onboarding.clearSelection()
}

const handleStartSelected = () => {
  onboarding.startSelectedModules()
}

const handleClose = () => {
  onboarding.closeTutorialHub()
}

const selectedCount = computed(() => onboarding.selectedModules.length)
const hasSelection = computed(() => selectedCount.value > 0)
</script>

<template>
  <div class="tutorial-hub-overlay" @click.self="handleClose">
    <div class="tutorial-hub">
      <!-- 头部 -->
      <header class="tutorial-hub__header">
        <div class="header-content">
          <h2>📚 功能教程中心</h2>
          <p>选择你想了解的功能，开始学习</p>
        </div>
        <n-button quaternary circle size="small" class="close-btn" @click="handleClose">
          <template #icon>
            <n-icon :component="CloseIcon" />
          </template>
        </n-button>
      </header>

      <!-- 进度条 -->
      <div class="tutorial-hub__progress">
        <div class="progress-label">
          <span>学习进度</span>
          <span class="progress-value">{{ onboarding.completionPercentage }}%</span>
        </div>
        <n-progress
          :percentage="onboarding.completionPercentage"
          :show-indicator="false"
          :height="6"
          :border-radius="3"
        />
      </div>

      <!-- 模块列表 -->
      <div class="tutorial-hub__content">
        <div
          v-for="cat in TUTORIAL_CATEGORIES"
          :key="cat.id"
          class="category-section"
        >
          <div class="category-header">
            <h3 class="category-title">{{ cat.label }}</h3>
            <span class="category-desc">{{ cat.description }}</span>
          </div>

          <div class="module-grid">
            <div
              v-for="mod in groupedModules[cat.id]"
              :key="mod.id"
              class="module-card"
              :class="{
                'module-card--selected': isModuleSelected(mod.id),
                'module-card--completed': isModuleCompleted(mod.id),
              }"
              @click="handleToggleModule(mod.id)"
            >
              <div class="module-card__checkbox">
                <n-checkbox :checked="isModuleSelected(mod.id)" @click.stop />
              </div>
              <div class="module-card__content">
                <div class="module-card__title">
                  <span>{{ mod.title }}</span>
                  <n-tag v-if="isModuleCompleted(mod.id)" size="tiny" type="success">
                    已学习
                  </n-tag>
                </div>
                <div class="module-card__desc">{{ mod.description }}</div>
                <div class="module-card__meta">
                  <span>{{ mod.steps.length }} 步骤</span>
                  <span>约 {{ formatDuration(mod.estimatedTime) }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 底部操作 -->
      <footer class="tutorial-hub__footer">
        <div class="footer-quick-select">
          <n-button size="small" quaternary @click="handleSelectRecommended">
            推荐入门
          </n-button>
          <n-button size="small" quaternary @click="handleSelectAll">
            全选
          </n-button>
          <n-button v-if="hasSelection" size="small" quaternary @click="handleClearSelection">
            清空
          </n-button>
        </div>
        <div class="footer-actions">
          <n-button @click="handleClose">稍后再说</n-button>
          <n-button
            type="primary"
            :disabled="!hasSelection"
            @click="handleStartSelected"
          >
            开始学习 ({{ selectedCount }})
          </n-button>
        </div>
      </footer>
    </div>
  </div>
</template>

<style scoped>
.tutorial-hub-overlay {
  position: fixed;
  inset: 0;
  z-index: 10000;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
}

.tutorial-hub {
  background: var(--sc-bg-elevated, #fff);
  border-radius: 16px;
  max-width: 720px;
  width: 95%;
  max-height: 85vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  animation: slideUp 0.3s ease-out;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.tutorial-hub__header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 1.5rem 1.5rem 1rem;
  border-bottom: 1px solid var(--sc-border-mute, #e5e7eb);
}

.header-content h2 {
  margin: 0 0 0.25rem;
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--sc-text-primary, #1f2937);
}

.header-content p {
  margin: 0;
  font-size: 0.875rem;
  color: var(--sc-text-secondary, #6b7280);
}

.close-btn {
  margin-top: -0.25rem;
  margin-right: -0.5rem;
}

.tutorial-hub__progress {
  padding: 1rem 1.5rem;
  background: var(--sc-bg-surface, #f9fafb);
}

.progress-label {
  display: flex;
  justify-content: space-between;
  font-size: 0.8125rem;
  color: var(--sc-text-secondary, #6b7280);
  margin-bottom: 0.5rem;
}

.progress-value {
  font-weight: 600;
  color: var(--primary-color, #3388de);
}

.tutorial-hub__content {
  flex: 1;
  overflow-y: auto;
  padding: 1rem 1.5rem;
}

.category-section {
  margin-bottom: 1.5rem;
}

.category-section:last-child {
  margin-bottom: 0.5rem;
}

.category-header {
  margin-bottom: 0.75rem;
}

.category-title {
  margin: 0;
  font-size: 1rem;
  font-weight: 600;
  color: var(--sc-text-primary, #374151);
}

.category-desc {
  font-size: 0.75rem;
  color: var(--sc-text-secondary, #9ca3af);
}

.module-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 0.75rem;
}

.module-card {
  display: flex;
  gap: 0.75rem;
  padding: 0.875rem;
  background: var(--sc-bg-surface, #f9fafb);
  border: 1px solid var(--sc-border-mute, #e5e7eb);
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.module-card:hover {
  border-color: var(--primary-color, #3388de);
  background: var(--sc-bg-elevated, #fff);
}

.module-card--selected {
  border-color: var(--primary-color, #3388de);
  background: rgba(51, 136, 222, 0.06);
}

.module-card--completed {
  opacity: 0.7;
}

.module-card__checkbox {
  flex-shrink: 0;
  padding-top: 2px;
}

.module-card__content {
  flex: 1;
  min-width: 0;
}

.module-card__title {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-weight: 600;
  font-size: 0.9375rem;
  color: var(--sc-text-primary, #1f2937);
  margin-bottom: 0.25rem;
}

.module-card__desc {
  font-size: 0.8125rem;
  color: var(--sc-text-secondary, #6b7280);
  margin-bottom: 0.5rem;
  line-height: 1.4;
}

.module-card__meta {
  display: flex;
  gap: 1rem;
  font-size: 0.75rem;
  color: var(--sc-text-secondary, #9ca3af);
}

.tutorial-hub__footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 1.5rem;
  border-top: 1px solid var(--sc-border-mute, #e5e7eb);
  background: var(--sc-bg-surface, #f9fafb);
  border-radius: 0 0 16px 16px;
}

.footer-quick-select {
  display: flex;
  gap: 0.25rem;
}

.footer-actions {
  display: flex;
  gap: 0.75rem;
}

/* 移动端适配 */
@media (max-width: 640px) {
  .tutorial-hub {
    max-height: 95vh;
    border-radius: 16px 16px 0 0;
    margin-top: auto;
  }

  .tutorial-hub__header {
    padding: 1.25rem 1rem 0.75rem;
  }

  .tutorial-hub__progress {
    padding: 0.75rem 1rem;
  }

  .tutorial-hub__content {
    padding: 0.75rem 1rem;
  }

  .module-grid {
    grid-template-columns: 1fr;
  }

  .tutorial-hub__footer {
    flex-direction: column;
    gap: 0.75rem;
    padding: 1rem;
  }

  .footer-quick-select,
  .footer-actions {
    width: 100%;
    justify-content: center;
  }

  .footer-actions .n-button {
    flex: 1;
  }
}
</style>
