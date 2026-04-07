#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
国家标准爬虫 - 电力相关标准数据爬取（接口2：openstd.samr.gov.cn）
"""

import requests
import pandas as pd
import time
import re
from bs4 import BeautifulSoup

# 基础URL
BASE_URL = "https://openstd.samr.gov.cn/bzgk/gb/std_list"
DETAIL_BASE_URL = "https://openstd.samr.gov.cn/bzgk/gb/newGbInfo?hcno="

# 请求头
HEADERS = {
    'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36',
    'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8',
    'Accept-Language': 'zh-CN,zh;q=0.9,en;q=0.8',
}

def get_page_data(page=1, page_size=50, delay=0.5):
    """获取单页数据（HTML格式）"""
    params = {
        'r': str(time.time()),
        'page': page,
        'pageSize': page_size,
        'p.p1': 0,
        'p.p2': '电力',
        'p.p7': 3,
        'p.p90': 'circulation_date',
        'p.p91': 'desc'
    }

    try:
        response = requests.get(BASE_URL, params=params, headers=HEADERS, timeout=30)
        response.encoding = 'utf-8'
        return response.text
    except Exception as e:
        print(f"请求第 {page} 页失败: {e}")
        return None

def parse_list_data(html):
    """解析HTML列表数据"""
    soup = BeautifulSoup(html, 'html.parser')

    # 找数据表格（包含"标准号"的表格）
    records = []
    for table in soup.find_all('table'):
        rows = table.find_all('tr')
        if rows and len(rows) > 1:
            # 检查是否是数据表格
            header_text = rows[0].get_text(strip=True)
            if '标准号' in header_text:
                for row in rows[1:]:
                    cells = row.find_all('td')
                    if len(cells) >= 8:
                        # 提取onclick中的hcno
                        onclick = ''
                        for cell in cells:
                            a = cell.find('a')
                            if a and a.get('onclick'):
                                onclick = a.get('onclick')
                                break

                        hcno_match = re.search(r"showInfo\s*\(\s*['\"]([^'\"]+)['\"]", onclick)
                        hcno = hcno_match.group(1) if hcno_match else ''

                        record = {
                            '序号': cells[0].get_text(strip=True),
                            '标准号': cells[1].get_text(strip=True),
                            '是否采标': cells[2].get_text(strip=True),
                            '标准名称': cells[3].get_text(strip=True),
                            '类别': cells[4].get_text(strip=True),
                            '状态': cells[5].get_text(strip=True),
                            '发布日期': cells[6].get_text(strip=True).replace(' 00:00:00.', ''),
                            '实施日期': cells[7].get_text(strip=True).replace(' 00:00:00.', ''),
                            'hcno': hcno
                        }
                        records.append(record)
                break

    # 提取总数信息
    total = 0
    for div in soup.find_all('div'):
        text = div.get_text(strip=True)
        match = re.search(r'共(\d+)条标准', text)
        if match:
            total = int(match.group(1))
            break

    return records, total

def get_detail_info(hcno):
    """获取详情页信息"""
    detail_url = DETAIL_BASE_URL + hcno

    try:
        response = requests.get(detail_url, headers=HEADERS, timeout=30)
        response.encoding = 'utf-8'
        soup = BeautifulSoup(response.text, 'html.parser')

        text = soup.get_text(separator='\n', strip=True)
        lines = text.split('\n')

        detail_info = {
            'detail_url': detail_url,
            '英文标准名称': '',
            '中国标准分类号': '',
            '国际标准分类号': '',
            '主管部门': '',
            '归口部门': '',
            '发布单位': ''
        }

        # 按行解析字段
        for i, line in enumerate(lines):
            line = line.strip()

            # 格式：字段名：值（同一行）
            if '英文标准名称' in line:
                match = re.search(r'英文标准名称[：:]+(.+)', line)
                if match:
                    detail_info['英文标准名称'] = match.group(1).strip()
                elif i + 1 < len(lines):
                    detail_info['英文标准名称'] = lines[i+1].strip()

            elif '中国标准分类号（CCS）' in line or '中国标准分类号' in line:
                match = re.search(r'[分类号]+[：:]+([A-Z]\d+)', line)
                if match:
                    detail_info['中国标准分类号'] = match.group(1)
                elif i + 1 < len(lines):
                    next_val = lines[i+1].strip()
                    if re.match(r'[A-Z]\d+', next_val):
                        detail_info['中国标准分类号'] = next_val

            elif '国际标准分类号（ICS）' in line or '国际标准分类号' in line:
                match = re.search(r'[分类号]+[：:]+(\d+\.?\d*)', line)
                if match:
                    detail_info['国际标准分类号'] = match.group(1)
                elif i + 1 < len(lines):
                    next_val = lines[i+1].strip()
                    if re.match(r'\d+', next_val):
                        detail_info['国际标准分类号'] = next_val

            elif '主管部门' in line:
                match = re.search(r'主管部门[：:]+(.+)', line)
                if match:
                    detail_info['主管部门'] = match.group(1).strip()
                elif i + 1 < len(lines):
                    detail_info['主管部门'] = lines[i+1].strip()

            elif '归口部门' in line:
                match = re.search(r'归口部门[：:]+(.+)', line)
                if match:
                    detail_info['归口部门'] = match.group(1).strip()
                elif i + 1 < len(lines):
                    detail_info['归口部门'] = lines[i+1].strip()

            elif '发布单位' in line:
                match = re.search(r'发布单位[：:]+(.+)', line)
                if match:
                    detail_info['发布单位'] = match.group(1).strip()
                elif i + 1 < len(lines):
                    detail_info['发布单位'] = lines[i+1].strip()

        return detail_info

    except Exception as e:
        print(f"获取详情 {hcno} 失败: {e}")
        return {'detail_url': DETAIL_BASE_URL + hcno, 'error': str(e)}

def crawl_all_data(page_size=50, delay=0.5):
    """爬取所有数据"""
    all_data = []

    print("开始爬取国家标准（接口2）电力相关数据...")

    # 获取第一页
    html = get_page_data(1, page_size)
    if not html:
        print("无法获取第一页数据")
        return []

    records, total = parse_list_data(html)
    print(f"总记录数: {total}")

    all_data.extend(records)

    # 计算总页数并继续爬取
    total_pages = (total + page_size - 1) // page_size

    for page in range(2, total_pages + 1):
        print(f"处理第 {page}/{total_pages} 页...")
        time.sleep(delay)

        html = get_page_data(page, page_size)
        if html:
            records, _ = parse_list_data(html)
            all_data.extend(records)

    print(f"共获取 {len(all_data)} 条记录")
    return all_data

def process_and_export(all_data, output_file="国家标准电力数据_接口2.xlsx"):
    """处理数据并导出到Excel"""
    if not all_data:
        print("没有数据可导出")
        return None

    print(f"\n开始获取详情信息...")

    result_data = []
    total = len(all_data)

    for index, record in enumerate(all_data, 1):
        if index % 10 == 0 or index == 1:
            print(f"处理第 {index}/{total} 条记录...")

        hcno = record.get('hcno', '')
        if not hcno:
            continue

        # 获取详情
        time.sleep(0.3)
        detail_info = get_detail_info(hcno)

        # 合并数据
        row = {
            '详情链接': detail_info.get('detail_url', DETAIL_BASE_URL + hcno),
            '标准号': record.get('标准号', ''),
            '标准中文名称': record.get('标准名称', ''),
            '英文标准名称': detail_info.get('英文标准名称', ''),
            '是否采标': record.get('是否采标', ''),
            '类别': record.get('类别', ''),
            '标准状态': record.get('状态', ''),
            '发布日期': record.get('发布日期', ''),
            '实施日期': record.get('实施日期', ''),
            '中国标准分类号': detail_info.get('中国标准分类号', ''),
            '国际标准分类号': detail_info.get('国际标准分类号', ''),
            '主管部门': detail_info.get('主管部门', ''),
            '归口部门': detail_info.get('归口部门', ''),
            '发布单位': detail_info.get('发布单位', '')
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
    all_data = crawl_all_data(page_size=50, delay=0.5)

    if all_data:
        df = process_and_export(all_data, "国家标准电力数据_接口2.xlsx")

        if df is not None and len(df) > 0:
            print("\n数据预览:")
            pd.set_option('display.max_columns', None)
            pd.set_option('display.width', None)
            pd.set_option('display.max_colwidth', 30)
            print(df.head(5))

if __name__ == "__main__":
    main()