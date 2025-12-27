import torch
import torch.nn as nn
import torch.nn.functional as F


class TimeMixLayer(nn.Module):
    """时间混合层：捕获时间维度上的依赖关系"""
    def __init__(self, d_model, seq_len, dropout=0.1):
        super(TimeMixLayer, self).__init__()
        self.d_model = d_model
        self.seq_len = seq_len
        
        # 时间注意力权重（可学习）
        self.time_weights = nn.Parameter(torch.randn(seq_len, seq_len))
        # 前馈网络
        self.ffn = nn.Sequential(
            nn.Linear(d_model, d_model * 4),
            nn.GELU(),
            nn.Dropout(dropout),
            nn.Linear(d_model * 4, d_model)
        )
        self.norm1 = nn.LayerNorm(d_model)
        self.norm2 = nn.LayerNorm(d_model)
        self.dropout = nn.Dropout(dropout)

    def forward(self, x):
        # x: [batch_size, seq_len, d_model]
        batch_size, seq_len, d_model = x.shape
        
        # 时间混合：通过权重矩阵融合不同时间步信息
        time_attn = F.softmax(self.time_weights, dim=-1)  # [seq_len, seq_len]
        x_time = torch.matmul(time_attn, x)  # [batch_size, seq_len, d_model]
        
        # 残差连接 + 归一化
        x = self.norm1(x + self.dropout(x_time))
        
        # 前馈网络
        x_ffn = self.ffn(x)
        x = self.norm2(x + self.dropout(x_ffn))
        
        return x


class FeatureMixLayer(nn.Module):
    """特征混合层：捕获特征维度上的依赖关系"""
    def __init__(self, d_model, dropout=0.1):
        super(FeatureMixLayer, self).__init__()
        self.d_model = d_model
        
        # 特征注意力权重（可学习）
        self.feature_weights = nn.Parameter(torch.randn(d_model, d_model))
        # 前馈网络
        self.ffn = nn.Sequential(
            nn.Linear(d_model, d_model * 4),
            nn.GELU(),
            nn.Dropout(dropout),
            nn.Linear(d_model * 4, d_model)
        )
        self.norm1 = nn.LayerNorm(d_model)
        self.norm2 = nn.LayerNorm(d_model)
        self.dropout = nn.Dropout(dropout)

    def forward(self, x):
        # x: [batch_size, seq_len, d_model]
        batch_size, seq_len, d_model = x.shape
        
        # 特征混合：通过权重矩阵融合不同特征信息
        x_trans = x.transpose(1, 2)  # [batch_size, d_model, seq_len]
        feature_attn = F.softmax(self.feature_weights, dim=-1)  # [d_model, d_model]
        x_feat = torch.matmul(feature_attn, x_trans).transpose(1, 2)  # [batch_size, seq_len, d_model]
        
        # 残差连接 + 归一化
        x = self.norm1(x + self.dropout(x_feat))
        
        # 前馈网络
        x_ffn = self.ffn(x)
        x = self.norm2(x + self.dropout(x_ffn))
        
        return x


class DataEmbedding(nn.Module):
    """数据嵌入层：将原始时序数据映射到高维空间"""
    def __init__(self, c_in, d_model, seq_len, dropout=0.1):
        super(DataEmbedding, self).__init__()
        self.value_embedding = nn.Linear(c_in, d_model)  # 特征嵌入
        self.position_embedding = nn.Parameter(torch.randn(1, seq_len, d_model))  # 位置嵌入
        self.dropout = nn.Dropout(dropout)

    def forward(self, x):
        # x: [batch_size, seq_len, c_in]
        x = self.value_embedding(x)  # [batch_size, seq_len, d_model]
        x += self.position_embedding  # 加入位置信息
        return self.dropout(x)


class Model(nn.Module):
    """TimeMixer模型主类"""
    def __init__(self, seq_len, enc_in, pred_len, d_model, dropout, e_layers, dec_in):
        super(Model, self).__init__()
        self.seq_len = seq_len
        self.pred_len = pred_len
        self.d_model = d_model
        
        # 嵌入层
        self.embedding = DataEmbedding(
            c_in=enc_in,
            d_model=d_model,
            seq_len=seq_len,
            dropout=dropout
        )
        
        # 时间混合层堆叠
        self.time_mix_layers = nn.ModuleList([
            TimeMixLayer(
                d_model=d_model,
                seq_len=seq_len,
                dropout=dropout
            ) for _ in range(e_layers)
        ])
        
        # 特征混合层堆叠
        self.feature_mix_layers = nn.ModuleList([
            FeatureMixLayer(
                d_model=d_model,
                dropout=dropout
            ) for _ in range(e_layers)
        ])
        
        # 输出投影层（将隐藏状态映射到预测长度）
        self.projection = nn.Linear(seq_len, pred_len)
        self.final_projection = nn.Linear(d_model, dec_in)

    def forward(self, x_enc):
        # x_enc: [batch_size, seq_len, enc_in]
        
        # 数据归一化（参考iTransformer的非平稳处理）
        means = x_enc.mean(1, keepdim=True).detach()
        x_enc = x_enc - means
        stdev = torch.sqrt(
            torch.var(x_enc, dim=1, keepdim=True, unbiased=False) + 1e-5)
        x_enc /= stdev
        
        # 嵌入层
        enc_out = self.embedding(x_enc)  # [batch_size, seq_len, d_model]
        
        # 时间混合
        for time_layer in self.time_mix_layers:
            enc_out = time_layer(enc_out)
        
        # 特征混合
        for feat_layer in self.feature_mix_layers:
            enc_out = feat_layer(enc_out)
        
        # 投影到预测长度
        dec_out = self.projection(enc_out.transpose(1, 2)).transpose(1, 2)  # [batch_size, pred_len, d_model]
        
        # 还原归一化
        dec_out = dec_out * \
                  (stdev[:, 0, :].unsqueeze(1).repeat(1, self.pred_len, 1))
        dec_out = dec_out + \
                  (means[:, 0, :].unsqueeze(1).repeat(1, self.pred_len, 1))
        
        # 映射到输出特征维度
        dec_out = self.final_projection(dec_out)  # [batch_size, pred_len, dec_in]
        
        return dec_out