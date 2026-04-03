<template>
  <div class="main-wrapper">
    <!-- 左侧分类栏 -->
    <div class="category-sidebar">
      <div v-if="categories.length === 0" class="category-loading">
        加载中...
      </div>
      <div
        v-for="cat in categories"
        :key="cat.id"
        :class="['category-item', { active: selectedCategory === cat.id }]"
        @click="selectCategory(cat.id)"
      >
        {{ cat.name }}
      </div>
    </div>

    <!-- 右侧政策列表 -->
    <div class="content-area">
      <!-- Loading状态 -->
      <div v-if="loading" class="loading-container">
        <div class="loading-spinner"></div>
        <div class="loading-text">加载中...</div>
      </div>

      <!-- 空数据提示 -->
      <div v-else-if="policyList.length === 0" class="empty-container">
        <svg
          class="empty-icon"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="1.5"
        >
          <path
            d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
          />
        </svg>
        <div class="empty-text">暂无相关政策数据</div>
      </div>

      <!-- 列表 -->
      <div v-else>
        <div
          v-for="item in policyList"
          :key="item.id"
          class="list-card"
          @click="goToDetail(item)"
        >
          <div class="code">
            {{ item.standard_no }}
            <span :class="['status-tag', getStatusClass(item.status)]">{{
              item.status
            }}</span>
          </div>
          <div class="title">{{ item.standard_name }}</div>
          <div class="org">{{ getPublisher(item) }}</div>
          <div class="date">{{ item.publish_date || "未知日期" }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from "vue";
import { useRouter, useRoute } from "vue-router";
import { mobileApi } from "@/api/policy";

const props = defineProps<{
  policyType?: string;
  keyword?: string;
  highlightCategoryId?: number; // 需要高亮的分类ID
}>();

const router = useRouter();
const route = useRoute();

const categories = ref<any[]>([]);
const selectedCategory = ref(0);
const policyList = ref<any[]>([]);
const loading = ref(false);

const currentType = () => {
  if (props.policyType) return props.policyType;
  return (route.params.type as string) || "national";
};

const fetchCategories = async () => {
  try {
    const res = (await mobileApi.getCategories()) as any;
    const data = res.data?.data || res.data || {};
    const list = data.list || [];
    categories.value = [{ id: 0, name: "全部" }, ...list];
  } catch (e) {
    console.error("获取分类失败:", e);
    categories.value = [{ id: 0, name: "全部" }];
  }
};

const fetchPolicyList = async () => {
  loading.value = true;
  try {
    const params: any = { page: 1, pageSize: 50 };
    if (selectedCategory.value > 0) {
      params.categoryId = selectedCategory.value;
    }
    if (props.keyword) {
      params.keyword = props.keyword;
    }

    const res = (await mobileApi.getListByType(currentType(), params)) as any;
    const data = res.data?.data || res.data || {};
    policyList.value = data.list || [];
  } catch (e) {
    console.error("获取政策列表失败:", e);
    policyList.value = [];
  } finally {
    loading.value = false;
  }
};

const selectCategory = (catId: number) => {
  selectedCategory.value = catId;
  fetchPolicyList();
};

const getStatusClass = (status: string) => {
  if (status === "现行") return "status-current";
  if (status === "即将实施") return "status-future";
  return "status-obsolete";
};

const getPublisher = (item: any) => {
  return (
    item.publisher ||
    item.approve_dept ||
    item.technical_dept ||
    item.department ||
    "未知机构"
  );
};

const goToDetail = (item: any) => {
  router.push(`/m/detail/${currentType()}/${item.id}`);
};

watch(
  () => props.policyType,
  () => {
    selectedCategory.value = 0;
    fetchPolicyList();
  },
);

watch(
  () => props.keyword,
  (newVal) => {
    if (newVal) {
      selectedCategory.value = 0;
      fetchPolicyList();
    }
  },
);

// 监听高亮分类ID变化
watch(
  () => props.highlightCategoryId,
  (newVal) => {
    if (newVal && newVal > 0) {
      selectedCategory.value = newVal;
      fetchPolicyList();
    }
  },
  { immediate: true }
);

onMounted(() => {
  fetchCategories();
  fetchPolicyList();
});
</script>

<style scoped>
.main-wrapper {
  display: flex;
  height: calc(100vh - 100px);
}

.category-sidebar {
  width: 6rem;
  background: #fff;
  border-right: 1px solid #f2f3f5;
  overflow-y: auto;
  flex-shrink: 0;
}

.category-loading {
  padding: 20px 12px;
  text-align: center;
  color: #999;
  font-size: 12px;
}

.category-item {
  padding: 14px 12px;
  font-size: 14px;
  color: #666;

  cursor: pointer;
  transition: all 0.2s;
}

.category-item:active {
  background: #f5f5f5;
}

.category-item.active {
  background: #fef7f4;
  color: #d92a2b;

  font-weight: 500;
}

.content-area {
  flex: 1;
  padding: 12px 0;
  overflow-y: auto;
  background: #fafafa;
  scrollbar-width: thin;
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 0;
}

.loading-spinner {
  width: 28px;
  height: 28px;
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
  margin-top: 10px;
  color: #999;
  font-size: 13px;
}

.empty-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
}

.empty-icon {
  width: 48px;
  height: 48px;
  color: #ddd;
}

.empty-text {
  margin-top: 12px;
  color: #999;
  font-size: 14px;
}

.list-card {
  padding: 12px 16px;
  margin-bottom: 8px;
  background: #fff;
  display: flex;
  flex-direction: column;
  gap: 4px;
  cursor: pointer;
  transition: background 0.2s;
}

.list-card:active {
  background: #f5f5f5;
}

.list-card .title {
  font-size: 15px;
  font-weight: 500;
  color: #333;
  line-height: 1.4;
}

.list-card .code {
  font-size: 14px;
  color: #333;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.list-card .org {
  font-size: 12px;
  color: #d92a2b;
}

.list-card .date {
  font-size: 12px;
  color: #999;
}
</style>

