<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useChatStore } from '@/stores/chat';
import { useMessage } from 'naive-ui';

const props = defineProps<{ worldId: string }>();
const chat = useChatStore();
const message = useMessage();
const invites = ref<any[]>([]);
const loading = ref(false);

const loadInvites = async () => {
  if (!props.worldId) return;
  loading.value = true;
  try {
    const resp = await chat.loadWorldSections(props.worldId, ['invites']);
    const list = Array.isArray(resp.invites) ? resp.invites : [];
    invites.value = list.length ? [list[0]] : [];
    latestInvite.value = invites.value[0] || null;
  } catch (e) {
    message.error('加载邀请失败');
  } finally {
    loading.value = false;
  }
};

const showInviteModal = ref(false);
const latestInvite = ref<any>(null);
const showCreateModal = ref(false);
const inviteForm = ref({ ttlMinutes: 0, maxUse: 0, memo: '' });

const resetForm = () => {
  inviteForm.value = { ttlMinutes: 0, maxUse: 0, memo: '' };
};

const saveInvite = async () => {
  if (!props.worldId) return;
  try {
    const ttl = Math.max(0, Number(inviteForm.value.ttlMinutes) || 0);
    const maxUse = Math.max(0, Number(inviteForm.value.maxUse) || 0);
    const payload: any = {
      ttlMinutes: ttl,
      maxUse: maxUse,
      memo: inviteForm.value.memo?.trim() || undefined,
    };
    const resp = await chat.createWorldInvite(props.worldId, payload);
    latestInvite.value = resp.invite || null;
    showInviteModal.value = true;
    showCreateModal.value = false;
    message.success('已创建邀请');
    await loadInvites();
  } catch (e: any) {
    message.error(e?.response?.data?.message || '创建邀请失败');
  }
};

const copySlug = async (slug: string) => {
  try {
    if (navigator?.clipboard?.writeText) {
      await navigator.clipboard.writeText(slug);
    } else {
      const textarea = document.createElement('textarea');
      textarea.value = slug;
      textarea.style.position = 'fixed';
      textarea.style.opacity = '0';
      document.body.appendChild(textarea);
      textarea.focus();
      textarea.select();
      document.execCommand('copy');
      document.body.removeChild(textarea);
    }
    message.success('已复制邀请码');
  } catch (e) {
    message.error('复制失败，请手动选择后复制');
  }
};

const buildInviteLink = (slug: string) => {
  const origin = typeof window !== 'undefined' ? window.location.origin : '';
  return `${origin}/#/invite/${slug}`;
};

onMounted(loadInvites);
</script>

<template>
  <div class="space-y-2">
    <div class="flex justify-between items-center">
      <h3 class="font-bold">邀请列表</h3>
      <n-button size="small" type="primary" @click="() => { resetForm(); showCreateModal = true; }">创建邀请</n-button>
    </div>
    <n-list bordered>
      <n-spin :show="loading">
        <n-empty v-if="!invites.length" description="暂无邀请" />
        <n-list-item v-for="item in invites" :key="item.id">
          <div class="flex items-center justify-between">
            <div>
              <div class="text-sm font-medium break-all">{{ item.slug }}</div>
              <div class="text-xs text-gray-500">使用 {{ item.usedCount }} / {{ item.maxUse || '∞' }}</div>
            </div>
            <n-space size="small">
              <n-button size="tiny" @click="() => { latestInvite.value = item; showInviteModal.value = true; }">查看</n-button>
              <n-button size="tiny" @click="copySlug(item.slug)">复制邀请码</n-button>
            </n-space>
          </div>
        </n-list-item>
      </n-spin>
    </n-list>
    <n-modal v-model:show="showInviteModal" preset="dialog" title="邀请详情" style="max-width:420px">
      <n-space vertical v-if="latestInvite">
        <n-input readonly :value="buildInviteLink(latestInvite.slug)" />
        <n-statistic label="可用次数" :value="latestInvite.maxUse || '无限'" />
        <n-statistic label="已使用" :value="latestInvite.usedCount" />
      </n-space>
      <template #action>
        <n-space>
          <n-button quaternary @click="showInviteModal = false">关闭</n-button>
          <n-button type="primary" @click="copySlug(buildInviteLink(latestInvite?.slug ?? ''))">复制链接</n-button>
        </n-space>
      </template>
    </n-modal>
    <n-modal v-model:show="showCreateModal" preset="dialog" title="创建邀请" style="max-width:520px">
      <n-form label-placement="left" label-width="96">
        <n-form-item label="有效期(分钟)">
          <n-space>
            <n-input-number v-model:value="inviteForm.ttlMinutes" :min="0" :step="30" placeholder="0 表示永久" />
            <n-radio-group v-model:value="inviteForm.ttlMinutes" size="small">
              <n-space>
                <n-radio-button :value="0">永久</n-radio-button>
                <n-radio-button :value="30">30 分钟</n-radio-button>
                <n-radio-button :value="60">1 小时</n-radio-button>
                <n-radio-button :value="60 * 24">1 天</n-radio-button>
              </n-space>
            </n-radio-group>
          </n-space>
        </n-form-item>
        <n-form-item label="可用次数">
          <n-space>
            <n-input-number v-model:value="inviteForm.maxUse" :min="0" :step="1" placeholder="0 表示无限" />
            <n-radio-group v-model:value="inviteForm.maxUse" size="small">
              <n-space>
                <n-radio-button :value="0">无限</n-radio-button>
                <n-radio-button :value="1">1 次</n-radio-button>
                <n-radio-button :value="5">5 次</n-radio-button>
                <n-radio-button :value="10">10 次</n-radio-button>
              </n-space>
            </n-radio-group>
          </n-space>
        </n-form-item>
        <n-form-item label="备注">
          <n-input v-model:value="inviteForm.memo" type="textarea" autosize placeholder="可选，方便区分用途" />
        </n-form-item>
      </n-form>
      <template #action>
        <n-space>
          <n-button quaternary @click="showCreateModal = false">取消</n-button>
          <n-button type="primary" @click="saveInvite">保存</n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>
