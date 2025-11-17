<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router';
import { onMounted, ref, computed } from 'vue';
import { useChatStore } from '@/stores/chat';
import { useUserStore } from '@/stores/user';
import { useMessage } from 'naive-ui';

const route = useRoute();
const router = useRouter();
const chat = useChatStore();
const user = useUserStore();
const message = useMessage();

const slug = computed(() => (route.params.slug as string) || (route.query.invite as string) || '');
const status = ref<'pending' | 'processing' | 'success' | 'error'>('pending');
const errorMessage = ref('');

const processInvite = async () => {
  if (!slug.value) {
    status.value = 'error';
    errorMessage.value = '缺少邀请码';
    return;
  }
  if (!user.token) {
    router.replace({ name: 'user-signin', query: { redirect: route.fullPath } });
    return;
  }
  status.value = 'processing';
  try {
    await chat.ensureConnectionReady();
    const resp = await chat.consumeWorldInvite(slug.value);
    const worldId = resp.world?.id;
    if (worldId) {
      await chat.switchWorld(worldId, { force: true });
      status.value = 'success';
      message.success('已加入世界');
      try {
        await router.replace({ name: 'home' });
      } catch (err) {
        console.warn('router replace failed', err);
      }
      if (router.currentRoute.value.name !== 'home') {
        window.location.hash = '#/';
      }
    } else {
      status.value = 'error';
      errorMessage.value = '加入失败，世界信息缺失';
    }
  } catch (e: any) {
    status.value = 'error';
    errorMessage.value = e?.response?.data?.message || '加入失败';
  }
};

onMounted(processInvite);
</script>

<template>
  <div class="w-full h-full flex items-center justify-center p-6">
    <div class="text-center space-y-4">
      <n-spin :show="status === 'processing'">
        <template v-if="status === 'pending' || status === 'processing'">
          <p>正在验证邀请链接...</p>
        </template>
        <template v-else-if="status === 'success'">
          <p>邀请成功，正在跳转...</p>
        </template>
        <template v-else>
          <n-alert type="error" title="加入失败">{{ errorMessage }}</n-alert>
        </template>
      </n-spin>
    </div>
  </div>
</template>
