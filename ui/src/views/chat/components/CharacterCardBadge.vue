<script setup lang="ts">
import { computed } from 'vue';
import { useCharacterCardStore } from '@/stores/characterCard';
import { useDisplayStore } from '@/stores/display';
import { useChatStore } from '@/stores/chat';
import { renderCardTemplate, getWorldCardTemplate } from '@/utils/characterCardTemplate';

const props = defineProps<{
  identityId?: string;
  identityColor?: string;
}>();

const cardStore = useCharacterCardStore();
const displayStore = useDisplayStore();
const chatStore = useChatStore();

const boundCardId = computed(() => {
  if (!props.identityId) return '';
  return cardStore.getBoundCardId(props.identityId) || '';
});

const card = computed(() => {
  if (!boundCardId.value) return null;
  return cardStore.getCardById(boundCardId.value);
});

const template = computed(() => {
  const worldId = chatStore.currentWorldId;
  return displayStore.settings.characterCardBadgeTemplateByWorld?.[worldId] ?? getWorldCardTemplate(worldId);
});

const renderedContent = computed(() => {
  const channelId = chatStore.curChannel?.id || '';
  let attrs: Record<string, any> | undefined;
  if (channelId && boundCardId.value) {
    const activeIdentityId = chatStore.getActiveIdentityId(channelId);
    if (activeIdentityId && props.identityId && activeIdentityId === props.identityId) {
      attrs = cardStore.activeCards[channelId]?.attrs;
    } else {
      const activeId = cardStore.getActiveCardId(channelId);
      if (activeId && activeId === boundCardId.value) {
        attrs = cardStore.activeCards[channelId]?.attrs;
      }
    }
  }
  attrs = attrs || card.value?.attrs;
  if (!attrs) return '';
  return renderCardTemplate(template.value, attrs);
});

const isVisible = computed(() => {
  return displayStore.settings.characterCardBadgeEnabled && !!renderedContent.value;
});

const badgeStyle = computed(() => {
  if (!props.identityColor) return {};
  return {
    backgroundColor: `${props.identityColor}12`,
    color: props.identityColor,
    borderColor: `${props.identityColor}33`,
  };
});
</script>

<template>
  <span
    v-if="isVisible"
    class="character-card-badge"
    :style="badgeStyle"
    v-html="renderedContent"
  ></span>
</template>

<style scoped>
.character-card-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.3em;
  font-size: 0.68em;
  line-height: 1.2;
  padding: 0.08em 0.36em;
  border-radius: 6px;
  border: 1px solid rgba(128, 128, 128, 0.2);
  margin-left: 0.5em;
  vertical-align: middle;
  white-space: nowrap;
}
</style>
