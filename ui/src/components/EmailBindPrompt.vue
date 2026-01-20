<script setup lang="ts">
import { computed } from 'vue';
import { NModal, NButton, NIcon } from 'naive-ui';
import { Mail } from '@vicons/tabler';
import { useI18n } from 'vue-i18n';

const props = withDefaults(defineProps<{
  show: boolean;
}>(), {
  show: false,
});

const emit = defineEmits<{
  (event: 'update:show', value: boolean): void;
  (event: 'bind'): void;
  (event: 'dismiss'): void;
  (event: 'skip'): void;
}>();

const { t } = useI18n();

const showModal = computed({
  get: () => props.show,
  set: (val: boolean) => emit('update:show', val),
});

const handleBind = () => {
  emit('bind');
  showModal.value = false;
};

const handleDismiss = () => {
  emit('dismiss');
  showModal.value = false;
};

const handleSkip = () => {
  emit('skip');
  showModal.value = false;
};
</script>

<template>
  <n-modal
    v-model:show="showModal"
    preset="card"
    :bordered="false"
    :closable="true"
    :mask-closable="true"
    class="email-bind-modal"
    style="width: 360px; max-width: 90vw;"
    @close="handleSkip"
  >
    <div class="email-prompt">
      <div class="email-prompt__icon-wrap">
        <n-icon :component="Mail" size="48" color="var(--sc-primary, #3b82f6)" />
      </div>

      <h3 class="email-prompt__title">{{ t('emailBindPrompt.title') }}</h3>
      <p class="email-prompt__desc">{{ t('emailBindPrompt.description') }}</p>

      <div class="email-prompt__actions">
        <n-button type="primary" @click="handleBind">
          {{ t('emailBindPrompt.action') }}
        </n-button>
        <n-button quaternary @click="handleDismiss">
          {{ t('emailBindPrompt.dismiss') }}
        </n-button>
        <n-button text size="small" @click="handleSkip">
          {{ t('emailBindPrompt.later') }}
        </n-button>
      </div>
    </div>
  </n-modal>
</template>

<style scoped>
.email-prompt {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  padding: 1rem 0.5rem;
}

.email-prompt__icon-wrap {
  margin-bottom: 1.25rem;
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: var(--sc-bg-secondary, #f3f4f6);
  display: flex;
  align-items: center;
  justify-content: center;
}

.email-prompt__title {
  font-size: 1.25rem;
  font-weight: 600;
  margin: 0 0 0.5rem;
  color: var(--sc-text-primary, #1f2937);
}

.email-prompt__desc {
  font-size: 0.875rem;
  color: var(--sc-text-secondary, #6b7280);
  margin: 0 0 1.5rem;
  line-height: 1.5;
}

.email-prompt__actions {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  width: 100%;
}

.email-prompt__actions .n-button {
  width: 100%;
}
</style>
