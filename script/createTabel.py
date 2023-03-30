import pandas as pd
import os
if __name__ == '__main__':
    # 导入 Excel 文件
    df = pd.read_excel('/Users/cboy/goProject/awesomeProject10/gochat/script/goods list.xls')

    # 选择需要的列
    df = df[['发货日期', '公司名称', '金额']]

    # 计算每年的总金额
    df['年份'] = df['发货日期'].apply(lambda x: x.year)
    df = df.groupby(['年份', '公司名称']).sum()

    # 打印结果
    filepath = '/Users/cboy/goProject/awesomeProject10/gochat/script/result.xls'  # 修改成实际的文件路径
    new_filepath = os.path.splitext(filepath)[0] + '_汇总.xlsx'
    df.to_excel(new_filepath)
