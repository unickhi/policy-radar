<template>
  <div class="webview-page">
    <!-- 顶部导航 -->
    <div class="webview-header">
      <div class="header-left" @click="goBack">
        <svg class="back-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M19 12H5M12 19l-7-7 7-7" />
        </svg>
      </div>
      <div class="header-title">{{ title || '政策溯源' }}</div>
      <div class="header-right" @click="openInBrowser">
        <svg class="browser-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6" />
          <polyline points="15 3 21 3 21 9" />
          <line x1="10" y1="14" x2="21" y2="3" />
        </svg>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-container">
      <div class="loading-spinner"></div>
      <div class="loading-text">加载中...</div>
    </div>

    <!-- iframe -->
    <iframe
      ref="iframeRef"
      :src="url"
      class="webview-iframe"
      @load="onLoad"
      sandbox="allow-scripts allow-same-origin allow-forms allow-popups"
    ></iframe>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";

const route = useRoute();
const router = useRouter();

const url = ref("");
const title = ref("");
const loading = ref(true);
const iframeRef = ref<HTMLIFrameElement | null>(null);

onMounted(() => {
  url.value = (route.query.url as string) || "";
  title.value = (route.query.title as string) || "政策溯源";
});

const goBack = () => {
  router.back();
};

const openInBrowser = () => {
  if (url.value) {
    window.open(url.value, "_blank");
  }
};

const onLoad = () => {
  loading.value = false;
};
</script>

<style scoped>
.webview-page {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: #fff;
}

.webview-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 48px;
  padding: 0 16px;
  background: linear-gradient(90deg, rgba(235, 67, 28, 1) 0%, rgba(210, 13, 5, 1) 100%);
  color: #fff;
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-left,
.header-right {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}

.back-icon,
.browser-icon {
  width: 22px;
  height: 22px;
}

.header-title {
  flex: 1;
  text-align: center;
  font-size: 17px;
  font-weight: 600;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.loading-container {
  position: absolute;
  top: 48px;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: #fff;
  z-index: 10;
}

.loading-spinner {
  width: 36px;
  height: 36px;
  border: 3px solid #f3f3f3;
  border-top: 3px solid #d92a2b;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

.loading-text {
  margin-top: 12px;
  color: #999;
  font-size: 14px;
}

.webview-iframe {
  flex: 1;
  width: 100%;
  border: none;
  background: #fff;
}
</style>