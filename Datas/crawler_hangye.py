#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
行业标准爬虫 - 电力相关行业标准数据爬取（hbba.sacinfo.org.cn）
"""

import requests
import pandas as pd
import time
import re
from bs4 import BeautifulSoup
from datetime import datetime

# 基础URL
BASE_URL = "https://hbba.sacinfo.org.cn/stdQueryList"
DETAIL_BASE_URL = "https://hbba.sacinfo.org.cn/stdDetail/"

# 请求头
HEADERS = {
    'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36',
    'Accept': 'application/json, text/javascript, */*',
    'Accept-Language': 'zh-CN,zh;q=0.9,en;q=0.8',
    'Content-Type': 'application/x-www-form-urlencoded',
}

def convert_timestamp(ts):
    """转换时间戳为日期字符串"""
    if ts:
        try:
            dt = datetime.fromtimestamp(ts / 1000)
            return dt.strftime('%Y-%m-%d')
        except:
            return ''
    return ''

def get_page_data(page=1, page_size=100, delay=0.5):
    """获取单页数据（POST请求）"""
    data = {
        'current': page,
        'size': page_size,
        'key': '电力',
        'ministry': '',
        'industry': '电力',
        'pubdate': '-36',
        'date': '',
        'status[]': ['即将实施', '现行']
    }

    try:
        response = requests.post(BASE_URL, data=data, headers=HEADERS, timeout=30)
        response.raise_for_status()
        return response.json()
    except Exception as e:
        print(f"请求第 {page} 页失败: {e}")
        return None

def get_detail_info(pk):
    """获取详情页信息"""
    detail_url = DETAIL_BASE_URL + pk

    try:
        headers = HEADERS.copy()
        headers['Accept'] = 'text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8'
        headers.pop('Content-Type', None)

        response = requests.get(detail_url, headers=headers, timeout=30)
        response.encoding = 'utf-8'
        soup = BeautifulSoup(response.text, 'html.parser')

        text = soup.get_text(separator='\n', strip=True)
        lines = text.split('\n')

        detail_info = {
            'detail_url': detail_url,
            '标准号': '',
            '发布日期': '',
            '实施日期': '',
            '制修订': '',
            '中国标准分类号': '',
            '国际标准分类号': '',
            '技术归口': '',
            '批准发布部门': '',
            '行业分类': '',
            '标准类别': ''
        }

        # 按行解析字段
        for i, line in enumerate(lines):
            line = line.strip()

            if '标准号' in line and len(line) < 20:
                # 可能在同一行或下一行
                match = re.search(r'标准号[：:\s]+([A-Z]+/T?\s*[\d\-]+)', line)
                if match:
                    detail_info['标准号'] = match.group(1).strip()
                elif i + 1 < len(lines):
                    next_line = lines[i+1].strip()
                    if re.match(r'[A-Z]+/T?\s*[\d—\-]+', next_line):
                        detail_info['标准号'] = next_line.replace('—', '-')

            elif '发布日期' in line and len(line) < 15:
                match = re.search(r'(\d{4}[-/]\d{2}[-/]\d{2})', line)
                if match:
                    detail_info['发布日期'] = match.group(1)
                elif i + 1 < len(lines):
                    next_line = lines[i+1].strip()
                    match = re.search(r'(\d{4}[-/]\d{2}[-/]\d{2})', next_line)
                    if match:
                        detail_info['发布日期'] = match.group(1)

            elif '实施日期' in line and len(line) < 15:
                match = re.search(r'(\d{4}[-/]\d{2}[-/]\d{2})', line)
                if match:
                    detail_info['实施日期'] = match.group(1)
                elif i + 1 < len(lines):
                    next_line = lines[i+1].strip()
                    match = re.search(r'(\d{4}[-/]\d{2}[-/]\d{2})', next_line)
                    if match:
                        detail_info['实施日期'] = match.group(1)

            elif '制修订' in line and len(line) < 10:
                if i + 1 < len(lines):
                    next_line = lines[i+1].strip()
                    if next_line in ['制定', '修订', '修改']:
                        detail_info['制修订'] = next_line

            elif '中国标准分类号' in line and len(line) < 25:
                match = re.search(r'([A-Z]\d+)', line)
                if match:
                    detail_info['中国标准分类号'] = match.group(1)
                elif i + 1 < len(lines):
                    next_line = lines[i+1].strip()
                    if re.match(r'[A-Z]\d+', next_line):
                        detail_info['中国标准分类号'] = next_line

            elif '国际标准分类号' in line and len(line) < 25:
                match = re.search(r'(\d+\.?\d*\.?\d*)', line)
                if match:
                    detail_info['国际标准分类号'] = match.group(1)
                elif i + 1 < len(lines):
                    next_line = lines[i+1].strip()
                    if re.match(r'\d+', next_line):
                        detail_info['国际标准分类号'] = next_line

            elif '技术归口' in line and len(line) < 15:
                if i + 1 < len(lines):
                    next_line = lines[i+1].strip()
                    if len(next_line) > 5 and len(next_line) < 100:
                        detail_info['技术归口'] = next_line

            elif '批准发布部门' in line and len(line) < 20:
                if i + 1 < len(lines):
                    next_line = lines[i+1].strip()
                    if len(next_line) > 3 and len(next_line) < 50:
                        detail_info['批准发布部门'] = next_line

            elif '行业分类' in line and len(line) < 15:
                if i + 1 < len(lines):
                    next_line = lines[i+1].strip()
                    if len(next_line) > 5 and len(next_line) < 50:
                        detail_info['行业分类'] = next_line

            elif '标准类别' in line and len(line) < 15:
                if i + 1 < len(lines):
                    next_line = lines[i+1].strip()
                    if len(next_line) > 2 and len(next_line) < 20:
                        detail_info['标准类别'] = next_line

        return detail_info

    except Exception as e:
        print(f"获取详情 {pk} 失败: {e}")
        return {'detail_url': DETAIL_BASE_URL + pk, 'error': str(e)}

def crawl_all_data(page_size=100, delay=0.5):
    """爬取所有数据"""
    all_data = []

    print("开始爬取行业标准电力相关数据...")

    # 获取第一页
    first_page = get_page_data(1, page_size)
    if not first_page:
        print("无法获取第一页数据")
        return []

    total = first_page.get('total', 0)
    records = first_page.get('records', [])
    print(f"总记录数: {total}")

    all_data.extend(records)

    # 计算总页数并继续爬取
    total_pages = (total + page_size - 1) // page_size

    for page in range(2, total_pages + 1):
        print(f"处理第 {page}/{total_pages} 页...")
        time.sleep(delay)

        page_data = get_page_data(page, page_size)
        if page_data:
            records = page_data.get('records', [])
            all_data.extend(records)

    print(f"共获取 {len(all_data)} 条记录")
    return all_data

def process_and_export(all_data, output_file="行业标准电力数据.xlsx"):
    """处理数据并导出到Excel"""
    if not all_data:
        print("没有数据可导出")
        return None

    print(f"\n开始获取详情信息...")

    result_data = []
    total = len(all_data)

    for index, record in enumerate(all_data, 1):
        if index % 20 == 0 or index == 1:
            print(f"处理第 {index}/{total} 条记录...")

        pk = record.get('pk', '')
        if not pk:
            continue

        # 从列表数据提取基础信息
        basic_info = {
            '标准号': record.get('code', ''),
            '标准中文名称': record.get('chName', ''),
            '发布日期': convert_timestamp(record.get('issueDate')),
            '实施日期': convert_timestamp(record.get('actDate')),
            '标准状态': record.get('status', ''),
            '行业': record.get('industry', ''),
            '批准发布部门': record.get('chargeDept', ''),
            '代替标准': record.get('reviseStdCodes', '')
        }

        # 获取详情补充信息
        time.sleep(0.3)
        detail_info = get_detail_info(pk)

        # 合并数据（详情信息补充基础信息）
        row = {
            '详情链接': detail_info.get('detail_url', DETAIL_BASE_URL + pk),
            '标准号': detail_info.get('标准号') or basic_info['标准号'],
            '标准中文名称': basic_info['标准中文名称'],
            '发布日期': detail_info.get('发布日期') or basic_info['发布日期'],
            '实施日期': detail_info.get('实施日期') or basic_info['实施日期'],
            '制修订': detail_info.get('制修订', ''),
            '中国标准分类号': detail_info.get('中国标准分类号', ''),
            '国际标准分类号': detail_info.get('国际标准分类号', ''),
            '技术归口': detail_info.get('技术归口', ''),
            '批准发布部门': detail_info.get('批准发布部门') or basic_info['批准发布部门'],
            '行业分类': detail_info.get('行业分类') or basic_info['行业'],
            '标准类别': detail_info.get('标准类别', ''),
            '标准状态': basic_info['标准状态'],
            '代替标准': basic_info['代替标准']
        }

        result_data.append(row)

    # 创建DataFrame并导出
    df = pd.DataFrame(result_data)
    df.to_excel(output_file, index=False, engine='openpyxl')

    print(f"\n数据已导出到: {output_file}")
    print(f"共 {len(df)} 条记录")

    return df

def main():
    """主函数"""
    all_data = crawl_all_data(page_size=100, delay=0.5)

    if all_data:
        df = process_and_export(all_data, "行业标准电力数据.xlsx")

        if df is not None and len(df) > 0:
            print("\n数据预览:")
            pd.set_option('display.max_columns', None)
            pd.set_option('display.width', None)
            pd.set_option('display.max_colwidth', 30)
            print(df.head(5))

if __name__ == "__main__":
    main()