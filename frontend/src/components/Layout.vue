<template>
  <div class="h-screen flex bg-gray-50">
    <!-- 侧边栏 -->
    <aside class="w-56 sidebar flex-shrink-0">
      <div class="h-full flex flex-col">
        <!-- Logo区域 -->
        <div class="h-14 flex items-center justify-center bg-white/10">
          <h1 class="text-white text-lg font-bold">政策雷达-电力</h1>
        </div>

        <!-- 菜单 -->
        <el-menu
          :default-active="currentRoute"
          class="sidebar-menu flex-1 border-none"
          router
        >
          <el-menu-item index="/admin/hot-update">
            <el-icon><Refresh /></el-icon>
            <span>政策热更新</span>
          </el-menu-item>

          <!-- 数据管理 -->
          <el-sub-menu index="policy">
            <template #title>
              <el-icon><Document /></el-icon>
              <span>Data - 电力标准</span>
            </template>
            <el-menu-item index="/admin/national">国标</el-menu-item>
            <el-menu-item index="/admin/industry">行标</el-menu-item>
            <el-menu-item index="/admin/local">地标</el-menu-item>
          </el-sub-menu>

          <el-menu-item index="/admin/category">
            <el-icon><Folder /></el-icon>
            <span>政策分类</span>
          </el-menu-item>

          <el-menu-item index="/admin/recommend">
            <el-icon><Star /></el-icon>
            <span>推荐政策</span>
          </el-menu-item>

          <el-menu-item index="/admin/dashboard">
            <el-icon><DataLine /></el-icon>
            <span>数据看板</span>
          </el-menu-item>
        </el-menu>

        <div class="p-3 text-white/50 text-xs text-center">v1.2.0</div>
      </div>
    </aside>

    <!-- 主内容区 -->
    <main class="flex-1 flex flex-col overflow-hidden">
      <!-- 顶部栏 -->
      <header
        class="h-14 bg-white border-b flex items-center justify-between px-5"
      >
        <div class="flex items-center text-gray-500">
          <el-icon class="mr-1"><HomeFilled /></el-icon>
          <span class="text-sm">{{ currentTitle }}</span>
        </div>
      </header>

      <!-- 内容区 -->
      <div class="flex-1 p-5 overflow-auto bg-gray-50">
        <router-view />
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { useRoute } from "vue-router";

const route = useRoute();
const currentRoute = computed(() => route.path);
const currentTitle = computed(() => route.meta.title || "政策管理");
</script>

<style scoped>
.sidebar {
  background-color: #db373b;
}

.sidebar-menu {
  background-color: transparent;
}

.sidebar-menu :deep(.el-menu-item),
.sidebar-menu :deep(.el-sub-menu__title) {
  color: #fff;
  height: 48px;
  line-height: 48px;
}

.sidebar-menu :deep(.el-menu-item:hover),
.sidebar-menu :deep(.el-sub-menu__title:hover) {
  background-color: rgba(255, 255, 255, 0.1);
}

.sidebar-menu :deep(.el-menu-item.is-active) {
  background-color: rgba(255, 255, 255, 0.15);
  color: #fff;
}

.sidebar-menu :deep(.el-sub-menu.is-active > .el-sub-menu__title) {
  color: #fff;
}

/* 减少圆角 */
:deep(.el-button) {
  border-radius: 2px !important;
}

:deep(.el-input__wrapper),
:deep(.el-select__wrapper),
:deep(.el-textarea__inner) {
  border-radius: 2px !important;
}

:deep(.el-drawer) {
  border-radius: 0 !important;
}

:deep(.el-table) {
  border-radius: 0 !important;
}

:deep(.el-tag) {
  border-radius: 2px !important;
}
</style>

