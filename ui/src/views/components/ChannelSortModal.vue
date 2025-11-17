<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { useMessage } from 'naive-ui';
import { useChatStore } from '@/stores/chat';
import type { SChannel } from '@/types';

interface SortNode {
  id: string;
  name: string;
  parentId: string;
  sortOrder: number;
  note?: string;
  membersCount?: number;
  permType?: string;
  children: SortNode[];
}

type SortRow =
  | { type: 'node'; key: string; node: SortNode; parentId: string; depth: number; index: number }
  | { type: 'tail'; key: string; parentId: string; depth: number };

const props = defineProps<{ show: boolean }>();
const emit = defineEmits<{ (e: 'update:show', value: boolean): void }>();

const chat = useChatStore();
const message = useMessage();

const visible = computed({
  get: () => props.show,
  set: (value: boolean) => emit('update:show', value),
});

const treeData = ref<SortNode[]>([]);
const originalOrders = ref<Record<string, string[]>>({});
const draggingId = ref<string | null>(null);
const dragOverKey = ref<string | null>(null);
const saving = ref(false);

watch(
  () => props.show,
  (val) => {
    if (val) {
      initData();
    } else {
      resetDragState();
    }
  },
  { immediate: true },
);

const initData = () => {
  const list = normalizeChannels(chat.channelTree as SChannel[] ?? [], '');
  treeData.value = list;
  originalOrders.value = {};
  captureOriginalOrders(list, '');
};

const normalizeChannels = (channels: SChannel[] | undefined, parentId: string): SortNode[] => {
  if (!Array.isArray(channels)) return [];
  return channels
    .filter((item) => !item.isPrivate && item.permType !== 'private')
    .map((item) => ({
      id: item.id,
      name: item.name || '未命名频道',
      parentId,
      sortOrder: item.sortOrder ?? 0,
      note: item.note ?? '',
      membersCount: item.membersCount,
      permType: item.permType || 'public',
      children: normalizeChannels(item.children as SChannel[] | undefined, item.id),
    }));
};

const captureOriginalOrders = (nodes: SortNode[], parentId: string) => {
  originalOrders.value[parentId] = nodes.map((node) => node.id);
  nodes.forEach((node) => {
    if (node.children?.length) {
      captureOriginalOrders(node.children, node.id);
    } else {
      originalOrders.value[node.id] = [];
    }
  });
};

const rows = computed<SortRow[]>(() => {
  const flatten = (nodes: SortNode[], depth: number, parentId: string): SortRow[] => {
    const current: SortRow[] = [];
    nodes.forEach((node, index) => {
      current.push({ type: 'node', key: node.id, node, parentId, depth, index });
      if (node.children?.length) {
        current.push(...flatten(node.children, depth + 1, node.id));
      }
    });
    if (nodes.length > 0) {
      const tailKey = `tail-${parentId || 'root'}`;
      current.push({ type: 'tail', key: tailKey, parentId, depth });
    }
    return current;
  };
  return flatten(treeData.value, 0, '');
});

const isDirty = computed(() => checkDirty(treeData.value, ''));

const checkDirty = (nodes: SortNode[], parentId: string): boolean => {
  if (!nodes.length) return false;
  const currentIds = nodes.map((node) => node.id);
  const originalIds = originalOrders.value[parentId] || [];
  if (!arraysEqual(currentIds, originalIds)) {
    return true;
  }
  return nodes.some((node) => (node.children?.length ? checkDirty(node.children, node.id) : false));
};

const arraysEqual = (a: string[], b: string[]) => {
  if (a.length !== b.length) return false;
  return a.every((item, index) => item === b[index]);
};

const resetDragState = () => {
  draggingId.value = null;
  dragOverKey.value = null;
};

const handleDragStart = (row: SortRow) => {
  if (row.type !== 'node') return;
  draggingId.value = row.node.id;
};

const handleDragEnter = (row: SortRow) => {
  if (!draggingId.value) return;
  dragOverKey.value = row.key;
};

const handleDrop = (row: SortRow) => {
  if (!draggingId.value) return;
  if (row.type === 'node') {
    moveNodeBefore(row.parentId, row.index, false);
  } else if (row.type === 'tail') {
    moveNodeBefore(row.parentId, Number.POSITIVE_INFINITY, true);
  }
};

const moveNodeBefore = (parentId: string, targetIndex: number, isTail = false) => {
  const dragged = draggingId.value;
  if (!dragged) return;
  const sourceInfo = findNodeInfo(dragged);
  if (!sourceInfo) return;
  if (sourceInfo.parentId !== parentId) {
    message.warning('暂不支持跨层级拖动，请在同一层内排序');
    return;
  }
  const list = getListByParent(parentId);
  if (!list) return;
  const from = sourceInfo.index;
  const cappedTarget = isTail ? list.length : targetIndex;
  const [item] = list.splice(from, 1);
  let insertIndex = cappedTarget;
  if (insertIndex > list.length) {
    insertIndex = list.length;
  }
  if (from < insertIndex) {
    insertIndex -= 1;
    if (insertIndex < 0) insertIndex = 0;
  }
  list.splice(insertIndex, 0, item);
  dragOverKey.value = null;
};

const findNodeInfo = (
  nodeId: string,
  nodes: SortNode[] = treeData.value,
  parentId = '',
): { node: SortNode; parentId: string; index: number } | null => {
  for (let index = 0; index < nodes.length; index += 1) {
    const node = nodes[index];
    if (node.id === nodeId) {
      return { node, parentId, index };
    }
    if (node.children?.length) {
      const child = findNodeInfo(nodeId, node.children, node.id);
      if (child) return child;
    }
  }
  return null;
};

const getListByParent = (parentId: string): SortNode[] | null => {
  if (!parentId) return treeData.value;
  const parentInfo = findNodeInfo(parentId);
  return parentInfo?.node.children || null;
};

const closeModal = () => {
  visible.value = false;
};

interface SortUpdate {
  id: string;
  name: string;
  note?: string;
  permType?: string;
  sortOrder: number;
}

const generateSequentialOrders = (count: number) => {
  const base = count * 100;
  return Array.from({ length: count }, (_, index) => base - index * 100);
};

const collectUpdates = (nodes: SortNode[], parentId: string, bucket: SortUpdate[]) => {
  if (!nodes.length) return;
  const original = originalOrders.value[parentId] || [];
  const current = nodes.map((node) => node.id);
  const changed = !arraysEqual(current, original);
  if (changed) {
    const orders = generateSequentialOrders(nodes.length);
    nodes.forEach((node, index) => {
      const nextOrder = orders[index];
      if (node.sortOrder !== nextOrder) {
        bucket.push({
          id: node.id,
          sortOrder: nextOrder,
          name: node.name,
          note: node.note ?? '',
          permType: node.permType,
        });
        node.sortOrder = nextOrder;
      }
    });
  }
  nodes.forEach((node) => {
    if (node.children?.length) {
      collectUpdates(node.children, node.id, bucket);
    }
  });
};

const saveReorder = async () => {
  if (saving.value) return false;
  if (!isDirty.value) {
    message.info('顺序未发生变化');
    closeModal();
    return false;
  }
  const updates: SortUpdate[] = [];
  collectUpdates(treeData.value, '', updates);
  if (!updates.length) {
    message.info('顺序未发生变化');
    closeModal();
    return false;
  }
  saving.value = true;
  try {
    for (const update of updates) {
      await chat.channelInfoEdit(update.id, {
        sortOrder: update.sortOrder,
        name: update.name,
        note: update.note ?? '',
        permType: update.permType,
      });
    }
    if (chat.currentWorldId) {
      await chat.channelList(chat.currentWorldId, true);
    }
    message.success('频道排序已更新');
    closeModal();
  } catch (error: any) {
    message.error(error?.message || '保存排序失败');
  } finally {
    saving.value = false;
  }
  return false;
};

const refreshFromServer = async () => {
  if (!chat.currentWorldId) return;
  await chat.channelList(chat.currentWorldId, true);
  initData();
};
</script>

<template>
  <n-modal v-model:show="visible" preset="dialog" title="频道排序" :positive-text="'保存'" :negative-text="'取消'"
    :positive-button-props="{ disabled: !isDirty || !treeData.length, loading: saving }" @positive-click="saveReorder"
    @negative-click="closeModal">
    <div class="space-y-3">
      <n-alert type="info" title="拖动提示" :closable="false">
        拖动同一父级中的频道即可调整顺序，暂不支持跨层级拖动。
      </n-alert>
      <div class="flex gap-2">
        <n-button size="tiny" @click="refreshFromServer" :loading="saving">
          重新获取
        </n-button>
        <n-button size="tiny" tertiary @click="initData" :disabled="saving">
          恢复初始顺序
        </n-button>
      </div>
      <div v-if="!treeData.length" class="py-6">
        <n-empty description="暂无可排序的频道" />
      </div>
      <div v-else class="channel-sort-list" @dragover.prevent>
        <div v-for="row in rows" :key="row.key">
          <template v-if="row.type === 'node'">
            <div class="channel-sort-item" :class="{
              dragging: draggingId === row.node.id,
              'drop-target': dragOverKey === row.key,
            }" draggable="true"
              @dragstart="handleDragStart(row)" @dragend="resetDragState" @dragenter.prevent="handleDragEnter(row)"
              @dragover.prevent @drop.prevent="handleDrop(row)" :style="{ paddingLeft: `${row.depth * 20 + 12}px` }">
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-2">
                  <div class="font-medium">{{ row.node.name }}</div>
                  <n-tag size="small" v-if="row.node.permType === 'non-public'" type="warning" round>
                    非公开
                  </n-tag>
                </div>
                <div class="text-xs text-gray-500">
                  {{ row.node.membersCount ? `${row.node.membersCount}人` : '' }}
                </div>
              </div>
            </div>
          </template>
          <template v-else>
            <div class="channel-sort-dropzone" :class="{ 'drop-target': dragOverKey === row.key }"
              @dragenter.prevent="handleDragEnter(row)" @dragover.prevent @drop.prevent="handleDrop(row)"
              :style="{ paddingLeft: `${row.depth * 20 + 12}px` }">
              拖到这里放置到该分组末尾
            </div>
          </template>
        </div>
      </div>
    </div>
  </n-modal>
</template>

<style scoped lang="scss">
.channel-sort-list {
  max-height: 60vh;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
}

.channel-sort-item {
  border: 1px solid var(--n-border-color);
  border-radius: 0.375rem;
  padding: 0.5rem 0.75rem;
  background-color: var(--sc-bg-2);
  cursor: grab;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
}

.channel-sort-item.dragging {
  opacity: 0.6;
  box-shadow: 0 0 0 2px var(--n-primary-color-hover);
}

.channel-sort-dropzone {
  border: 1px dashed var(--n-border-color);
  border-radius: 0.375rem;
  padding: 0.35rem 0.75rem;
  color: var(--n-text-color-disabled);
  font-size: 0.85rem;
}

.drop-target {
  border-color: var(--n-primary-color);
  background-color: rgba(46, 164, 79, 0.1);
}
</style>
