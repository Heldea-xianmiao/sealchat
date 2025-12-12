<script setup lang="ts">
/**
 * OnboardingResumePrompt - 未完成引导继续提示
 * 当用户有未完成的引导时显示
 */
import { computed } from 'vue'
import { NButton } from 'naive-ui'
import { useOnboardingStore } from '@/stores/onboarding'

const onboarding = useOnboardingStore()

const currentModule = computed(() => onboarding.currentModuleConfig)

const handleResume = () => {
  onboarding.resumeOnboarding()
}

const handleRestart = () => {
  onboarding.openTutorialHub()
}

const handleSkip = () => {
  onboarding.skip()
}
</script>

<template>
  <div class="resume-prompt-overlay" @click.self="handleSkip">
    <div class="resume-prompt">
      <div class="resume-icon">📖</div>
      <h2 class="resume-title">继续学习？</h2>
      <p class="resume-desc">
        你有一个未完成的教程：<strong>{{ currentModule?.title }}</strong>
      </p>

      <div class="resume-actions">
        <n-button type="primary" @click="handleResume">
          继续学习
        </n-button>
        <n-button @click="handleRestart">
          重新选择
        </n-button>
      </div>

      <button class="resume-skip" @click="handleSkip">
        跳过引导
      </button>
    </div>
  </div>
</template>

<style scoped>
.resume-prompt-overlay {
  position: fixed;
  inset: 0;
  z-index: 10000;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
}

.resume-prompt {
  background: var(--sc-bg-elevated, #fff);
  border-radius: 16px;
  padding: 2rem;
  max-width: 360px;
  width: 90%;
  text-align: center;
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

.resume-icon {
  font-size: 48px;
  margin-bottom: 1rem;
}

.resume-title {
  margin: 0 0 0.5rem;
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--sc-text-primary, #1f2937);
}

.resume-desc {
  margin: 0 0 1.5rem;
  font-size: 1rem;
  color: var(--sc-text-secondary, #6b7280);
}

.resume-desc strong {
  color: var(--primary-color, #3388de);
}

.resume-actions {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  margin-bottom: 1.5rem;
}

.resume-actions .n-button {
  width: 100%;
}

.resume-skip {
  background: none;
  border: none;
  color: var(--sc-text-secondary, #9ca3af);
  font-size: 0.875rem;
  cursor: pointer;
  padding: 0.5rem;
  transition: color 0.15s;
}

.resume-skip:hover {
  color: var(--sc-text-primary, #6b7280);
  text-decoration: underline;
}
</style>
