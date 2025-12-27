import torch
import torch.nn as nn


class Model(nn.Module):
    """
    LSTM模型，用于时间序列预测任务
    参考iTransformer的结构设计，包含数据归一化、LSTM编码和解码投影等模块
    """

    def __init__(self, configs):
        super(Model, self).__init__()
        self.task_name = configs.task_name
        self.seq_len = configs.seq_len  # 输入序列长度
        self.pred_len = configs.pred_len  # 预测序列长度
        self.input_dim = configs.enc_in  # 输入特征维度
        self.hidden_dim = configs.d_model  # 隐藏层维度
        self.num_layers = configs.e_layers  # LSTM层数
        self.dropout = configs.dropout  # dropout率

        # LSTM编码器
        self.lstm = nn.LSTM(
            input_size=self.input_dim,
            hidden_size=self.hidden_dim,
            num_layers=self.num_layers,
            dropout=self.dropout,
            batch_first=True  # 批量维度在前
        )

        # 输出投影层，将LSTM输出映射到预测长度
        self.projection = nn.Linear(self.hidden_dim, self.pred_len)
        # 特征维度映射（与iTransformer保持一致）
        self.projection2 = nn.Linear(configs.enc_in, configs.dec_in)

    def forward(self, x_enc):
        # 数据归一化（参考Non-stationary Transformer的处理方式）
        means = x_enc.mean(1, keepdim=True).detach()
        x_enc = x_enc - means
        stdev = torch.sqrt(torch.var(x_enc, dim=1, keepdim=True, unbiased=False) + 1e-5)
        x_enc /= stdev

        batch_size, seq_len, input_dim = x_enc.shape

        # LSTM前向传播
        # 输出形状: (batch_size, seq_len, hidden_dim)
        lstm_out, _ = self.lstm(x_enc)

        # 取最后一个时间步的输出进行预测
        # 也可以使用所有时间步输出的均值，根据任务需求调整
        last_out = lstm_out[:, -1, :]  # (batch_size, hidden_dim)

        # 投影到预测长度
        dec_out = self.projection(last_out)  # (batch_size, pred_len)
        # 调整维度以匹配输入格式 (batch_size, pred_len, input_dim)
        dec_out = dec_out.unsqueeze(-1).repeat(1, 1, input_dim)

        # 反归一化
        dec_out = dec_out * stdev.repeat(1, self.pred_len, 1)
        dec_out = dec_out + means.repeat(1, self.pred_len, 1)

        # 特征维度映射
        dec_out = self.projection2(dec_out)

        return dec_out