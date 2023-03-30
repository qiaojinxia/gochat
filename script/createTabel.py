import pandas as pd
import mysql.connector

# 从Excel文件中读取数据
df = pd.read_excel('path/to/file.xlsx')

# 建立与MySQL的连接
config = {
    'user': '<username>',
    'password': '<password>',
    'host': '<host>',
    'database': '<database>'
}
cnx = mysql.connector.connect(**config)

# 创建游标对象
cursor = cnx.cursor()

# 将数据插入MySQL中
df.to_sql(con=cnx, name='table_name', if_exists='replace', index=False)

# 提交更改并关闭游标和连接
cursor.close()
cnx.commit()
cnx.close()