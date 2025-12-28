import os
import warnings
from collections import defaultdict
import pytorch_lightning as pl
import pandas as pd
import numpy as np
import torch, random
import torch.optim as optim

from model import iTransformer, LSTM, MixFormer
from config import Config
from cfg.config import settings

def outputing(inputs):#inputs为list,含有pass_data, future_data两项,pass_data形状为torch.Size([1, 72, 13]),future_data的形状为torch.Size([1, 24, 12],其中第一维度为batch_size,可以随意调整,第二维为时间,第三维为特征
    configs = Config()
    model = MixFormer.Model(configs)

    path = settings.model_path
    state_dict = torch.load(path)
    model.load_state_dict(state_dict)
    model.eval()

    pass_data, future_data = inputs

    p_data_mean = torch.mean(pass_data, dim=0)
    p_data_std = torch.std(pass_data, dim=0)
    pass_data = (pass_data - p_data_mean) / (p_data_std + 1e-5)

    f_data_mean = torch.mean(future_data, dim=0)
    f_data_std = torch.std(future_data, dim=0)
    future_data = (future_data - f_data_mean) / (f_data_std + 1e-5)

    power_data = pass_data[:,-24:,-1:]
    data_mean = torch.mean(power_data, dim=0)
    data_std = torch.std(power_data, dim=0)

    outputs = model(pass_data, future_data)

    outputs = outputs * (data_std + 1e-5) + data_mean
    # max_val = torch.max(outputs)  
    # threshold = max_val * 0.15
    # outputs[outputs < threshold] = 0


    return outputs#形状为torch.Size([1, 24, 1]