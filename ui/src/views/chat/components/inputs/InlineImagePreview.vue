<script setup lang="ts">
import { computed } from 'vue';
import { NButton, NIcon, NSpin } from 'naive-ui';
import { X, AlertCircle } from '@vicons/tabler';

interface ImageInfo {
  status: 'uploading' | 'uploaded' | 'failed';
  previewUrl?: string;
  error?: string;
}

const props = defineProps<{
  images: Record<string, ImageInfo>;
}>();

const emit = defineEmits<{
  (event: 'remove', markerId: string): void;
}>();

const imageList = computed(() => {
  return Object.entries(props.images).map(([id, info]) => ({
    id,
    ...info,
  }));
});

const handleRemove = (markerId: string) => {
  emit('remove', markerId);
};
</script>

<template>
  <div v-if="imageList.length > 0" class="inline-image-preview">
    <div class="inline-image-preview__list">
      <div v-for="image in imageList" :key="image.id" class="inline-image-preview__item" :class="{
          'inline-image-preview__item--uploading': image.status === 'uploading',
          'inline-image-preview__item--failed': image.status === 'failed',
        }">
        <div class="inline-image-preview__thumbnail">
          <img v-if="image.previewUrl" :src="image.previewUrl" :alt="image.id" class="inline-image-preview__img" />
          <div v-else class="inline-image-preview__placeholder">
            <n-icon size="24" color="#9ca3af">
              <AlertCircle />
            </n-icon>
          </div>

          <!-- 上传中遮罩 -->
          <div v-if="image.status === 'uploading'" class="inline-image-preview__overlay">
            <n-spin size="small" stroke="#ffffff" />
          </div>

          <!-- 失败遮罩 -->
          <div v-if="image.status === 'failed'"
            class="inline-image-preview__overlay inline-image-preview__overlay--error">
            <n-icon size="20" color="#ffffff">
              <AlertCircle />
            </n-icon>
            <span class="inline-image-preview__error-text">{{ image.error || '上传失败' }}</span>
          </div>

          <!-- 删除按钮 -->
          <button class="inline-image-preview__remove" @click="handleRemove(image.id)" title="删除图片">
            <n-icon size="14">
              <X />
            </n-icon>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.inline-image-preview {
  margin-top: 0.5rem;
  padding: 0.5rem;
  background-color: #f9fafb;
  border-radius: 0.5rem;
}

.inline-image-preview__list {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.inline-image-preview__item {
  position: relative;
  width: 5rem;
  height: 5rem;
  border-radius: 0.5rem;
  overflow: hidden;
  border: 2px solid #e5e7eb;
  background-color: #ffffff;
  transition: border-color 0.2s ease;

  &:hover {
    border-color: #3b82f6;
  }

  &--uploading {
    border-color: #3b82f6;
  }

  &--failed {
    border-color: #ef4444;
  }
}

.inline-image-preview__thumbnail {
  position: relative;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.inline-image-preview__img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.inline-image-preview__placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  background-color: #f3f4f6;
}

.inline-image-preview__overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background-color: rgba(0, 0, 0, 0.5);
  gap: 0.25rem;

  &--error {
    background-color: rgba(239, 68, 68, 0.8);
  }
}

.inline-image-preview__error-text {
  font-size: 0.625rem;
  color: #ffffff;
  text-align: center;
  padding: 0 0.25rem;
}

.inline-image-preview__remove {
  position: absolute;
  top: 0.25rem;
  right: 0.25rem;
  width: 1.5rem;
  height: 1.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: rgba(0, 0, 0, 0.6);
  border: none;
  border-radius: 50%;
  color: #ffffff;
  cursor: pointer;
  opacity: 0;
  transition: opacity 0.2s ease, background-color 0.2s ease;

  &:hover {
    background-color: rgba(239, 68, 68, 0.9);
  }
}

.inline-image-preview__item:hover .inline-image-preview__remove {
  opacity: 1;
}
</style>
