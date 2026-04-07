#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
国家标准数据合并脚本 - 合并接口1和接口2数据
"""

import pandas as pd

def merge_national_standards():
    """合并两个国家标准数据源"""

    # 读取两个数据文件
    df1 = pd.read_excel('电力标准数据.xlsx')
    df2 = pd.read_excel('国家标准电力数据_接口2.xlsx')

    print("=== 数据源信息 ===")
    print(f"接口1数据: {len(df1)} 条")
    print(f"接口2数据: {len(df2)} 条")

    # 重命名链接字段
    df1 = df1.rename(columns={'详情链接': '链接1'})
    df2 = df2.rename(columns={'详情链接': '链接2'})

    # 移除df1中与df2重复的字段（保留链接1和详情简介）
    df1_keep_cols = ['标准号', '链接1', '详情简介', '标准性质']

    # df2保留的字段
    df2_keep_cols = ['标准号', '链接2', '标准中文名称', '英文标准名称', '是否采标',
                     '类别', '标准状态', '发布日期', '实施日期',
                     '中国标准分类号', '国际标准分类号', '主管部门', '归口部门', '发布单位']

    # 提取需要的列
    df1_subset = df1[df1_keep_cols]
    df2_subset = df2[df2_keep_cols]

    # 根据标准号合并
    merged_df = pd.merge(df1_subset, df2_subset, on='标准号', how='outer')

    # 整理列顺序
    final_columns = [
        '链接1', '链接2',
        '标准号', '标准中文名称', '英文标准名称',
        '发布日期', '实施日期', '标准状态', '标准性质', '类别', '是否采标',
        '中国标准分类号', '国际标准分类号',
        '主管部门', '归口部门', '发布单位',
        '详情简介'
    ]

    merged_df = merged_df[final_columns]

    # 导出合并后的数据
    output_file = '国家标准电力数据_合并.xlsx'
    merged_df.to_excel(output_file, index=False, engine='openpyxl')

    print(f"\n=== 合并结果 ===")
    print(f"合并后数据: {len(merged_df)} 条")
    print(f"输出文件: {output_file}")
    print(f"\n字段列表 ({len(final_columns)} 个):")
    for i, col in enumerate(final_columns, 1):
        print(f"  {i}. {col}")

    # 数据预览
    print("\n=== 数据预览 ===")
    pd.set_option('display.max_columns', None)
    pd.set_option('display.width', None)
    pd.set_option('display.max_colwidth', 30)
    print(merged_df.head(5))

    # 检查空值情况
    print("\n=== 字段完整性检查 ===")
    for col in final_columns:
        null_count = merged_df[col].isna().sum()
        empty_count = (merged_df[col] == '').sum() if merged_df[col].dtype == 'object' else 0
        total_missing = null_count + empty_count
        fill_rate = (len(merged_df) - total_missing) / len(merged_df) * 100
        print(f"  {col}: 填充率 {fill_rate:.1f}% ({len(merged_df) - total_missing}/{len(merged_df)})")

    return merged_df

if __name__ == "__main__":
    merge_national_standards()