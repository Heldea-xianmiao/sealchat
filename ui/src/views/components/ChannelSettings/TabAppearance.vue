<script lang="ts" setup>
import { ref, computed, watch, onUnmounted, type PropType } from 'vue';
import {
  NButton,
  NSpace,
  NSlider,
  NRadioGroup,
  NRadio,
  NUpload,
  NIcon,
  NColorPicker,
  NSwitch,
  NCard,
  NSpin,
  useMessage,
  type UploadFileInfo,
} from 'naive-ui';
import { Photo as ImageIcon, Trash } from '@vicons/tabler';
import VueCropper from 'vue-cropperjs';
import 'cropperjs/dist/cropper.css';
import { compressImage } from '@/composables/useImageCompressor';
import { uploadImageAttachment } from '@/views/chat/composables/useAttachmentUploader';
import { useChatStore } from '@/stores/chat';
import { resolveAttachmentUrl } from '@/composables/useAttachmentResolver';
import type { SChannel, ChannelBackgroundSettings } from '@/types';

const props = defineProps({
  channel: {
    type: Object as PropType<SChannel>,
  },
});

const emit = defineEmits<{
  (e: 'update'): void;
}>();

const message = useMessage();
const chat = useChatStore();
const defaultSettings: ChannelBackgroundSettings = {
  mode: 'cover',
  opacity: 30,
  blur: 0,
  brightness: 100,
  overlayColor: undefined,
  overlayOpacity: 0,
};

const parseSettings = (input?: ChannelBackgroundSettings | string): ChannelBackgroundSettings => {
  if (!input) return { ...defaultSettings };
  if (typeof input !== 'string') {
    return { ...defaultSettings, ...input };
  }
  try {
    return { ...defaultSettings, ...JSON.parse(input) };
  } catch {
    return { ...defaultSettings };
  }
};

const backgroundAttachmentId = ref<string>('');
const settings = ref<ChannelBackgroundSettings>({ ...defaultSettings });
const saving = ref(false);

const cropperVisible = ref(false);
const cropperFile = ref<File | null>(null);
const cropperRef = ref<InstanceType<typeof VueCropper> | null>(null);
const cropperImageUrl = ref('');
const cropperProcessing = ref(false);

const enableOverlay = ref(false);

const backgroundUrl = computed(() => {
  if (!backgroundAttachmentId.value) return '';
  return resolveAttachmentUrl(backgroundAttachmentId.value);
});

const previewStyle = computed(() => {
  if (!backgroundUrl.value) return {};
  const s = settings.value;
  let bgSize = 'cover';
  let bgRepeat = 'no-repeat';
  let bgPosition = 'center';
  switch (s.mode) {
    case 'contain':
      bgSize = 'contain';
      break;
    case 'tile':
      bgSize = 'auto';
      bgRepeat = 'repeat';
      break;
    case 'center':
      bgSize = 'auto';
      bgPosition = 'center';
      break;
  }
  return {
    backgroundImage: `url(${backgroundUrl.value})`,
    backgroundSize: bgSize,
    backgroundRepeat: bgRepeat,
    backgroundPosition: bgPosition,
    opacity: s.opacity / 100,
    filter: `blur(${s.blur}px) brightness(${s.brightness}%)`,
  };
});

const overlayStyle = computed(() => {
  if (!enableOverlay.value || !settings.value.overlayColor) return {};
  return {
    backgroundColor: settings.value.overlayColor,
    opacity: (settings.value.overlayOpacity ?? 0) / 100,
  };
});

watch(
  () => props.channel,
  (ch) => {
    if (ch) {
      backgroundAttachmentId.value = ch.backgroundAttachmentId || '';
      settings.value = parseSettings(ch.backgroundSettings);
      enableOverlay.value = !!settings.value.overlayColor && (settings.value.overlayOpacity ?? 0) > 0;
    }
  },
  { immediate: true }
);

const handleFileChange = (options: { file: UploadFileInfo }) => {
  const file = options.file.file;
  if (!file) return;
  cropperFile.value = file;
  const reader = new FileReader();
  reader.onload = (e) => {
    cropperImageUrl.value = e.target?.result as string;
    cropperVisible.value = true;
  };
  reader.readAsDataURL(file);
};

const handleCropConfirm = async () => {
  if (!cropperRef.value) return;
  cropperProcessing.value = true;
  try {
    const cropper = cropperRef.value.cropper;
    if (!cropper) throw new Error('Cropper not initialized');
    const canvas = cropper.getCroppedCanvas({
      maxWidth: 1920,
      maxHeight: 1080,
      imageSmoothingEnabled: true,
      imageSmoothingQuality: 'high',
    });
    const blob = await new Promise<Blob>((resolve, reject) => {
      canvas.toBlob((b) => {
        if (b) resolve(b);
        else reject(new Error('Failed to create blob'));
      }, 'image/png');
    });
    const pngFile = new File([blob], 'background.png', { type: 'image/png' });
    const compressed = await compressImage(pngFile, { maxWidth: 1920, maxHeight: 1080 });
    const result = await uploadImageAttachment(compressed, { channelId: props.channel?.id });
    let attachId = result.attachmentId || '';
    if (attachId.startsWith('id:')) {
      attachId = attachId.slice(3);
    }
    backgroundAttachmentId.value = attachId;
    cropperVisible.value = false;
    cropperFile.value = null;
    cropperImageUrl.value = '';
  } catch (err: any) {
    message.error(err?.message || '上传失败');
  } finally {
    cropperProcessing.value = false;
  }
};

const handleCropCancel = () => {
  cropperVisible.value = false;
  cropperFile.value = null;
  cropperImageUrl.value = '';
};

const removeBackground = () => {
  backgroundAttachmentId.value = '';
};

const resetSettings = () => {
  settings.value = { ...defaultSettings };
  enableOverlay.value = false;
};

const handleSave = async () => {
  if (!props.channel?.id) return;
  saving.value = true;
  try {
    const finalSettings: ChannelBackgroundSettings = { ...settings.value };
    if (!enableOverlay.value) {
      finalSettings.overlayColor = undefined;
      finalSettings.overlayOpacity = 0;
    }
    await chat.channelBackgroundEdit(props.channel.id, {
      backgroundAttachmentId: backgroundAttachmentId.value,
      backgroundSettings: JSON.stringify(finalSettings),
    });
    message.success('保存成功');
    emit('update');
  } catch (err: any) {
    message.error(err?.message || '保存失败');
  } finally {
    saving.value = false;
  }
};

onUnmounted(() => {
  cropperImageUrl.value = '';
});
</script>

<template>
  <div class="tab-appearance">
    <div class="section">
      <h4 class="section-title">背景图片</h4>
      <div class="upload-area">
        <div v-if="backgroundUrl" class="current-bg">
          <img :src="backgroundUrl" alt="当前背景" class="bg-thumb" />
          <NButton size="small" type="error" quaternary @click="removeBackground">
            <template #icon><NIcon :component="Trash" /></template>
            移除
          </NButton>
        </div>
        <NUpload
          accept="image/*"
          :show-file-list="false"
          :custom-request="() => {}"
          @change="handleFileChange"
        >
          <NButton type="primary" size="small">
            <template #icon><NIcon :component="ImageIcon" /></template>
            {{ backgroundUrl ? '更换图片' : '上传图片' }}
          </NButton>
        </NUpload>
      </div>
    </div>

    <div class="section" v-if="backgroundUrl">
      <h4 class="section-title">显示设置</h4>
      <div class="settings-grid">
        <div class="setting-row">
          <span class="setting-label">显示模式</span>
          <NRadioGroup v-model:value="settings.mode" size="small">
            <NRadio value="cover">铺满</NRadio>
            <NRadio value="contain">适应</NRadio>
            <NRadio value="tile">平铺</NRadio>
            <NRadio value="center">居中</NRadio>
          </NRadioGroup>
        </div>
        <div class="setting-row">
          <span class="setting-label">透明度 {{ settings.opacity }}%</span>
          <NSlider v-model:value="settings.opacity" :min="5" :max="100" :step="5" />
        </div>
        <div class="setting-row">
          <span class="setting-label">模糊 {{ settings.blur }}px</span>
          <NSlider v-model:value="settings.blur" :min="0" :max="20" :step="1" />
        </div>
        <div class="setting-row">
          <span class="setting-label">亮度 {{ settings.brightness }}%</span>
          <NSlider v-model:value="settings.brightness" :min="50" :max="150" :step="5" />
        </div>
        <div class="setting-row">
          <span class="setting-label">颜色叠加</span>
          <NSwitch v-model:value="enableOverlay" size="small" />
        </div>
        <template v-if="enableOverlay">
          <div class="setting-row">
            <span class="setting-label">叠加颜色</span>
            <NColorPicker v-model:value="settings.overlayColor" :show-alpha="false" size="small" />
          </div>
          <div class="setting-row">
            <span class="setting-label">叠加透明度 {{ settings.overlayOpacity ?? 0 }}%</span>
            <NSlider v-model:value="settings.overlayOpacity" :min="0" :max="100" :step="5" />
          </div>
        </template>
      </div>
      <NButton size="small" quaternary @click="resetSettings">重置为默认</NButton>
    </div>

    <div class="section" v-if="backgroundUrl">
      <h4 class="section-title">预览</h4>
      <NCard class="preview-card">
        <div class="preview-container">
          <div class="preview-bg" :style="previewStyle"></div>
          <div class="preview-overlay" :style="overlayStyle"></div>
          <div class="preview-content">
            <div class="preview-msg preview-msg--other">
              <div class="preview-avatar">A</div>
              <div class="preview-bubble">这是一条示例消息</div>
            </div>
            <div class="preview-msg preview-msg--self">
              <div class="preview-bubble">我的回复内容</div>
              <div class="preview-avatar">我</div>
            </div>
          </div>
        </div>
      </NCard>
    </div>

    <div class="actions">
      <NButton type="primary" :loading="saving" @click="handleSave">保存外观设置</NButton>
    </div>

    <n-modal v-model:show="cropperVisible" preset="card" title="裁剪背景图" style="max-width: 600px;">
      <div class="cropper-wrapper" v-if="cropperImageUrl">
        <VueCropper
          ref="cropperRef"
          :src="cropperImageUrl"
          :aspect-ratio="16 / 9"
          :view-mode="1"
          drag-mode="move"
          :auto-crop-area="0.9"
          :background="true"
          :guides="true"
          class="cropper-instance"
        />
      </div>
      <template #footer>
        <NSpace justify="end">
          <NButton @click="handleCropCancel">取消</NButton>
          <NButton type="primary" :loading="cropperProcessing" @click="handleCropConfirm">确认</NButton>
        </NSpace>
      </template>
    </n-modal>
  </div>
</template>

<style lang="scss" scoped>
.tab-appearance {
  padding: 1rem 0;
}

.section {
  margin-bottom: 1.5rem;
}

.section-title {
  font-size: 0.95rem;
  font-weight: 600;
  margin-bottom: 0.75rem;
  color: var(--n-text-color);
}

.upload-area {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.current-bg {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.bg-thumb {
  width: 80px;
  height: 45px;
  object-fit: cover;
  border-radius: 4px;
  border: 1px solid var(--n-border-color);
}

.settings-grid {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  margin-bottom: 0.5rem;
}

.setting-row {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.setting-label {
  min-width: 140px;
  font-size: 0.875rem;
  color: var(--n-text-color-3);
}

.preview-card {
  padding: 0;
}

.preview-container {
  position: relative;
  height: 180px;
  border-radius: 4px;
  overflow: hidden;
  background: #1a1a1a;
}

.preview-bg {
  position: absolute;
  inset: 0;
  z-index: 1;
}

.preview-overlay {
  position: absolute;
  inset: 0;
  z-index: 2;
  pointer-events: none;
}

.preview-content {
  position: relative;
  z-index: 3;
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.preview-msg {
  display: flex;
  align-items: flex-start;
  gap: 0.5rem;

  &--self {
    justify-content: flex-end;
  }
}

.preview-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.75rem;
  color: #fff;
  flex-shrink: 0;
}

.preview-bubble {
  background: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(4px);
  padding: 0.5rem 0.75rem;
  border-radius: 8px;
  color: #fff;
  font-size: 0.875rem;
  max-width: 200px;
}

.actions {
  display: flex;
  justify-content: flex-end;
  padding-top: 1rem;
  border-top: 1px solid var(--n-border-color);
}

.cropper-wrapper {
  width: 100%;
  height: 350px;
  background: #333;
  border-radius: 4px;
  overflow: hidden;
}

.cropper-instance {
  width: 100%;
  height: 100%;
}
</style>
