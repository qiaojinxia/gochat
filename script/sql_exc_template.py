code = '''
import init
import pandas as pd
import openpyxl
from pandasql import sqldf
try:
{table_codes}
    
    # 使用pandasSql分析表格，添加新列
    pysqldf = lambda q: sqldf(q, globals())
    
    query = """
    {query}
    """
    result = pysqldf(query)
    
    # 文件名加上_result后缀
    result_file_path = '{file_path}' + '_result.xlsx'
    # 保存文件
    result.to_excel('data/' + result_file_path, index=False)
    
    # 输出文件名
    print(result_file_path)
except Exception as e:
    # 输出错误信息到控制台
    print('Error:', e)
'''
import sys

# 从控制台接收参数1文件目录，参数2 sql语句
tables = eval(sys.argv[1])
query = sys.argv[2]
file_path = sys.argv[3]

table_codes = ''
for table_name, table_path in tables.items():
    table_codes += f"    {table_name} = pd.read_excel('{table_path}',engine='openpyxl_wo_formatting')\n"

code = code.replace("{table_codes}", table_codes)
code = code.replace("{file_path}", file_path)
code = code.replace("{query}", query)


exec(code)