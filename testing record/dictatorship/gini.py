import pandas as pd
import numpy as np
# 读取包含多个子表格的Excel文件
excel_file_path = 'd5.xlsx'
excel_data = pd.read_excel(excel_file_path, sheet_name=None)

# sheet_name=None 会将所有的子表格读取到一个字典中，其中字典的键是子表格的名称

# 创建一个新的字典用于存储修改后的DataFrame
modified_excel_data = {}

# 遍历字典，每个键值对对应一个子表格的名称和数据
for sheet_name, df in excel_data.items():
    # 创建一个新的DataFrame用于存储修改后的数据
    modified_df = df.copy()

    # 计算每行数据的 Gini Index
    gini_values = []
    for index, row in df.iterrows():
        data = row.tolist()
        data = data[1:]
        n = len(data)
        gini_sum = 0
        u = sum(data) / n
        for i in range(n):
            for j in range(n):
                gini_sum = gini_sum + abs(data[i] - data[j])
        gini_index = 1 - (1 / (2 * u * n * n)) * gini_sum
        gini_values.append(gini_index)

    # 将 Gini Index 添加到 DataFrame 的右侧
    modified_df['Gini Index'] = gini_values
    print(np.mean(gini_values))
    # 将修改后的DataFrame存储到新的字典中
    modified_excel_data[sheet_name] = modified_df

# 将修改后的数据存储到新的 Excel 文件中
with pd.ExcelWriter('gini5.xlsx') as writer:
    for sheet_name, modified_df in modified_excel_data.items():
        modified_df.to_excel(writer, sheet_name=sheet_name, index=False)

