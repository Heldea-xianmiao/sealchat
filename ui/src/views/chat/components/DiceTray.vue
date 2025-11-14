<template>
  <div class="dice-tray">
    <div class="dice-tray__header">
      <div>
        默认骰：<strong>{{ currentDefaultDice }}</strong>
      </div>
      <n-button v-if="canEditDefault" size="tiny" text type="primary" @click="modalVisible = true">
        修改
      </n-button>
    </div>
    <div class="dice-tray__body">
      <div class="dice-tray__column dice-tray__column--quick">
        <div class="dice-tray__section-title">快捷骰</div>
        <div class="dice-tray__quick-grid">
          <button
            v-for="faces in quickFaces"
            :key="faces"
            type="button"
            class="dice-tray__quick-btn"
            @click="handleQuickSelect(faces)"
          >
            <span>d{{ faces }}</span>
            <span v-if="quickSelections[faces]" class="dice-tray__quick-count">×{{ quickSelections[faces] }}</span>
          </button>
        </div>
        <div v-if="hasQuickSelection" class="dice-tray__quick-summary">
          <div class="dice-tray__quick-expression">{{ quickExpression }}</div>
          <div class="dice-tray__quick-tools">
            <span class="dice-tray__quick-total">共 {{ quickTotal }} 次</span>
            <n-button text size="tiny" @click="clearQuickSelection">清空</n-button>
          </div>
        </div>
      </div>
      <div class="dice-tray__column dice-tray__column--form">
        <div class="dice-tray__section-title">自定义</div>
        <div class="dice-tray__form">
          <n-form-item label="数量">
            <n-input-number v-model:value="count" :min="1" size="small" />
          </n-form-item>
          <n-form-item label="面数">
            <n-input-number v-model:value="sides" :min="1" size="small" />
          </n-form-item>
          <n-form-item label="修正">
            <n-input-number v-model:value="modifier" size="small" />
          </n-form-item>
          <n-form-item label="理由">
            <n-input v-model:value="reason" size="small" placeholder="可选，例如攻击" />
          </n-form-item>
          <div class="dice-tray__actions">
            <n-button size="small" :disabled="!canSubmit" @click="handleInsert">
              插入到输入框
            </n-button>
            <n-button type="primary" size="small" :disabled="!canSubmit" @click="handleRoll">
              立即掷骰
            </n-button>
          </div>
        </div>
      </div>
    </div>
    <div v-if="hasHistory" class="dice-tray__history">
      <div class="dice-tray__section-title">最近检定</div>
      <div class="dice-tray__history-list">
        <div v-for="item in displayedHistory" :key="item.id" class="dice-tray__history-entry">
          <button type="button" class="dice-tray__history-roll" @click="handleHistoryRoll(item)">
            <span class="dice-tray__history-label">{{ formatHistoryLabel(item.expr) }}</span>
          </button>
          <button
            type="button"
            class="dice-tray__history-fav"
            :class="{ 'is-active': item.favorite }"
            :aria-pressed="item.favorite"
            @click.stop="toggleFavorite(item.id)"
          >
            <span v-if="item.favorite">★</span>
            <span v-else>☆</span>
          </button>
        </div>
      </div>
    </div>
  </div>
  <n-modal
    v-model:show="modalVisible"
    preset="card"
    class="dice-settings-modal"
    :bordered="false"
    title="修改默认骰"
  >
    <n-form size="small" label-placement="left" :show-feedback="false">
      <n-form-item label="面数">
        <n-input v-model:value="defaultDiceInput" placeholder="例如 d20" />
      </n-form-item>
      <n-alert v-if="defaultDiceError" type="warning" :show-icon="false">
        {{ defaultDiceError }}
      </n-alert>
      <div class="dice-tray__settings-actions">
        <n-button @click="modalVisible = false">取消</n-button>
        <n-button type="primary" :disabled="!!defaultDiceError" @click="handleSaveDefault">
          保存
        </n-button>
      </div>
    </n-form>
  </n-modal>
</template>

<script setup lang="ts">
import { computed, ref, watch, onMounted } from 'vue';
import { ensureDefaultDiceExpr, isValidDefaultDiceExpr } from '@/utils/dice';

const props = withDefaults(defineProps<{
  defaultDice?: string
  canEditDefault?: boolean
}>(), {
  defaultDice: 'd20',
  canEditDefault: false,
});

const emit = defineEmits<{
  (event: 'insert', expr: string): void
  (event: 'roll', expr: string): void
  (event: 'update-default', expr: string): void
}>();

const quickFaces = [2, 4, 6, 8, 10, 12, 20, 100];
const quickSelections = ref<Record<number, number>>({});
const count = ref(1);
const sides = ref<number | null>(null);
const modifier = ref(0);
const reason = ref('');
const modalVisible = ref(false);
const defaultDiceInput = ref(ensureDefaultDiceExpr(props.defaultDice));

type DiceHistoryItem = {
  id: string;
  expr: string;
  favorite: boolean;
  timestamp: number;
};

const HISTORY_KEY = 'sealchat:dice-history';
const HISTORY_DISPLAY_LIMIT = 4;
const HISTORY_STORAGE_LIMIT = 12;

const historyItems = ref<DiceHistoryItem[]>([]);

const isClient = typeof window !== 'undefined';

const sortByTimestampDesc = (a: DiceHistoryItem, b: DiceHistoryItem) => b.timestamp - a.timestamp;

const persistHistory = () => {
  if (!isClient) return;
  try {
    window.localStorage.setItem(HISTORY_KEY, JSON.stringify(historyItems.value));
  } catch (error) {
    console.warn('无法保存骰子历史', error);
  }
};

const loadHistory = () => {
  if (!isClient) return;
  try {
    const raw = window.localStorage.getItem(HISTORY_KEY);
    if (!raw) return;
    const parsed = JSON.parse(raw);
    if (!Array.isArray(parsed)) return;
    historyItems.value = parsed
      .filter((item: any) => typeof item?.expr === 'string')
      .map((item: any): DiceHistoryItem => ({
        id: item.id || `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
        expr: item.expr,
        favorite: !!item.favorite,
        timestamp: typeof item.timestamp === 'number' ? item.timestamp : Date.now(),
      }))
      .sort(sortByTimestampDesc);
  } catch (error) {
    console.warn('无法加载骰子历史', error);
  }
};

onMounted(() => {
  loadHistory();
});

const pruneHistory = () => {
  const favorites = historyItems.value.filter((item) => item.favorite).sort(sortByTimestampDesc);
  const nonFavorites = historyItems.value.filter((item) => !item.favorite).sort(sortByTimestampDesc);
  const allowedNonFavorites = Math.max(0, HISTORY_STORAGE_LIMIT - favorites.length);
  const keptNonFavorites = nonFavorites.slice(0, allowedNonFavorites);
  historyItems.value = [...favorites, ...keptNonFavorites].sort(sortByTimestampDesc);
};

const recordHistory = (expr: string) => {
  const trimmed = expr.trim();
  if (!trimmed) return;
  const existingIndex = historyItems.value.findIndex((item) => item.expr === trimmed);
  const favorite = existingIndex !== -1 ? historyItems.value[existingIndex].favorite : false;
  if (existingIndex !== -1) {
    historyItems.value.splice(existingIndex, 1);
  }
  const newItem: DiceHistoryItem = {
    id: `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
    expr: trimmed,
    favorite,
    timestamp: Date.now(),
  };
  historyItems.value = [newItem, ...historyItems.value];
  pruneHistory();
  persistHistory();
};

const toggleFavorite = (itemId: string) => {
  historyItems.value = historyItems.value.map((item) =>
    item.id === itemId ? { ...item, favorite: !item.favorite } : item,
  );
  pruneHistory();
  persistHistory();
};

const displayedHistory = computed(() => {
  const favorites = historyItems.value.filter((item) => item.favorite).sort(sortByTimestampDesc);
  const nonFavorites = historyItems.value.filter((item) => !item.favorite).sort(sortByTimestampDesc);
  const remainingSlots = Math.max(0, HISTORY_DISPLAY_LIMIT - favorites.length);
  const recentNonFavorites = remainingSlots > 0 ? nonFavorites.slice(0, remainingSlots) : [];
  return [...favorites, ...recentNonFavorites];
});

const hasHistory = computed(() => displayedHistory.value.length > 0);

const formatHistoryLabel = (expr: string) => expr.replace(/^\.r/, 'r').replace(/\s+/g, ' ');

const currentDefaultDice = computed(() => ensureDefaultDiceExpr(props.defaultDice));

watch(() => props.defaultDice, (value) => {
  defaultDiceInput.value = ensureDefaultDiceExpr(value);
  if (!sides.value) {
    sides.value = parseInt(defaultDiceInput.value.slice(1), 10) || 20;
  }
}, { immediate: true });

const sanitizedReason = computed(() => reason.value.trim());

const quickExpression = computed(() => {
  const entries = Object.entries(quickSelections.value).filter(([, count]) => count > 0);
  if (!entries.length) {
    return '';
  }
  return entries
    .map(([faces, count]) => `.r${count}d${faces}`)
    .join(' + ');
});

const quickTotal = computed(() =>
  Object.values(quickSelections.value).reduce((sum, count) => sum + count, 0)
);

const hasQuickSelection = computed(() => quickTotal.value > 0);

const expression = computed(() => {
  if (!count.value || !sides.value) {
    return '';
  }
  const amount = Math.max(1, Math.floor(count.value));
  const face = Math.max(1, Math.floor(sides.value));
  const parts = [`.r${amount}d${face}`];
  if (modifier.value) {
    const delta = Math.trunc(modifier.value);
    if (delta > 0) {
      parts.push(`+${delta}`);
    } else {
      parts.push(`${delta}`);
    }
  }
  if (sanitizedReason.value) {
    parts.push(`#${sanitizedReason.value}`);
  }
  return parts.join(' ');
});

const combinedExpression = computed(() => quickExpression.value || expression.value);

const canSubmit = computed(() => !!combinedExpression.value);

const handleQuickSelect = (faces: number) => {
  quickSelections.value = {
    ...quickSelections.value,
    [faces]: (quickSelections.value[faces] || 0) + 1,
  };
};

const clearQuickSelection = () => {
  quickSelections.value = {};
};

const handleInsert = () => {
  if (canSubmit.value && combinedExpression.value) {
    emit('insert', combinedExpression.value);
    if (hasQuickSelection.value) {
      clearQuickSelection();
    }
  }
};

const handleRoll = () => {
  if (canSubmit.value && combinedExpression.value) {
    emit('roll', combinedExpression.value);
    recordHistory(combinedExpression.value);
    if (hasQuickSelection.value) {
      clearQuickSelection();
    }
  }
};

const handleHistoryRoll = (item: DiceHistoryItem) => {
  emit('roll', item.expr);
  recordHistory(item.expr);
};

const defaultDiceError = computed(() => {
  if (!defaultDiceInput.value) {
    return '请输入默认骰，例如 d20';
  }
  if (!isValidDefaultDiceExpr(defaultDiceInput.value)) {
    return '格式不正确，示例：d20';
  }
  return '';
});

const handleSaveDefault = () => {
  if (defaultDiceError.value) {
    return;
  }
  emit('update-default', ensureDefaultDiceExpr(defaultDiceInput.value));
  modalVisible.value = false;
};
</script>

<style scoped>
.dice-tray {
  min-width: 320px;
  max-width: 480px;
  padding: 12px;
  background: var(--sc-bg-elevated, #fff);
  border: 1px solid var(--sc-border-strong, #e5e7eb);
  border-radius: 12px;
  color: var(--sc-fg-primary, #111);
}

.dice-tray__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
  font-size: 13px;
}

.dice-tray__body {
  display: flex;
  gap: 12px;
}

.dice-tray__column {
  flex: 1;
  padding: 8px;
  border-radius: 10px;
  background: var(--sc-bg-layer, #fafafa);
}

.dice-tray__column--quick {
  flex: 0 0 140px;
}

.dice-tray__section-title {
  font-size: 12px;
  color: var(--sc-fg-muted, #666);
  margin-bottom: 6px;
}

.dice-tray__quick-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 6px;
}

.dice-tray__quick-btn {
  position: relative;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border: 1px solid var(--sc-border-mute, #d1d5db);
  border-radius: 8px;
  padding: 0.3rem 0;
  font-size: 0.9rem;
  background: var(--sc-bg-layer, #fff);
  color: var(--sc-fg-primary, #111);
  transition: background 0.2s ease, color 0.2s ease;
}

.dice-tray__quick-btn:hover {
  background: rgba(15, 23, 42, 0.04);
}

.dice-tray__quick-count {
  position: absolute;
  top: -0.35rem;
  right: -0.35rem;
  font-size: 0.65rem;
  background: var(--sc-accent, #2563eb);
  color: #fff;
  border-radius: 999px;
  padding: 0.05rem 0.35rem;
  box-shadow: 0 2px 6px rgba(15, 23, 42, 0.2);
}

.dice-tray__quick-summary {
  margin-top: 0.5rem;
  padding: 0.4rem 0.5rem;
  border-radius: 6px;
  background: rgba(15, 23, 42, 0.04);
  font-size: 0.85rem;
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.dice-tray__quick-expression {
  word-break: break-all;
  font-family: var(--sc-code-font, 'SFMono-Regular', Menlo, Consolas, monospace);
}

.dice-tray__quick-tools {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.8rem;
}

.dice-tray__form :deep(.n-form-item) {
  margin-bottom: 8px;
}

.dice-tray__actions {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
  margin-top: 8px;
}

.dice-tray__settings-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 12px;
}

.dice-tray__history {
  margin-top: 0.75rem;
  padding-top: 0.75rem;
  border-top: 1px solid var(--sc-border-mute, #e2e8f0);
  color: var(--sc-fg-primary, #111);
}

.dice-tray__history-list {
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
}

.dice-tray__history-entry {
  display: flex;
  align-items: center;
  gap: 0.4rem;
}

.dice-tray__history-roll {
  flex: 1;
  border: 1px solid var(--sc-border-mute, #d1d5db);
  border-radius: 6px;
  background: var(--sc-bg-layer, #f8fafc);
  color: var(--sc-fg-primary, #111);
  padding: 0.35rem 0.5rem;
  font-size: 0.85rem;
  text-align: left;
  transition: background 0.2s ease, color 0.2s ease;
}

.dice-tray__history-roll:hover {
  background: rgba(15, 23, 42, 0.08);
}

.dice-tray__history-label {
  display: inline-block;
  width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.dice-tray__history-fav {
  width: 1.75rem;
  height: 1.75rem;
  border-radius: 999px;
  border: 1px solid var(--sc-border-mute, #d1d5db);
  background: transparent;
  color: var(--sc-fg-muted, #6b7280);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 1rem;
  transition: background 0.2s ease, color 0.2s ease, border-color 0.2s ease;
}

.dice-tray__history-fav.is-active {
  color: var(--sc-accent, #2563eb);
  border-color: currentColor;
  background: rgba(37, 99, 235, 0.08);
}

:global([data-display-palette='night']) .dice-tray {
  background: var(--sc-bg-elevated, #2a282a);
  border-color: var(--sc-border-strong, rgba(255, 255, 255, 0.12));
  color: var(--sc-fg-primary, #eee);
}

:global([data-display-palette='night']) .dice-tray__column {
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.08);
}

:global([data-display-palette='night']) .dice-tray__column--quick {
  background: rgba(255, 255, 255, 0.03);
}

:global([data-display-palette='night']) .dice-tray__column--form {
  background: rgba(255, 255, 255, 0.06);
}

:global([data-display-palette='night']) .dice-tray__quick-btn {
  border-color: rgba(255, 255, 255, 0.2);
  background: rgba(15, 23, 42, 0.35);
  color: var(--sc-fg-primary, #eee);
}

:global([data-display-palette='night']) .dice-tray__quick-btn:hover {
  background: rgba(255, 255, 255, 0.12);
}

:global([data-display-palette='night']) .dice-tray__quick-count {
  background: var(--sc-accent-night, #60a5fa);
  color: #0f172a;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.35);
}

:global([data-display-palette='night']) .dice-tray__quick-summary {
  background: rgba(255, 255, 255, 0.08);
}

:global([data-display-palette='night']) .dice-tray__history {
  border-top-color: rgba(255, 255, 255, 0.12);
  color: var(--sc-fg-primary, #f8fafc);
}

:global([data-display-palette='night']) .dice-tray__history-roll {
  border-color: rgba(255, 255, 255, 0.2);
  background: rgba(15, 23, 42, 0.35);
  color: var(--sc-fg-primary, #f8fafc);
}

:global([data-display-palette='night']) .dice-tray__history-roll:hover {
  background: rgba(255, 255, 255, 0.12);
}

:global([data-display-palette='night']) .dice-tray__history-fav {
  border-color: rgba(255, 255, 255, 0.2);
  color: rgba(226, 232, 240, 0.8);
}

:global([data-display-palette='night']) .dice-tray__history-fav.is-active {
  color: var(--sc-accent-night, #60a5fa);
  border-color: var(--sc-accent-night, #60a5fa);
  background: rgba(96, 165, 250, 0.15);
}

.dice-settings-modal :global(.n-card__content) {
  padding-top: 0;
}

.dice-settings-modal :global(.n-card) {
  background: var(--sc-bg-elevated, #fff);
  color: var(--sc-fg-primary, #111);
  max-width: 360px;
  width: min(360px, 90vw);
  margin: 0 auto;
}

::global([data-display-palette='night']) .dice-settings-modal :global(.n-card) {
  background: var(--sc-bg-elevated, #2a282a);
  color: #fff;
  border: 1px solid rgba(255, 255, 255, 0.12);
}
</style>
