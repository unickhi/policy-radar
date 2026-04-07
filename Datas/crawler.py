#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
国家标准爬虫 - 电力相关标准数据爬取
"""

import requests
import pandas as pd
import time
import re
from bs4 import BeautifulSoup

# API基础URL
BASE_API_URL = "https://std.samr.gov.cn/gb/search/gbQueryPage"
DETAIL_BASE_URL = "https://std.samr.gov.cn/gb/search/gbDetailed?id="

# 请求头
HEADERS = {
    'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36',
    'Accept': 'application/json, text/javascript, */*; q=0.01',
    'Accept-Language': 'zh-CN,zh;q=0.9,en;q=0.8',
    'Referer': 'https://std.samr.gov.cn/gb/search/gbHomePage',
}

def get_page_data(search_text="电力", page_number=1, page_size=50):
    """获取单页数据"""
    params = {
        'searchText': search_text,
        'ics': '',
        'state': '',
        'ISSUE_DATE': '-36',
        'sortOrder': 'asc',
        'pageSize': page_size,
        'pageNumber': page_number,
        '_': int(time.time() * 1000)
    }

    try:
        response = requests.get(BASE_API_URL, params=params, headers=HEADERS, timeout=30)
        response.raise_for_status()
        return response.json()
    except Exception as e:
        print(f"请求第 {page_number} 页失败: {e}")
        return None

def clean_html_tags(text):
    """清理HTML标签"""
    if not text:
        return ''
    # 移除<sacinfo>标签
    text = re.sub(r'</?sacinfo>', '', text)
    return text.strip()

def get_detail_summary(detail_id):
    """获取详情页摘要信息"""
    detail_url = DETAIL_BASE_URL + detail_id

    try:
        response = requests.get(detail_url, headers=HEADERS, timeout=30)
        response.encoding = 'utf-8'
        soup = BeautifulSoup(response.text, 'html.parser')

        text = soup.get_text(separator='\n', strip=True)
        lines = text.split('\n')

        summary_parts = []

        # 提取归口单位 (格式: TC44（全国变压器标准化技术委员会）归口)
        for i, line in enumerate(lines):
            if '归口' in line and 'TC' in line:
                # 合并前后几行获取完整信息
                context = ' '.join(lines[max(0, i-1):min(len(lines), i+3)])
                context = re.sub(r'\s+', ' ', context)
                # 提取归口单位名称
                match = re.search(r'（([^）]+)）归口', context)
                if match:
                    summary_parts.append(f"归口单位: {match.group(1)}")
                break

        # 提取主管部门
        for i, line in enumerate(lines):
            if '主管部门' in line:
                # 获取主管部门名称
                context = ' '.join(lines[i:min(len(lines), i+3)])
                match = re.search(r'主管部门[为\s]*([^。\n]+)', context)
                if match:
                    name = match.group(1).strip()
                    if name and name not in ['为', '']:
                        summary_parts.append(f"主管部门: {name}")
                break

        # 提取起草单位
        drafting_units = []
        for i, line in enumerate(lines):
            if '起草单位' in line:
                # 获取后续几行的起草单位
                for j in range(i+1, min(len(lines), i+10)):
                    unit = lines[j].strip()
                    if unit and len(unit) > 3 and '起草' not in unit and '人' not in unit:
                        drafting_units.append(unit)
                    if '起草人' in unit:
                        break
                if drafting_units:
                    summary_parts.append(f"起草单位: {', '.join(drafting_units[:3])}")
                break

        # 提取采标情况
        for i, line in enumerate(lines):
            if '采标情况' in line or '修改采用' in line or '等同采用' in line:
                context = ' '.join(lines[i:min(len(lines), i+3)])
                match = re.search(r'(修改采用|等同采用)[^。]+', context)
                if match:
                    summary_parts.append(f"采标: {match.group().strip()}")
                break

        # 提取代替标准
        for i, line in enumerate(lines):
            if '代替了以下标准' in line or '代替标准' in line:
                for j in range(i+1, min(len(lines), i+5)):
                    if re.match(r'[A-Z]+/T?\s*\d+', lines[j]):
                        summary_parts.append(f"代替标准: {lines[j].strip()}")
                        break
                break

        if summary_parts:
            return '; '.join(summary_parts)

        return ''

    except Exception as e:
        return f'获取详情失败: {str(e)[:50]}'

def crawl_all_data(search_text="电力", page_size=50, delay=0.5):
    """爬取所有数据"""
    all_data = []

    print(f"开始爬取'{search_text}'相关标准数据...")

    # 先获取第一页确定总页数
    first_page = get_page_data(search_text, 1, page_size)
    if not first_page:
        print("无法获取第一页数据,退出爬取")
        return []

    total_records = first_page.get('total', 0)
    records = first_page.get('rows', [])

    if total_records:
        total_pages = (total_records + page_size - 1) // page_size
        print(f"总记录数: {total_records}, 总页数: {total_pages}")
    else:
        total_pages = 1 if len(records) < page_size else 100
        print(f"未找到总记录数,将爬取直到没有数据")

    # 处理第一页
    print(f"\n处理第 1 页...")
    all_data.extend(records)

    # 继续爬取后续页面
    page_number = 2
    while True:
        if total_records and page_number > total_pages:
            break

        print(f"处理第 {page_number} 页...")
        time.sleep(delay)

        page_data = get_page_data(search_text, page_number, page_size)
        if not page_data:
            break

        records = page_data.get('rows', [])
        if not records:
            print(f"第 {page_number} 页没有数据,停止爬取")
            break

        all_data.extend(records)
        page_number += 1

        if page_number > 200:
            print("达到最大页数限制(200页)")
            break

    print(f"\n共获取 {len(all_data)} 条记录")
    return all_data

def process_and_export(all_data, output_file="电力标准数据.xlsx"):
    """处理数据并导出到Excel"""
    if not all_data:
        print("没有数据可导出")
        return None

    print(f"\n开始处理并获取详情信息...")

    result_data = []
    total = len(all_data)

    for index, record in enumerate(all_data, 1):
        if index % 10 == 0 or index == 1:
            print(f"处理第 {index}/{total} 条记录...")

        detail_id = record.get('id', '')
        if not detail_id:
            continue

        # 从API数据中提取字段
        detail_url = DETAIL_BASE_URL + detail_id

        # 清理名称中的HTML标签
        standard_name = clean_html_tags(record.get('C_C_NAME', ''))

        # 获取详情摘要
        time.sleep(0.3)  # 添加延迟避免请求过快
        summary = get_detail_summary(detail_id)

        row = {
            '详情链接': detail_url,
            '标准号': record.get('C_STD_CODE', ''),
            '标准中文名称': standard_name,
            '发布日期': record.get('ISSUE_DATE', ''),
            '实施日期': record.get('ACT_DATE', ''),
            '标准状态': record.get('STATE', ''),
            '标准性质': record.get('STD_NATURE', ''),
            '详情简介': summary
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
    # 爬取所有数据
    all_data = crawl_all_data(search_text="电力", page_size=50, delay=0.3)

    # 处理并导出
    if all_data:
        df = process_and_export(all_data, "电力标准数据.xlsx")

        if df is not None and len(df) > 0:
            print("\n数据预览:")
            pd.set_option('display.max_columns', None)
            pd.set_option('display.width', None)
            pd.set_option('display.max_colwidth', 50)
            print(df.head(10))

if __name__ == "__main__":
    main()