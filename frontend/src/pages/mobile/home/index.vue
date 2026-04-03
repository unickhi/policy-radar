<template>
  <div>
    <div
      class="top"
      style="
        position: sticky;
        top: 0;
        background: linear-gradient(
          90deg,
          rgba(235, 67, 28, 1) 0%,
          rgba(210, 13, 5, 1) 100%
        );
      "
    >
      <!-- 顶部导航 -->
      <div class="top-nav">
        <div class="nav-left">
          <div class="nav-title">政策雷达-电力</div>
        </div>
        <div class="nav-right">
          <input
            class="search-input"
            placeholder="搜索标准编号/名称"
            v-model="searchKeyword"
            @keyup.enter="handleSearch"
          />
          <svg
            class="nav-icon"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            @click="handleSearch"
          >
            <circle cx="11" cy="11" r="8"></circle>
            <path d="m21 21-4.35-4.35" />
          </svg>
        </div>
      </div>

      <!-- 标签切换 -->

      <div class="tab-switcher">
        <div
          v-for="(tab, index) in tabs"
          :key="index"
          :class="['tab-item', { active: activeTab === index }]"
          @click="switchTab(index)"
        >
          {{ tab.name }}
        </div>
      </div>
    </div>

    <!-- 首页内容 -->
    <div v-if="activeTab === 0">
      <!-- Loading状态 -->
      <div v-if="loading" class="loading-container">
        <div class="loading-spinner"></div>
        <div class="loading-text">加载中...</div>
      </div>

      <template v-else>
        <!-- Banner -->
        <div
          class="banner-section"
          v-if="banners.length > 0"
          @click="handleBannerClick(banners[0])"
        >
          <div class="banner-title">{{ banners[0].title || "政策推荐" }}</div>
          <div class="banner-desc">
            {{ (banners[0].content || "").slice(0, 60) }}...
          </div>
        </div>

        <!-- 国标 -->
        <div class="list-section">
          <div class="section-header">
            <span>国标 - 电力 (GB)</span>
            <span class="section-more" @click="goToCategory('national')"
              >更多 &gt;</span
            >
          </div>
          <div v-if="nationalList.length === 0" class="empty-tip">
            暂无国标数据
          </div>
          <div
            v-for="item in nationalList.slice(0, 3)"
            :key="item.id"
            class="standard-card"
            @click="goToDetail('national', item.id)"
          >
            <div class="standard-header">
              <span class="standard-code"
                >{{ item.standard_no + " · " }}
                <span class="time-text">{{ getTimeDisplay(item) }}</span>
              </span>

              <span :class="['status-tag', getStatusClass(item.status)]">{{
                item.status
              }}</span>
            </div>
            <div class="standard-desc">{{ item.standard_name }}</div>
            <div class="standard-meta">
              {{ item.publisher || "未知机构" }} ·
              {{ item.publish_date || "未知日期" }}
            </div>
          </div>
        </div>

        <!-- 地标 -->
        <div class="list-section">
          <div class="section-header">
            <span>地标 - 电力 (DB)</span>
            <span class="section-more" @click="goToCategory('local')"
              >更多 &gt;</span
            >
          </div>
          <div v-if="localList.length === 0" class="empty-tip">
            暂无地标数据
          </div>
          <div
            v-for="item in localList.slice(0, 3)"
            :key="item.id"
            class="standard-card"
            @click="goToDetail('local', item.id)"
          >
            <div class="standard-header">
              <span class="standard-code">{{ item.standard_no }}</span>
              <span v-if="item.status === '即将实施'" class="time-badge">{{
                getTimeDisplay(item)
              }}</span>
              <span :class="['status-tag', getStatusClass(item.status)]">{{
                item.status
              }}</span>
            </div>
            <div class="standard-desc">{{ item.standard_name }}</div>
            <div class="standard-meta">
              {{ item.publisher || item.department || "未知机构" }} ·
              {{ item.publish_date || "未知日期" }}
            </div>
          </div>
        </div>

        <!-- 行标 -->
        <div class="list-section">
          <div class="section-header">
            <span>行标 - 电力 (DL)</span>
            <span class="section-more" @click="goToCategory('industry')"
              >更多 &gt;</span
            >
          </div>
          <div v-if="industryList.length === 0" class="empty-tip">
            暂无行标数据
          </div>
          <div
            v-for="item in industryList.slice(0, 3)"
            :key="item.id"
            class="standard-card"
            @click="goToDetail('industry', item.id)"
          >
            <div class="standard-header">
              <span class="standard-code">{{ item.standard_no }}</span>
              <span v-if="item.status === '即将实施'" class="time-badge">{{
                getTimeDisplay(item)
              }}</span>
              <span :class="['status-tag', getStatusClass(item.status)]">{{
                item.status
              }}</span>
            </div>
            <div class="standard-desc">{{ item.standard_name }}</div>
            <div class="standard-meta">
              {{ item.approve_dept || "未知机构" }} ·
              {{ item.publish_date || "未知日期" }}
            </div>
          </div>
        </div>
      </template>
    </div>

    <!-- 分类页内容（其他Tab） -->
    <div v-else>
      <CategoryView
        :policy-type="tabs[activeTab].type"
        :keyword="searchKeywordForCategory"
        :highlight-category-id="searchCategoryId"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import { mobileApi } from "@/api/policy";
import CategoryView from "@/pages/mobile/category/index.vue";

const router = useRouter();

const tabs = [
  { name: "近期", type: "" },
  { name: "国标", type: "national" },
  { name: "地标", type: "local" },
  { name: "行标", type: "industry" },
];

const activeTab = ref(0);
const searchKeyword = ref("");
const searchKeywordForCategory = ref(""); // 传递给分类页的搜索关键词
const searchCategoryId = ref(0); // 搜索结果高亮的分类ID
const loading = ref(false);
const banners = ref<any[]>([]);
const nationalList = ref<any[]>([]);
const localList = ref<any[]>([]);
const industryList = ref<any[]>([]);

const fetchData = async () => {
  loading.value = true;
  try {
    const res = (await mobileApi.getHomeData()) as any;
    // axios响应在res.data中，API返回结构是{code, message, data}
    const data = res.data?.data || res.data || {};
    banners.value = data.banners || [];
    nationalList.value = data.national?.list || [];
    localList.value = data.local?.list || [];
    industryList.value = data.industry?.list || [];
  } catch (e) {
    console.error("获取首页数据失败:", e);
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  fetchData();
});

const switchTab = (index: number) => {
  activeTab.value = index;
  // 切换Tab时清除搜索关键词和分类
  if (index === 0) {
    searchKeywordForCategory.value = "";
    searchCategoryId.value = 0;
  }
};

const getStatusClass = (status: string) => {
  if (status === "现行") return "status-current";
  if (status === "即将实施") return "status-future";
  return "status-obsolete";
};

// 计算时间显示：即将实施显示倒计时，现行显示已实施时长
const getTimeDisplay = (item: any) => {
  const implementDate = item.implement_date;
  if (!implementDate) return "";

  // 清理日期格式
  const cleanDateStr = implementDate
    .replace(/0+$/, "")
    .replace(/(\d{4}-\d{1,2}-\d{1,2}).*/, "$1");

  try {
    const implDate = new Date(cleanDateStr);
    const now = new Date();

    if (isNaN(implDate.getTime())) return "";

    if (item.status === "即将实施") {
      // 计算距离实施还有多少时间
      const diff = implDate.getTime() - now.getTime();
      if (diff <= 0) return "即将生效";

      const days = Math.floor(diff / (1000 * 60 * 60 * 24));
      const months = Math.floor(days / 30);
      const remainingDays = days % 30;

      if (months > 0) {
        return `${months}个月${remainingDays.toString().padStart(2, "0")}天`;
      }
      return `${days.toString().padStart(2, "0")}天`;
    } else if (item.status === "现行") {
      // 计算已实施时长
      let years = now.getFullYear() - implDate.getFullYear();
      let months = now.getMonth() - implDate.getMonth();
      let days = now.getDate() - implDate.getDate();

      if (days < 0) {
        months--;
        days += new Date(now.getFullYear(), now.getMonth(), 0).getDate();
      }
      if (months < 0) {
        years--;
        months += 12;
      }

      if (years > 0) {
        return `已实施 ${years}年${months}个月${days}天`;
      }
      return `已实施 ${months}个月${days}天`;
    }
    return "";
  } catch (e) {
    return "";
  }
};

const goToDetail = (type: string, id: number, content?: string) => {
  if (content) {
    // 带推荐内容跳转
    router.push({
      path: `/m/detail/${type}/${id}`,
      query: { recommend: encodeURIComponent(content) },
    });
  } else {
    router.push(`/m/detail/${type}/${id}`);
  }
};

const goToCategory = (type: string) => {
  // 顺序：推荐(0), 国标(1), 地标(2), 行标(3)
  if (type === "national") activeTab.value = 1;
  else if (type === "local") activeTab.value = 2;
  else if (type === "industry") activeTab.value = 3;
};

// 根据搜索结果类型获取Tab索引
const getTabByType = (type: string): number => {
  switch (type) {
    case "national":
      return 1;
    case "local":
      return 2;
    case "industry":
      return 3;
    default:
      return 1; // 默认国标
  }
};

const handleSearch = async () => {
  if (!searchKeyword.value.trim()) {
    return;
  }

  loading.value = true;

  try {
    const res = (await mobileApi.search(searchKeyword.value.trim())) as any;
    const results = res.data?.data || res.data || [];

    if (results.length > 0) {
      // 有搜索结果，根据第一条结果的类型决定跳转到哪个Tab
      const firstResult = results[0];
      activeTab.value = getTabByType(firstResult.type);
      searchKeywordForCategory.value = searchKeyword.value.trim();
      // 如果有分类，高亮对应的分类
      searchCategoryId.value = firstResult.category_id || 0;
    } else {
      // 没有搜索结果，跳转到国标Tab
      activeTab.value = 1;
      searchKeywordForCategory.value = searchKeyword.value.trim();
      searchCategoryId.value = 0;
    }
  } catch (e) {
    console.error("搜索失败:", e);
    // 搜索失败，跳转到国标Tab
    activeTab.value = 1;
    searchKeywordForCategory.value = searchKeyword.value.trim();
    searchCategoryId.value = 0;
  } finally {
    loading.value = false;
  }
};

const handleBannerClick = (banner: any) => {
  if (banner.policy_id && banner.policy_type) {
    const type =
      banner.policy_type === "国标"
        ? "national"
        : banner.policy_type === "行标"
          ? "industry"
          : "local";
    goToDetail(type, banner.policy_id, banner.content);
  }
};
</script>

<style scoped>
.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 0;
}

.loading-spinner {
  width: 32px;
  height: 32px;
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

.empty-tip {
  text-align: center;
  padding: 30px;
  color: #999;
  font-size: 14px;
}
</style>

