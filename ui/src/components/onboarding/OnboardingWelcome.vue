<script setup lang="ts">
/**
 * OnboardingWelcome - 欢迎页面组件
 * 展示品牌信息和功能亮点，提供快速开始和自主探索选项
 */
import { NButton, NIcon } from 'naive-ui'
import { Message, User, Dice } from '@vicons/tabler'
import { useOnboardingStore } from '@/stores/onboarding'
import { RECOMMENDED_MODULES } from '@/config/tutorialModules'

const onboarding = useOnboardingStore()

const handleQuickStart = () => {
  onboarding.selectRecommendedModules(RECOMMENDED_MODULES)
  onboarding.startSelectedModules()
}

const handleExplore = () => {
  onboarding.openTutorialHub()
}

const handleSkip = () => {
  onboarding.skip()
}

const features = [
  { icon: Message, text: '沉浸式角色扮演聊天' },
  { icon: User, text: '灵活的角色身份系统' },
  { icon: Dice, text: '内置骰子与跑团工具' },
]
</script>

<template>
  <div class="onboarding-welcome-overlay" @click.self="handleSkip">
    <div class="onboarding-welcome">
      <!-- Logo -->
      <div class="welcome-logo">
        <img src="@/assets/logo.svg" alt="SealChat" onerror="this.style.display='none'" />
        <span class="welcome-logo-fallback">🦭</span>
      </div>

      <!-- 标题 -->
      <h1 class="welcome-title">欢迎使用 SealChat</h1>
      <p class="welcome-subtitle">让我们花几分钟了解平台的核心功能</p>

      <!-- 功能亮点 -->
      <div class="welcome-features">
        <div v-for="feature in features" :key="feature.text" class="welcome-feature">
          <n-icon :component="feature.icon" size="24" />
          <span>{{ feature.text }}</span>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="welcome-actions">
        <n-button type="primary" size="large" @click="handleQuickStart">
          🚀 快速开始
        </n-button>
        <n-button size="large" @click="handleExplore">
          📚 自主探索
        </n-button>
      </div>

      <button class="welcome-skip" @click="handleSkip">
        我已了解，跳过引导
      </button>
    </div>
  </div>
</template>

<style scoped>
.onboarding-welcome-overlay {
  position: fixed;
  inset: 0;
  z-index: 10000;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
}

.onboarding-welcome {
  background: var(--sc-bg-elevated, #fff);
  border-radius: 16px;
  padding: 2.5rem 2rem;
  max-width: 420px;
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

.welcome-logo {
  margin-bottom: 1.5rem;
}

.welcome-logo img {
  width: 72px;
  height: 72px;
}

.welcome-logo-fallback {
  font-size: 56px;
  display: block;
}

.welcome-title {
  font-size: 1.75rem;
  font-weight: 700;
  margin: 0 0 0.5rem;
  color: var(--sc-text-primary, #1f2937);
}

.welcome-subtitle {
  font-size: 1rem;
  color: var(--sc-text-secondary, #6b7280);
  margin: 0 0 2rem;
}

.welcome-features {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  margin-bottom: 2rem;
  text-align: left;
}

.welcome-feature {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem 1rem;
  background: var(--sc-bg-surface, #f9fafb);
  border-radius: 8px;
  color: var(--sc-text-primary, #374151);
}

.welcome-feature .n-icon {
  color: var(--primary-color, #3388de);
}

.welcome-actions {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  margin-bottom: 1.5rem;
}

.welcome-actions .n-button {
  width: 100%;
}

.welcome-skip {
  background: none;
  border: none;
  color: var(--sc-text-secondary, #9ca3af);
  font-size: 0.875rem;
  cursor: pointer;
  padding: 0.5rem;
  transition: color 0.15s;
}

.welcome-skip:hover {
  color: var(--sc-text-primary, #6b7280);
  text-decoration: underline;
}

/* 移动端适配 */
@media (max-width: 480px) {
  .onboarding-welcome {
    padding: 2rem 1.5rem;
    margin: 1rem;
  }

  .welcome-title {
    font-size: 1.5rem;
  }
}
</style>
