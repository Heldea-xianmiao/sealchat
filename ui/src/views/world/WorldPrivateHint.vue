<script setup lang="ts">
import { computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useUserStore } from '@/stores/user';

const route = useRoute();
const router = useRouter();
const user = useUserStore();

const worldId = computed(() => (typeof route.params.worldId === 'string' ? route.params.worldId.trim() : ''));
const redirectPath = computed(() => (typeof route.query.redirect === 'string' ? route.query.redirect : ''));
const loginLabel = computed(() => (user.token ? '切换账号' : '登录'));

const goLogin = () => {
  const redirect = redirectPath.value || (worldId.value ? `/${worldId.value}` : '/');
  router.push({ name: 'user-signin', query: { redirect } });
};

const goLobby = () => {
  router.push({ name: 'world-lobby' });
};

const goHome = () => {
  router.push({ name: 'home' });
};
</script>

<template>
  <div class="world-private-hint">
    <n-card class="world-private-hint__card" title="需要邀请才能访问">
      <div class="world-private-hint__content">
        <p>该世界为私有或未开放公开访问。</p>
        <p>如果你拥有成员邀请链接，请使用邀请链接加入后再打开。</p>
      </div>
      <div class="world-private-hint__actions">
        <n-button size="small" @click="goLobby">世界大厅</n-button>
        <n-button size="small" @click="goHome">返回首页</n-button>
        <n-button size="small" type="primary" @click="goLogin">{{ loginLabel }}</n-button>
      </div>
    </n-card>
  </div>
</template>

<style scoped>
.world-private-hint {
  min-height: 100vh;
  padding: 2rem;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--sc-bg-body, #f8fafc);
}

.world-private-hint__card {
  width: min(480px, 100%);
}

.world-private-hint__content {
  display: grid;
  gap: 0.5rem;
  color: var(--sc-text-secondary, #4b5563);
  line-height: 1.6;
}

.world-private-hint__actions {
  margin-top: 1.25rem;
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}
</style>
