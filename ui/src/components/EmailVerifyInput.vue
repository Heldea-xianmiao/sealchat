<script setup lang="ts">
import { ref, computed } from 'vue';
import { NInput, NButton, NInputGroup, useMessage } from 'naive-ui';

const props = defineProps<{
  modelValue: string;
  disabled?: boolean;
  loading?: boolean;
}>();

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void;
  (e: 'send'): Promise<boolean> | void;
}>();

const message = useMessage();
const sending = ref(false);
const countdown = ref(0);
let timer: ReturnType<typeof setInterval> | null = null;

const canSend = computed(() => !sending.value && countdown.value === 0 && !props.disabled);
const buttonText = computed(() => {
  if (sending.value) return '发送中...';
  if (countdown.value > 0) return `${countdown.value}s`;
  return '获取验证码';
});

const handleInput = (value: string) => {
  emit('update:modelValue', value);
};

const startCountdown = () => {
  countdown.value = 60;
  timer = setInterval(() => {
    countdown.value--;
    if (countdown.value <= 0) {
      clearInterval(timer!);
      timer = null;
    }
  }, 1000);
};

const handleSend = async () => {
  if (!canSend.value) return;

  sending.value = true;
  try {
    const result = await new Promise<boolean>((resolve) => {
      const maybePromise = emit('send');
      if (maybePromise && typeof maybePromise === 'object' && 'then' in maybePromise) {
        (maybePromise as Promise<boolean>).then(resolve).catch(() => resolve(false));
      } else {
        resolve(true);
      }
    });

    if (result !== false) {
      startCountdown();
    }
  } catch (e: any) {
    message.error(e?.response?.data?.error || '发送失败');
  } finally {
    sending.value = false;
  }
};
</script>

<template>
  <NInputGroup>
    <NInput
      :value="modelValue"
      @update:value="handleInput"
      placeholder="请输入验证码"
      :disabled="disabled || loading"
      maxlength="6"
      style="flex: 1"
    />
    <NButton
      type="primary"
      :disabled="!canSend"
      :loading="sending"
      @click="handleSend"
      style="width: 110px"
    >
      {{ buttonText }}
    </NButton>
  </NInputGroup>
</template>
