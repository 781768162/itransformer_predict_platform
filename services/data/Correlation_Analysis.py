import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
from matplotlib.colors import LinearSegmentedColormap

# 设置中文字体
plt.rcParams['font.sans-serif'] = ['WenQuanYi Zen Hei']
plt.rcParams['axes.unicode_minus'] = False

# 1. 读取数据文件
name = input("数据名：")
df = pd.read_csv(f'data\solar_data{name}_2014.csv')

# 2. 数据初步探索
print("=== 数据基本信息 ===")
print(f"数据形状: {df.shape}")
print(f"列名列表: {list(df.columns)}")
print("\n=== 数据前5行 ===")
print(df.head())
print("\n=== 数据类型 ===")
print(df.dtypes)
print("\n=== 缺失值统计 ===")
print(df.isnull().sum())

# 3. 数据预处理：选择数值型列进行相关性分析，处理缺失值
# 筛选数值型列
numeric_cols = df.select_dtypes(include=[np.number]).columns.tolist()
print(f"\n=== 用于相关性分析的数值型列 ===")
print(numeric_cols)

# 处理缺失值（使用均值填充，避免删除过多数据）
df_numeric = df[numeric_cols].fillna(df[numeric_cols].mean())

# 4. 计算相关系数矩阵（Pearson相关系数）
corr_matrix = df_numeric.corr(method='pearson')
print(f"\n=== 相关系数矩阵形状 ===")
print(f"矩阵维度: {corr_matrix.shape}")
print("\n=== 相关系数矩阵（前5行5列） ===")
print(corr_matrix.iloc[:5, :5])

# 5. 绘制相关性热力图
# 设置画布大小
fig, ax = plt.subplots(figsize=(12, 10))

# 定义自定义颜色映射（蓝-白-红）
colors = ['#2c7fb8', '#e0f3f8', '#ffffbf', '#fee090', '#fc8d59', '#d73027']
n_bins = 100
cmap = LinearSegmentedColormap.from_list('custom_cmap', colors, N=n_bins)

# 绘制热力图
im = ax.imshow(corr_matrix.values, cmap=cmap, aspect='auto', vmin=-1, vmax=1)

# 设置坐标轴标签
ax.set_xticks(range(len(corr_matrix.columns)))
ax.set_yticks(range(len(corr_matrix.columns)))
ax.set_xticklabels(corr_matrix.columns, rotation=45, ha='right', fontsize=10)
ax.set_yticklabels(corr_matrix.columns, fontsize=10)

# 添加数值标注（保留2位小数）
for i in range(len(corr_matrix.columns)):
    for j in range(len(corr_matrix.columns)):
        text = ax.text(j, i, f'{corr_matrix.iloc[i, j]:.2f}',
                       ha="center", va="center", color="black" if abs(corr_matrix.iloc[i, j]) < 0.7 else "white",
                       fontsize=8)

# 添加颜色条
cbar = plt.colorbar(im, ax=ax, shrink=0.8)
cbar.set_label('Pearson 相关系数', rotation=270, labelpad=20, fontsize=12)

# 设置标题
ax.set_title('太阳能数据（2014）变量相关性热力图', fontsize=14, fontweight='bold', pad=20)

# 调整布局
plt.tight_layout()

# 保存图片
plt.savefig(f'data/heatmap/data{name}.png', dpi=300, bbox_inches='tight')
plt.close()

print("\n=== 相关性分析完成 ===")
print(f"1. 相关系数矩阵已计算完成（{corr_matrix.shape[0]} x {corr_matrix.shape[1]}）")
print(f"2. 热力图已保存至: heatmap")
print(f"3. 分析变量数量: {len(numeric_cols)} 个数值型变量")