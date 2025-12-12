<script setup lang="ts">
/**
 * OnboardingSpotlight - 聚光灯高亮引导组件
 * 高亮目标元素并显示步骤提示卡片
 */
import { computed, ref, watch, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { NButton, NIcon } from 'naive-ui'
import { ChevronLeft, ChevronRight, X as CloseIcon } from '@vicons/tabler'
import { useOnboardingStore } from '@/stores/onboarding'

const onboarding = useOnboardingStore()

const currentModule = computed(() => onboarding.currentModuleConfig)
const currentStep = computed(() => onboarding.currentStepConfig)
const isFirstStep = computed(() => onboarding.isFirstStep)
const isLastStep = computed(() => onboarding.isLastStep)

// 目标元素位置
const targetRect = ref<DOMRect | null>(null)
const tooltipPlacement = ref<'top' | 'bottom' | 'left' | 'right' | 'center'>('center')

// 查找目标元素并计算位置
const updateTargetPosition = () => {
  const step = currentStep.value
  if (!step?.target) {
    targetRect.value = null
    tooltipPlacement.value = step?.placement || 'center'
    return
  }

  const el = document.querySelector(step.target) as HTMLElement | null
  if (!el) {
    targetRect.value = null
    tooltipPlacement.value = 'center'
    return
  }

  targetRect.value = el.getBoundingClientRect()
  tooltipPlacement.value = step.placement || 'bottom'

  // 滚动到可见区域
  el.scrollIntoView({ behavior: 'smooth', block: 'center' })
}

// 监听步骤变化
watch(
  [() => onboarding.progress.currentModule, () => onboarding.progress.currentStep],
  () => {
    nextTick(() => {
      updateTargetPosition()
    })
  },
  { immediate: true }
)

// 窗口大小变化时更新位置
let resizeObserver: ResizeObserver | null = null
onMounted(() => {
  window.addEventListener('resize', updateTargetPosition)
  window.addEventListener('scroll', updateTargetPosition, true)

  if (typeof ResizeObserver !== 'undefined') {
    resizeObserver = new ResizeObserver(updateTargetPosition)
    resizeObserver.observe(document.body)
  }
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', updateTargetPosition)
  window.removeEventListener('scroll', updateTargetPosition, true)
  resizeObserver?.disconnect()
})

// 遮罩样式（镂空目标区域）
const maskStyle = computed(() => {
  if (!targetRect.value) return {}

  const padding = 8
  const r = targetRect.value
  const x = r.left - padding
  const y = r.top - padding
  const w = r.width + padding * 2
  const h = r.height + padding * 2

  return {
    clipPath: `polygon(
      0% 0%, 0% 100%, 100% 100%, 100% 0%, 0% 0%,
      ${x}px ${y}px,
      ${x}px ${y + h}px,
      ${x + w}px ${y + h}px,
      ${x + w}px ${y}px,
      ${x}px ${y}px
    )`,
  }
})

// 高亮环样式
const ringStyle = computed(() => {
  if (!targetRect.value) return { display: 'none' }

  const padding = 6
  const r = targetRect.value
  return {
    left: `${r.left - padding}px`,
    top: `${r.top - padding}px`,
    width: `${r.width + padding * 2}px`,
    height: `${r.height + padding * 2}px`,
  }
})

// 提示卡片样式
const tooltipStyle = computed(() => {
  const gap = 16
  const tooltipWidth = 320
  const tooltipHeight = 200 // 估算高度

  if (!targetRect.value || tooltipPlacement.value === 'center') {
    return {
      left: '50%',
      top: '50%',
      transform: 'translate(-50%, -50%)',
      maxWidth: `${tooltipWidth}px`,
    }
  }

  const r = targetRect.value
  const vw = window.innerWidth
  const vh = window.innerHeight
  let left = 0
  let top = 0

  switch (tooltipPlacement.value) {
    case 'top':
      left = Math.min(vw - tooltipWidth - 16, Math.max(16, r.left + r.width / 2 - tooltipWidth / 2))
      top = r.top - tooltipHeight - gap
      if (top < 16) {
        top = r.bottom + gap
      }
      break
    case 'bottom':
      left = Math.min(vw - tooltipWidth - 16, Math.max(16, r.left + r.width / 2 - tooltipWidth / 2))
      top = r.bottom + gap
      if (top + tooltipHeight > vh - 16) {
        top = r.top - tooltipHeight - gap
      }
      break
    case 'left':
      left = r.left - tooltipWidth - gap
      top = Math.min(vh - tooltipHeight - 16, Math.max(16, r.top + r.height / 2 - tooltipHeight / 2))
      if (left < 16) {
        left = r.right + gap
      }
      break
    case 'right':
      left = r.right + gap
      top = Math.min(vh - tooltipHeight - 16, Math.max(16, r.top + r.height / 2 - tooltipHeight / 2))
      if (left + tooltipWidth > vw - 16) {
        left = r.left - tooltipWidth - gap
      }
      break
  }

  return {
    left: `${left}px`,
    top: `${top}px`,
    maxWidth: `${tooltipWidth}px`,
  }
})

const handlePrev = () => {
  onboarding.prevStep()
}

const handleNext = () => {
  onboarding.nextStep()
}

const handleSkip = () => {
  onboarding.skip()
}

const handleOpenHub = () => {
  onboarding.openTutorialHub()
}

// 进度信息
const progressText = computed(() => {
  const module = currentModule.value
  if (!module) return ''
  return `${onboarding.progress.currentStep + 1} / ${module.steps.length}`
})
</script>

<template>
  <div class="spotlight-overlay">
    <!-- 半透明遮罩（镂空目标区域） -->
    <div class="spotlight-mask" :style="maskStyle" />

    <!-- 高亮边框 -->
    <div v-if="targetRect" class="spotlight-ring" :style="ringStyle" />

    <!-- 提示卡片 -->
    <div class="spotlight-tooltip" :style="tooltipStyle">
      <!-- 模块标题 -->
      <div class="tooltip-header">
        <span class="tooltip-module">{{ currentModule?.title }}</span>
        <span class="tooltip-progress">{{ progressText }}</span>
      </div>

      <!-- 步骤标题 -->
      <h3 class="tooltip-title">{{ currentStep?.title }}</h3>

      <!-- 步骤内容 -->
      <p class="tooltip-content">{{ currentStep?.content }}</p>

      <!-- 步骤图片（如果有） -->
      <img
        v-if="currentStep?.image"
        :src="currentStep.image"
        alt=""
        class="tooltip-image"
      />

      <!-- 操作按钮 -->
      <div class="tooltip-actions">
        <n-button size="small" quaternary @click="handleSkip">
          跳过
        </n-button>
        <n-button size="small" quaternary @click="handleOpenHub">
          教程列表
        </n-button>
        <div class="tooltip-nav">
          <n-button
            size="small"
            :disabled="isFirstStep"
            @click="handlePrev"
          >
            <template #icon>
              <n-icon :component="ChevronLeft" />
            </template>
          </n-button>
          <n-button
            type="primary"
            size="small"
            @click="handleNext"
          >
            {{ isLastStep ? '完成' : '下一步' }}
            <template #icon v-if="!isLastStep">
              <n-icon :component="ChevronRight" />
            </template>
          </n-button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.spotlight-overlay {
  position: fixed;
  inset: 0;
  z-index: 10000;
  pointer-events: none;
}

.spotlight-mask {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  pointer-events: auto;
}

.spotlight-ring {
  position: fixed;
  border: 3px solid var(--primary-color, #3388de);
  border-radius: 8px;
  pointer-events: none;
  animation: pulse 2s ease-in-out infinite;
  z-index: 10001;
}

@keyframes pulse {
  0%, 100% {
    box-shadow: 0 0 0 0 rgba(51, 136, 222, 0.4);
  }
  50% {
    box-shadow: 0 0 0 8px rgba(51, 136, 222, 0);
  }
}

.spotlight-tooltip {
  position: fixed;
  background: var(--sc-bg-elevated, #fff);
  border-radius: 12px;
  padding: 1.25rem;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.25);
  pointer-events: auto;
  animation: fadeSlideIn 0.25s ease-out;
  z-index: 10002;
}

@keyframes fadeSlideIn {
  from {
    opacity: 0;
    transform: translateY(8px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.tooltip-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
}

.tooltip-module {
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--primary-color, #3388de);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.tooltip-progress {
  font-size: 0.75rem;
  color: var(--sc-text-secondary, #9ca3af);
}

.tooltip-title {
  margin: 0 0 0.5rem;
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--sc-text-primary, #1f2937);
}

.tooltip-content {
  margin: 0 0 1rem;
  font-size: 0.9375rem;
  line-height: 1.6;
  color: var(--sc-text-secondary, #4b5563);
  white-space: pre-line;
}

.tooltip-image {
  max-width: 100%;
  border-radius: 8px;
  margin-bottom: 1rem;
}

.tooltip-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 0.5rem;
}

.tooltip-nav {
  display: flex;
  gap: 0.5rem;
}

/* 移动端适配 */
@media (max-width: 640px) {
  .spotlight-tooltip {
    position: fixed !important;
    left: 16px !important;
    right: 16px !important;
    bottom: 16px !important;
    top: auto !important;
    transform: none !important;
    max-width: none !important;
    border-radius: 16px;
  }
}
</style>
