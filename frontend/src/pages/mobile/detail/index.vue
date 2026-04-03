<template>
  <div>
    <!-- 顶部导航 -->
    <div class="top-nav">
      <div class="nav-left">
        <svg
          class="nav-back"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          @click="goBack"
        >
          <path d="M19 12H5M12 19l-7-7 7-7" />
        </svg>
        <div class="nav-title">政策详情</div>
      </div>
    </div>

    <!-- Loading状态 -->
    <div v-if="loading" class="loading-container">
      <div class="loading-spinner"></div>
      <div class="loading-text">加载中...</div>
    </div>

    <!-- 详情内容 -->
    <div v-else class="detail-page">
      <!-- 标题区 -->
      <div class="detail-title">{{ detail.standard_name || "未知标准" }}</div>
      <div class="detail-subtitle">{{ detail.standard_no || "未知编号" }}</div>
      <div v-if="detail.english_name" class="detail-en-title">
        {{ detail.english_name }}
      </div>

      <!-- 状态标签组 -->
      <div class="detail-tags">
        <span
          :class="[
            'detail-tag',
            detail.status === '现行'
              ? 'current'
              : detail.status === '即将实施'
                ? 'future'
                : 'normal',
          ]"
        >
          {{ detail.status || "未知状态" }}
        </span>
        <span
          :class="[
            'detail-tag',
            detail.status === '即将实施' ? 'time-warning' : 'time-normal',
          ]"
        >
          {{ calculateDuration(detail.implement_date, detail.status) }}
        </span>
        <span class="detail-tag normal" v-if="detail.implement_date"
          >{{ cleanDate(detail.implement_date) }} 实施</span
        >
        <span class="detail-tag normal">{{
          detail.policy_type || "标准"
        }}</span>
        <span
          v-if="detail.source_link"
          class="policy-source-tag"
          @click="goToWebview"
        >
          政策溯源 ↗
        </span>
      </div>

      <!-- 信息网格 -->
      <div class="info-grid">
        <div class="info-item" v-if="detail.category_name">
          <div class="info-label">政策分类</div>
          <div class="info-value">{{ detail.category_name }}</div>
        </div>
        <div class="info-item" v-if="detail.nature">
          <div class="info-label">标准性质</div>
          <div class="info-value">{{ detail.nature }}</div>
        </div>
        <div class="info-item" v-if="detail.ccs_code">
          <div class="info-label">CCS分类号</div>
          <div class="info-value">{{ detail.ccs_code }}</div>
        </div>
        <div class="info-item" v-if="detail.ics_code">
          <div class="info-label">ICS分类号</div>
          <div class="info-value">{{ detail.ics_code }}</div>
        </div>
        <div class="info-item" v-if="getTechnicalDept()">
          <div class="info-label">归口单位</div>
          <div class="info-value">{{ getTechnicalDept() }}</div>
        </div>
        <div class="info-item" v-if="getPublisher()">
          <div class="info-label">发布机构</div>
          <div class="info-value">{{ getPublisher() }}</div>
        </div>
        <div class="info-item" v-if="detail.is_adopted">
          <div class="info-label">是否采标</div>
          <div class="info-value">{{ detail.is_adopted }}</div>
        </div>
        <div class="info-item" v-if="detail.category">
          <div class="info-label">类别</div>
          <div class="info-value">{{ detail.category }}</div>
        </div>
        <div class="info-item" v-if="detail.replace_standard">
          <div class="info-label">替代标准</div>
          <div class="info-value">{{ detail.replace_standard }}</div>
        </div>
      </div>

      <!-- 推荐分析内容（从Banner跳转带过来的） -->
      <div v-if="recommendContent" class="detail-intro recommend-section">
        <div class="intro-title" style="color: #d92a2b">推荐分析</div>
        <div class="recommend-content">{{ recommendContent }}</div>
      </div>

      <!-- 详情简介 -->
      <div v-if="detail.description" class="detail-intro">
        <div class="intro-title">标准详情简介</div>
        <div class="intro-content">{{ detail.description }}</div>
      </div>

      <!-- 政策拆分 -->
      <div v-if="policySplits.length > 0">
        <div class="split-title">政策拆分</div>
        <div
          v-for="(split, index) in policySplits"
          :key="index"
          class="markdown-area"
        >
          <div v-if="split.title" class="split-subtitle">{{ split.title }}</div>
          <div class="split-content">{{ split.content }}</div>
        </div>
      </div>

      <!-- 相关推荐 -->
      <div v-if="relatedList.length > 0" class="related-section">
        <div class="related-title">相关政策</div>
        <div
          v-for="item in relatedList"
          :key="item.id"
          class="related-item"
          @click="goToDetail(item.type, item.id)"
        >
          <div class="related-code">
            <span>{{ item.standard_no }}</span>
            <span>{{ item.publish_date }}</span>
          </div>
          <div class="related-name">{{ item.standard_name }}</div>
        </div>
      </div>
    </div>

    <!-- 底部下载按钮 -->
    <div v-if="detail.download_url" class="detail-footer">
      <a :href="detail.download_url" target="_blank" class="download-btn">
        <svg
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
        >
          <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
          <polyline points="7 10 12 15 17 10" />
          <line x1="12" y1="15" x2="12" y2="3" />
        </svg>
        <span>下载</span>
      </a>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRouter, useRoute } from "vue-router";
import { mobileApi } from "@/api/policy";

const router = useRouter();
const route = useRoute();

const loading = ref(false);
const detail = ref<any>({});
const policySplits = ref<any[]>([]);
const relatedList = ref<any[]>([]);
const recommendContent = ref("");

const fetchDetail = async () => {
  const type = route.params.type as string;
  const id = Number(route.params.id);

  loading.value = true;
  try {
    const res = (await mobileApi.getDetail(type, id)) as any;
    const data = res.data?.data || res.data || {};
    detail.value = data;
    policySplits.value = data.policy_splits || [];
    relatedList.value = data.related || [];
  } catch (e) {
    console.error("获取详情失败:", e);
  } finally {
    loading.value = false;
  }
};

const cleanDate = (dateStr: string) => {
  if (!dateStr) return "";
  // 去除末尾多余的0
  return dateStr.replace(/0+$/, "").replace(/(\d{4}-\d{1,2}-\d{1,2}).*/, "$1");
};

const calculateDuration = (implementDate: string, status?: string) => {
  if (!implementDate) return "未知实施时间";

  const cleanDateStr = cleanDate(implementDate);

  try {
    const implDate = new Date(cleanDateStr);
    const now = new Date();

    if (isNaN(implDate.getTime())) return "未知实施时间";

    // 即将实施：显示倒计时
    if (status === "即将实施") {
      const diff = implDate.getTime() - now.getTime();
      if (diff <= 0) return "即将生效";

      const days = Math.floor(diff / (1000 * 60 * 60 * 24));
      const months = Math.floor(days / 30);
      const remainingDays = days % 30;

      if (months > 0) {
        return `距实施 ${months}个月${remainingDays.toString().padStart(2, "0")}天`;
      }
      return `距实施 ${days.toString().padStart(2, "0")}天`;
    }

    // 现行或其他状态：显示已实施时长
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
  } catch (e) {
    return "未知实施时间";
  }
};

const getTechnicalDept = () => {
  return (
    detail.value.technical_dept ||
    detail.value.technical_owner ||
    detail.value.department ||
    ""
  );
};

const getPublisher = () => {
  return detail.value.publisher || detail.value.approve_dept || "";
};

const goBack = () => {
  router.back();
};

const goToDetail = (type: string, id: number) => {
  router.push(`/m/detail/${type}/${id}`);
};

const goToWebview = () => {
  if (detail.value.source_link) {
    router.push({
      path: "/m/webview",
      query: {
        url: detail.value.source_link,
        title: detail.value.standard_name || "政策溯源",
      },
    });
  }
};

onMounted(() => {
  // 获取推荐内容参数
  const recommend = route.query.recommend as string;
  if (recommend) {
    recommendContent.value = decodeURIComponent(recommend);
  }

  fetchDetail();
});
</script>

<style scoped>
.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 100px 0;
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
  margin-top: 14px;
  color: #999;
  font-size: 14px;
}

.detail-page {
  padding: 16px;
  padding-bottom: 100px;
}

.detail-title {
  font-size: 20px;
  font-weight: 600;
  color: #333;
  margin-bottom: 8px;
  line-height: 1.4;
}

.detail-subtitle {
  font-size: 16px;
  color: #666;
  margin-bottom: 4px;
}

.detail-en-title {
  font-size: 14px;
  color: #999;
  margin-bottom: 16px;
  line-height: 1.5;
}

.detail-tags {
  display: flex;
  gap: 8px;
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid #f2f3f5;
  flex-wrap: wrap;
}

.detail-tag {
  padding: 4px 12px;
  border-radius: 4px;
  font-size: 13px;
}

.detail-tag.current {
  border: 2px solid #d92a2b;
  background: #fff;
  color: #d92a2b;
  font-weight: 500;
}

.detail-tag.future {
  border: 2px solid #f57c00;
  background: #fff;
  color: #f57c00;
  font-weight: 500;
}

.detail-tag.normal {
  background: #f2f3f5;
  color: #666;
}

.detail-tag.time-warning {
  background: #fff8e1;
  color: #e65100;
  font-weight: 500;
}

.detail-tag.time-normal {
  background: #e3f2fd;
  color: #1565c0;
}

.policy-source-tag {
  padding: 4px 12px;
  border-radius: 4px;
  font-size: 13px;
  background: #e8f3ff;
  color: #007aff;
  font-weight: 500;
  cursor: pointer;
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.info-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px 0;
  margin-bottom: 24px;
  padding-bottom: 24px;
  border-bottom: 1px solid #f2f3f5;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.info-label {
  font-size: 13px;
  color: #999;
}

.info-value {
  font-size: 15px;
  font-weight: 500;
  color: #333;
}

.recommend-content {
  font-size: 14px;
  color: #666;
  line-height: 1.8;
  white-space: pre-wrap;
}

.detail-intro {
  margin-bottom: 24px;
  padding-bottom: 24px;
  border-bottom: 1px solid #f2f3f5;
}

.intro-title {
  font-size: 16px;
  font-weight: 500;
  color: #333;
}

.intro-content {
  font-size: 14px;
  color: #666;
  line-height: 1.8;
  white-space: pre-wrap;
}

.split-title {
  font-size: 16px;
  font-weight: 500;
  color: #333;
  margin-bottom: 12px;
}

.markdown-area {
  background: #fafafa;
  border: 1px solid #eee;
  border-radius: 8px;
  padding: 16px;
  margin-bottom: 12px;
}

.split-subtitle {
  font-size: 15px;
  font-weight: 500;
  color: #333;
  margin-bottom: 8px;
}

.split-content {
  font-size: 14px;
  line-height: 1.8;
  color: #333;
  white-space: pre-wrap;
}

.related-section {
  margin-bottom: 80px;
}

.related-title {
  font-size: 16px;
  font-weight: 500;
  color: #333;
  margin-bottom: 12px;
}

.related-item {
  padding: 12px 0;
  border-bottom: 1px solid #f2f3f5;
  cursor: pointer;
}

.related-item:last-child {
  border-bottom: none;
}

.related-item:active {
  background: #f9f9f9;
}

.related-code {
  font-size: 13px;
  color: #999;
  display: flex;
  justify-content: space-between;
  margin-bottom: 4px;
}

.related-name {
  font-size: 15px;
  font-weight: 500;
  color: #333;
}

.detail-footer {
  position: fixed;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 100%;
  max-width: 450px;
  background: #fff;
  padding: 12px 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-top: 1px solid #f2f3f5;
  z-index: 100;
}

.download-btn {
  width: 100%;
  height: 44px;
  background: #000;
  color: #fff;
  border-radius: 22px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  text-decoration: none;
}

.download-btn svg {
  width: 20px;
  height: 20px;
}
</style>

