<template>
  <div class="flex flex-col h-full gap-4">
    <!-- 头部控制区 -->
    <div class="bg-white p-4 border-b">
      <div class="flex justify-between items-center">
        <div class="flex gap-4 items-center">
          <span class="text-gray-600">数据源:</span>
          <el-radio-group v-model="dataSource">
            <el-radio value="national">国标数据</el-radio>
            <el-radio value="industry">行标数据</el-radio>
            <el-radio value="local">地标数据</el-radio>
          </el-radio-group>
        </div>
        <div class="flex gap-2 items-center">
          <el-input
            v-model="queryKeyword"
            placeholder="查询关键词，如：电力"
            class="w-48"
          />
          <el-button type="primary" @click="handleExecute" :loading="executing">
            <el-icon class="mr-1"><Play /></el-icon>执行爬取
          </el-button>
        </div>
      </div>
    </div>

    <!-- 主内容区：左侧脚本+日志(70%)，右侧结果(30%) -->
    <div class="flex-1 flex gap-4 overflow-hidden">
      <!-- 左侧：脚本编辑器 + 执行日志 -->
      <div class="flex flex-col gap-4" style="width: 70%">
        <!-- 脚本编辑器 -->
        <div class="bg-white border flex flex-col" style="height: 70%">
          <div class="px-4 py-2 border-b bg-gray-50 text-gray-600 font-medium">
            Python 爬取脚本
          </div>
          <div class="flex-1 relative">
            <div ref="editorContainer" class="absolute inset-0"></div>
          </div>
        </div>

        <!-- 执行日志 -->
        <div class="bg-white border flex flex-col" style="height: 30%">
          <div class="terminal-header">
            <div class="terminal-buttons">
              <span class="terminal-btn terminal-btn-red"></span>
              <span class="terminal-btn terminal-btn-yellow"></span>
              <span class="terminal-btn terminal-btn-green"></span>
            </div>
            <span class="terminal-title">Terminal - 执行日志</span>
          </div>
          <div class="terminal-body flex-1" ref="terminalBody">
            <div v-if="logs.length === 0" class="terminal-placeholder">
              点击"执行爬取"开始运行脚本...
            </div>
            <template v-else>
              <div class="terminal-line terminal-prompt">
                <span class="terminal-dollar">$ </span>
                <span>python crawler.py --keyword="{{ queryKeyword }}"</span>
              </div>
              <div v-for="(log, idx) in logs" :key="idx" class="terminal-line">
                <span v-if="log.type === 'error'" class="terminal-error">{{
                  log.message
                }}</span>
                <span v-else class="terminal-output">{{ log.message }}</span>
              </div>
              <div v-if="executing" class="terminal-line terminal-cursor">
                <span class="cursor-blink">▋</span>
              </div>
            </template>
          </div>
        </div>
      </div>

      <!-- 右侧：爬取结果 -->
      <div class="bg-white border flex flex-col" style="width: 30%">
        <div
          class="px-4 py-2 border-b bg-gray-50 flex justify-between items-center"
        >
          <span class="text-gray-600 font-medium">爬取结果</span>
          <span v-if="resultData.length > 0" class="text-sm text-gray-400"
            >共 {{ resultData.length }} 条</span
          >
        </div>
        <div class="flex-1 overflow-auto">
          <div
            v-if="resultData.length === 0"
            class="h-full flex items-center justify-center text-gray-400"
          >
            <div class="text-center">
              <el-icon class="text-4xl mb-2"><Document /></el-icon>
              <div>点击"执行爬取"获取数据</div>
            </div>
          </div>
          <el-table
            v-else
            :data="resultData"
            stripe
            border
            size="small"
            height="100%"
          >
            <el-table-column prop="standard_no" label="标准号" width="120" />
            <el-table-column
              prop="standard_name"
              label="标准名称"
              min-width="150"
              show-overflow-tooltip
            />
            <el-table-column prop="status" label="状态" width="60">
              <template #default="{ row }">
                <el-tag
                  size="small"
                  :type="row.status === '现行' ? 'success' : 'warning'"
                  >{{ row.status }}</el-tag
                >
              </template>
            </el-table-column>
          </el-table>
        </div>
        <!-- 导入操作 -->
        <div
          v-if="resultData.length > 0"
          class="p-3 border-t flex justify-between items-center bg-gray-50"
        >
          <el-tag size="small" type="info">{{
            importTarget === "national" ? "国标" : importTarget === "industry" ? "行标" : "地标"
          }}</el-tag>
          <el-button
            type="success"
            size="small"
            @click="handleImport"
            :loading="importing"
          >
            <el-icon class="mr-1"><Upload /></el-icon>确认导入
          </el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, nextTick } from "vue";
import { ElMessage } from "element-plus";
import * as monaco from "monaco-editor";
import { crawlerApi } from "@/api/policy";

const editorContainer = ref<HTMLElement | null>(null);
let editor: monaco.editor.IStandaloneCodeEditor | null = null;

const dataSource = ref("national");
const queryKeyword = ref("电力");
const executing = ref(false);
const importing = ref(false);
const resultData = ref<any[]>([]);
const importTarget = ref("national");
const logs = ref<any[]>([]);
const terminalBody = ref<HTMLElement | null>(null);

// 国标爬取脚本（参考Datas原始脚本逻辑）
const nationalScript = `#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
国家标准爬虫 - 参考原始Datas脚本逻辑
接口1: std.samr.gov.cn (JSON API)
接口2: openstd.samr.gov.cn (HTML + 详情页)
合并: 按标准号取并集
"""

import requests
import json
import time
import re
import sys
import os
from bs4 import BeautifulSoup

# ========== 接口配置 ==========
API1_URL = "https://std.samr.gov.cn/gb/search/gbQueryPage"
DETAIL1_URL = "https://std.samr.gov.cn/gb/search/gbDetailed?id="
API2_URL = "https://openstd.samr.gov.cn/bzgk/gb/std_list"
DETAIL2_URL = "https://openstd.samr.gov.cn/bzgk/gb/newGbInfo?hcno="

HEADERS = {
    'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36',
    'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8',
    'Accept-Language': 'zh-CN,zh;q=0.9',
}

def log(msg):
    print(msg, file=sys.stderr)

def clean_date(date_str):
    """清理日期格式"""
    if not date_str:
        return ''
    date_str = re.sub(r'\\s+00:00:00.*', '', date_str)
    match = re.match(r'(\\d{4})-(\\d{1,2})-(\\d{1,2})', date_str)
    if match:
        year, month, day = match.groups()
        return f"{year}-{month.zfill(2)}-{day.zfill(2)}"
    return date_str

def crawl_api1(keyword, page_size=50):
    """接口1: JSON API获取列表（参考Datas/crawler.py原始脚本）"""
    log(f"[接口1] 开始爬取关键词: {keyword}")
    all_data = []
    page = 1

    while True:
        params = {
            'searchText': keyword,
            'ics': '',
            'state': '',
            'ISSUE_DATE': '-36',  # 近36个月
            'sortOrder': 'asc',
            'pageSize': page_size,
            'pageNumber': page,
            '_': int(time.time() * 1000)
        }
        try:
            resp = requests.get(API1_URL, params=params, headers=HEADERS, timeout=30)
            data = resp.json()
            rows = data.get('rows', [])
            if not rows:
                break
            all_data.extend(rows)
            log(f"[接口1] 第{page}页: {len(rows)}条")
            if len(rows) < page_size:
                break
            page += 1
            time.sleep(0.3)
        except Exception as e:
            log(f"[接口1] 错误: {e}")
            break

    log(f"[接口1] 总计: {len(all_data)} 条")
    return all_data

def get_detail_api2(hcno):
    """获取接口2详情页信息（参考Datas/crawler_guobiao_v2.py原始脚本）"""
    if not hcno:
        return {}
    detail_url = DETAIL2_URL + hcno
    try:
        resp = requests.get(detail_url, headers=HEADERS, timeout=30)
        resp.encoding = 'utf-8'
        soup = BeautifulSoup(resp.text, 'html.parser')
        text = soup.get_text(separator='\\n', strip=True)
        lines = text.split('\\n')

        detail = {}
        for i, line in enumerate(lines):
            line = line.strip()

            # 英文标准名称
            if '英文标准名称' in line:
                match = re.search(r'英文标准名称[：:]+(.+)', line)
                if match:
                    detail['english_name'] = match.group(1).strip()
                elif i + 1 < len(lines):
                    detail['english_name'] = lines[i+1].strip()

            # 中国标准分类号（CCS）
            elif '中国标准分类号（CCS）' in line or '中国标准分类号' in line:
                match = re.search(r'[分类号]+[：:]+([A-Z]\\d+)', line)
                if match:
                    detail['ccs_code'] = match.group(1)
                elif i + 1 < len(lines):
                    next_val = lines[i+1].strip()
                    if re.match(r'[A-Z]\\d+', next_val):
                        detail['ccs_code'] = next_val

            # 国际标准分类号（ICS）
            elif '国际标准分类号（ICS）' in line or '国际标准分类号' in line:
                match = re.search(r'[分类号]+[：:]+(\\d+\\.?\\d*)', line)
                if match:
                    detail['ics_code'] = match.group(1)
                elif i + 1 < len(lines):
                    next_val = lines[i+1].strip()
                    if re.match(r'\\d+', next_val):
                        detail['ics_code'] = next_val

            # 主管部门
            elif '主管部门' in line:
                match = re.search(r'主管部门[：:]+(.+)', line)
                if match:
                    detail['department'] = match.group(1).strip()
                elif i + 1 < len(lines):
                    detail['department'] = lines[i+1].strip()

            # 归口部门
            elif '归口部门' in line:
                match = re.search(r'归口部门[：:]+(.+)', line)
                if match:
                    detail['technical_dept'] = match.group(1).strip()
                elif i + 1 < len(lines):
                    detail['technical_dept'] = lines[i+1].strip()

            # 发布单位
            elif '发布单位' in line:
                match = re.search(r'发布单位[：:]+(.+)', line)
                if match:
                    detail['publisher'] = match.group(1).strip()
                elif i + 1 < len(lines):
                    detail['publisher'] = lines[i+1].strip()

        return detail
    except Exception as e:
        log(f"[详情] 获取 {hcno} 失败: {e}")
        return {}

def crawl_api2(keyword, page_size=50):
    """接口2: HTML解析获取列表和详情"""
    log(f"[接口2] 开始爬取关键词: {keyword}")

    params = {
        'r': str(time.time()),
        'page': 1,
        'pageSize': page_size,
        'p.p1': 0,
        'p.p2': keyword,
        'p.p7': 3,
        'p.p90': 'circulation_date',
        'p.p91': 'desc'
    }

    try:
        resp = requests.get(API2_URL, params=params, headers=HEADERS, timeout=30)
        resp.encoding = 'utf-8'
        soup = BeautifulSoup(resp.text, 'html.parser')

        # 解析总数
        total = 0
        for div in soup.find_all('div'):
            text = div.get_text(strip=True)
            match = re.search(r'共(\\d+)条标准', text)
            if match:
                total = int(match.group(1))
                break

        log(f"[接口2] 总记录数: {total}")
        if total == 0:
            return []

        all_data = []
        total_pages = (total + page_size - 1) // page_size

        for page in range(1, min(total_pages + 1, 51)):
            params['page'] = page
            params['r'] = str(time.time())

            try:
                resp = requests.get(API2_URL, params=params, headers=HEADERS, timeout=30)
                resp.encoding = 'utf-8'
                soup = BeautifulSoup(resp.text, 'html.parser')

                records = []
                for table in soup.find_all('table'):
                    rows = table.find_all('tr')
                    if rows and '标准号' in rows[0].get_text():
                        for row in rows[1:]:
                            cells = row.find_all('td')
                            if len(cells) >= 8:
                                onclick = ''
                                for cell in cells:
                                    a = cell.find('a')
                                    if a and a.get('onclick'):
                                        onclick = a.get('onclick')
                                        break
                                hcno_match = re.search(r"showInfo\\s*\\(\\s*['\\"]([^'\\"]+)['\\"]", onclick)
                                hcno = hcno_match.group(1) if hcno_match else ''

                                records.append({
                                    'standard_no': cells[1].get_text(strip=True),
                                    'standard_name': cells[3].get_text(strip=True),
                                    'is_adopted': cells[2].get_text(strip=True),
                                    'category': cells[4].get_text(strip=True),
                                    'status': cells[5].get_text(strip=True),
                                    'publish_date': clean_date(cells[6].get_text(strip=True)),
                                    'implement_date': clean_date(cells[7].get_text(strip=True)),
                                    'hcno': hcno,
                                })
                        break

                all_data.extend(records)
                log(f"[接口2] 第{page}/{total_pages}页: {len(records)}条")
                time.sleep(0.3)

            except Exception as e:
                log(f"[接口2] 第{page}页错误: {e}")

        log(f"[接口2] 总计: {len(all_data)} 条")
        return all_data

    except Exception as e:
        log(f"[接口2] 错误: {e}")
        return []

def merge_data(api1_data, api2_data):
    """合并两个接口数据，按标准号取并集，并获取详情补充字段"""
    result = []
    api1_map = {d.get('C_STD_CODE', ''): d for d in api1_data}
    api2_map = {d.get('standard_no', ''): d for d in api2_data}

    # 所有标准号取并集
    all_std_nos = set(api1_map.keys()) | set(api2_map.keys())
    total = len(all_std_nos)
    log(f"[合并] 需处理 {total} 条记录（接口1: {len(api1_data)}, 接口2: {len(api2_data)}）")

    detail_count = 0
    for idx, std_no in enumerate(all_std_nos, 1):
        api1_item = api1_map.get(std_no, {})
        api2_item = api2_map.get(std_no, {})

        # 获取详情补充字段（仅当有hcno时）
        hcno = api2_item.get('hcno', '')
        detail = {}
        if hcno:
            detail_count += 1
            if idx % 10 == 1:
                log(f"[详情] 进度 {idx}/{total}, 已获取 {detail_count} 个详情页...")
            time.sleep(0.2)  # 降低延迟
            detail = get_detail_api2(hcno)

        merged = {
            'link1': DETAIL1_URL + api1_item.get('id', '') if api1_item.get('id') else '',
            'link2': DETAIL2_URL + hcno if hcno else '',
            'standard_no': std_no,
            'standard_name': api2_item.get('standard_name') or api1_item.get('C_C_NAME', '') or api1_item.get('C_STD_NAME', ''),
            'english_name': detail.get('english_name', '') or api1_item.get('ENGLISH_NAME', ''),
            'publish_date': api2_item.get('publish_date') or clean_date(api1_item.get('ISSUE_DATE', '')),
            'implement_date': api2_item.get('implement_date') or clean_date(api1_item.get('ACT_DATE', '')),
            'status': api2_item.get('status') or api1_item.get('STATE', ''),
            'nature': api1_item.get('STD_NATURE', ''),
            'category': api2_item.get('category') or api1_item.get('STD_TYPE', ''),
            'is_adopted': api2_item.get('is_adopted', ''),
            'ccs_code': detail.get('ccs_code', '') or api1_item.get('CCS_CODE', ''),
            'ics_code': detail.get('ics_code', '') or api1_item.get('ICS_CODE', ''),
            'department': detail.get('department', '') or api1_item.get('DEPT_NAME', ''),
            'technical_dept': detail.get('technical_dept', '') or api1_item.get('TECH_DEPT', ''),
            'publisher': detail.get('publisher', '') or api1_item.get('PUBLISH_DEPT', ''),
            'description': '',
            'download_url': '',
            'hcno': hcno,
        }
        result.append(merged)

    log(f"[合并] 完成，共获取 {detail_count} 个详情页")
    return result

def main(keyword="电力"):
    log(f"\\n{'='*50}")
    log(f"开始爬取国家标准 - 关键词: {keyword}")
    log(f"{'='*50}\\n")

    # 爬取两个接口
    api1_data = crawl_api1(keyword)
    api2_data = crawl_api2(keyword)

    # 合并数据
    result = merge_data(api1_data, api2_data)

    log(f"\\n合并后总计: {len(result)} 条国标数据")
    print(json.dumps(result, ensure_ascii=False))

if __name__ == "__main__":
    keyword = os.environ.get('KEYWORD', '电力')
    main(keyword)
`;

// 行标爬取脚本
const industryScript = `#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
行业标准爬虫 - 电力行业标准数据爬取 (hbba.sacinfo.org.cn)
"""

import requests
import json
import time
import sys
import os
from datetime import datetime

BASE_URL = "https://hbba.sacinfo.org.cn/stdQueryList"
DETAIL_URL = "https://hbba.sacinfo.org.cn/stdDetail/"

HEADERS = {
    'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36',
    'Accept': 'application/json, text/javascript, */*',
    'Content-Type': 'application/x-www-form-urlencoded',
}

def log(msg):
    """日志输出到stderr"""
    print(msg, file=sys.stderr)

def convert_timestamp(ts):
    """转换时间戳为日期"""
    if ts:
        try:
            return datetime.fromtimestamp(ts / 1000).strftime('%Y-%m-%d')
        except:
            return ''
    return ''

def crawl_industry(keyword, page_size=100):
    """爬取行业标准数据"""
    log(f"开始爬取行业标准 - 关键词: {keyword}")
    all_data = []
    page = 1

    while True:
        data = {
            'current': page,
            'size': page_size,
            'key': keyword,
            'ministry': '',
            'industry': '电力',
            'pubdate': '-36',  # 近36个月
            'date': '',
            'status[]': ['即将实施', '现行']
        }

        try:
            resp = requests.post(BASE_URL, data=data, headers=HEADERS, timeout=30)
            result = resp.json()
            records = result.get('records', [])

            if not records:
                break

            for r in records:
                all_data.append({
                    'standard_no': r.get('code', ''),
                    'standard_name': r.get('chName', ''),
                    'publish_date': convert_timestamp(r.get('issueDate')),
                    'implement_date': convert_timestamp(r.get('actDate')),
                    'status': r.get('status', ''),
                    'approve_dept': r.get('chargeDept', ''),
                    'replace_standard': r.get('reviseStdCodes', ''),
                    'detail_link': DETAIL_URL + r.get('pk', ''),
                })

            log(f"第{page}页: {len(records)}条")

            if len(records) < page_size:
                break
            page += 1
            time.sleep(0.3)

        except Exception as e:
            log(f"错误: {e}")
            break

    return all_data

def main(keyword="电力"):
    """主函数"""
    log(f"\\n{'='*50}")
    log(f"开始爬取行业标准 - 关键词: {keyword}")
    log(f"{'='*50}\\n")

    data = crawl_industry(keyword)
    log(f"\\n总计获取: {len(data)} 条行标数据")
    # 只输出JSON到stdout
    print(json.dumps(data, ensure_ascii=False))

if __name__ == "__main__":
    keyword = os.environ.get('KEYWORD', '电力')
    data = main(keyword)
`;

// 地标爬取脚本（dbba.sacinfo.org.cn 地方标准）
const localScript = `#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
地方标准爬虫 - dbba.sacinfo.org.cn
支持按省份(ministry参数)筛选，如：henanzjj(河南省)、shanghaizjj(上海)等
"""

import requests
import json
import time
import sys
import os
from datetime import datetime

BASE_URL = "https://dbba.sacinfo.org.cn/stdQueryList"
DETAIL_URL = "https://dbba.sacinfo.org.cn/stdDetail/"

HEADERS = {
    'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36',
    'Accept': 'application/json, text/javascript, */*',
    'Content-Type': 'application/x-www-form-urlencoded',
}

def log(msg):
    """日志输出到stderr"""
    print(msg, file=sys.stderr)

def convert_timestamp(ts):
    """转换时间戳为日期"""
    if ts:
        try:
            return datetime.fromtimestamp(ts / 1000).strftime('%Y-%m-%d')
        except:
            return ''
    return ''

def crawl_local(keyword, ministry='', page_size=100):
    """爬取地方标准数据
    Args:
        keyword: 搜索关键词，如：电力
        ministry: 省份代码，如：henanzjj(河南)、shanghaizjj(上海)、beijing(北京)等
                  留空则搜索所有省份
    """
    log(f"开始爬取地方标准 - 关键词: {keyword}, 省份: {ministry or '全部'}")
    all_data = []
    page = 1

    while True:
        data = {
            'current': page,
            'size': page_size,
            'key': keyword,
            'ministry': ministry,  # 留空搜索所有省份，或指定省份代码
            'industry': '',
            'pubdate': '-36',  # 近36个月
            'date': '',
            'status': ''  # 留空获取所有状态：现行、废止、即将实施等
        }

        try:
            resp = requests.post(BASE_URL, data=data, headers=HEADERS, timeout=30)
            result = resp.json()
            records = result.get('records', [])

            if not records:
                break

            for r in records:
                all_data.append({
                    'standard_no': r.get('code', ''),
                    'standard_name': r.get('chName', ''),
                    'publish_date': convert_timestamp(r.get('issueDate')),
                    'implement_date': convert_timestamp(r.get('actDate')),
                    'status': r.get('status', ''),
                    'department': r.get('chargeDept', ''),  # 批准部门
                    'publisher': r.get('industry', ''),  # 所属省份
                    'detail_link': DETAIL_URL + r.get('pk', ''),
                    'pk': r.get('pk', ''),
                })

            log(f"第{page}页: {len(records)}条")

            # 检查是否还有更多页
            total_pages = result.get('pages', 1)
            if page >= total_pages or len(records) < page_size:
                break
            page += 1
            time.sleep(0.3)

        except Exception as e:
            log(f"错误: {e}")
            break

    return all_data

def main(keyword="电力"):
    """主函数"""
    log(f"\\n{'='*50}")
    log(f"开始爬取地方标准 - 关键词: {keyword}")
    log(f"{'='*50}\\n")

    data = crawl_local(keyword)
    log(f"\\n总计获取: {len(data)} 条地标数据")
    # 只输出JSON到stdout
    print(json.dumps(data, ensure_ascii=False))

if __name__ == "__main__":
    keyword = os.environ.get('KEYWORD', '电力')
    data = main(keyword)
`;

const scripts: Record<string, string> = {
  national: nationalScript,
  industry: industryScript,
  local: localScript,
};

// 初始化编辑器
onMounted(async () => {
  await nextTick();
  if (editorContainer.value) {
    editor = monaco.editor.create(editorContainer.value, {
      value: scripts[dataSource.value],
      language: "python",
      theme: "vs-dark",
      minimap: { enabled: false },
      fontSize: 13,
      scrollBeyondLastLine: false,
      automaticLayout: true,
    });
  }
});

onUnmounted(() => {
  editor?.dispose();
  editor = null;
});

watch(dataSource, (val) => {
  if (editor) {
    editor.setValue(scripts[val]);
  }
  // 自动锁定导入目标
  importTarget.value = val;
  // 清空之前的结果
  resultData.value = [];
});

// 执行爬取
const handleExecute = async () => {
  if (!queryKeyword.value) {
    ElMessage.warning("请输入查询关键词");
    return;
  }

  executing.value = true;
  logs.value = [];
  resultData.value = [];
  addLog("开始执行爬取脚本...");

  try {
    const script = editor?.getValue() || scripts[dataSource.value];
    const res = (await crawlerApi.execute(script, queryKeyword.value)) as any;

    // 后端返回结构: { data: { data: [...], log: "...", success: true } }
    const data = res.data?.data || res.data;

    // 先显示执行日志
    if (res.data?.log) {
      const logLines = res.data.log.split("\n").filter((l: string) => l.trim());
      logLines.forEach((line: string) => addLog(line));
    }

    if (data && Array.isArray(data) && data.length > 0) {
      resultData.value = data;
      addLog(`爬取完成，获取 ${resultData.value.length} 条数据`);
    } else {
      addLog("警告: 爬取完成但未获取到数据");
    }
  } catch (e: any) {
    addLog(`执行失败: ${e.message}`, "error");
    // 错误时不显示模拟数据，只显示错误日志
    resultData.value = [];
  } finally {
    executing.value = false;
  }
};

// 导入数据
const handleImport = async () => {
  if (resultData.value.length === 0) return;
  importing.value = true;
  addLog("开始导入数据...");
  try {
    await crawlerApi.import(resultData.value, importTarget.value);
    addLog(`导入成功，共 ${resultData.value.length} 条`);
    ElMessage.success("导入成功");
    resultData.value = [];
  } catch (e: any) {
    addLog(`导入失败: ${e.message}`, "error");
  } finally {
    importing.value = false;
  }
};

const addLog = (message: string, type: "info" | "error" = "info") => {
  logs.value.push({ time: new Date().toLocaleTimeString(), message, type });
  // 自动滚动到底部
  nextTick(() => {
    if (terminalBody.value) {
      terminalBody.value.scrollTop = terminalBody.value.scrollHeight;
    }
  });
};
</script>

<style scoped>
/* 终端样式 */
.terminal-header {
  background: #323232;
  padding: 8px 12px;
  display: flex;
  align-items: center;
  border-bottom: 1px solid #404040;
  border-radius: 0;
}

.terminal-buttons {
  display: flex;
  gap: 6px;
}

.terminal-btn {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.terminal-btn-red {
  background: #ff5f56;
}

.terminal-btn-yellow {
  background: #ffbd2e;
}

.terminal-btn-green {
  background: #27ca40;
}

.terminal-title {
  color: #888;
  font-size: 12px;
  margin-left: 12px;
  font-family: "Monaco", "Menlo", "Consolas", monospace;
}

.terminal-body {
  background: #1e1e1e;
  padding: 12px 16px;
  overflow-y: auto;
  font-family: "Monaco", "Menlo", "Consolas", "Courier New", monospace;
  font-size: 12px;
  line-height: 1.6;
}

.terminal-placeholder {
  color: #666;
  font-style: italic;
}

.terminal-line {
  margin-bottom: 2px;
}

.terminal-prompt {
  color: #4ec9b0;
  margin-bottom: 8px;
}

.terminal-dollar {
  color: #6a9955;
}

.terminal-output {
  color: #d4d4d4;
}

.terminal-error {
  color: #f14c4c;
}

.terminal-cursor {
  margin-top: 4px;
}

.cursor-blink {
  color: #d4d4d4;
  animation: blink 1s step-end infinite;
}

@keyframes blink {
  0%,
  100% {
    opacity: 1;
  }
  50% {
    opacity: 0;
  }
}

/* 滚动条样式 */
.terminal-body::-webkit-scrollbar {
  width: 6px;
}

.terminal-body::-webkit-scrollbar-track {
  background: #2d2d2d;
}

.terminal-body::-webkit-scrollbar-thumb {
  background: #555;
  border-radius: 3px;
}

.terminal-body::-webkit-scrollbar-thumb:hover {
  background: #666;
}
</style>

