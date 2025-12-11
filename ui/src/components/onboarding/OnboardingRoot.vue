<script setup lang="ts">
/**
 * OnboardingRoot - 引导系统根组件
 * 管理不同引导状态下的组件显示
 */
import { computed, onMounted } from 'vue'
import { useOnboardingStore } from '@/stores/onboarding'
import { TUTORIAL_MODULES } from '@/config/tutorialModules'
import OnboardingWelcome from './OnboardingWelcome.vue'
import OnboardingTutorialHub from './OnboardingTutorialHub.vue'
import OnboardingSpotlight from './OnboardingSpotlight.vue'
import OnboardingResumePrompt from './OnboardingResumePrompt.vue'

const onboarding = useOnboardingStore()

// 注册模块列表
onMounted(() => {
  onboarding.registerModules(TUTORIAL_MODULES)
})

const showWelcome = computed(() => onboarding.showWelcome)
const showTutorialHub = computed(() => onboarding.isActive && onboarding.showTutorialHub)
const showSpotlight = computed(() => onboarding.showSpotlight)
const showResumePrompt = computed(() => onboarding.showResumePrompt)
</script>

<template>
  <Teleport to="body">
    <Transition name="fade">
      <OnboardingWelcome v-if="showWelcome" />
    </Transition>
    <Transition name="fade">
      <OnboardingTutorialHub v-if="showTutorialHub" />
    </Transition>
    <Transition name="fade">
      <OnboardingSpotlight v-if="showSpotlight" />
    </Transition>
    <Transition name="fade">
      <OnboardingResumePrompt v-if="showResumePrompt" />
    </Transition>
  </Teleport>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.25s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
