import os
import warnings
from collections import defaultdict
import pytorch_lightning as pl
import pandas as pd
import numpy as np
import torch,random
import torch.optim as optim

from model import iTransformer,LSTM,MixFormer
from config import settings

def outputing(inputs):#inputs为list,含有pass_data, future_data两项,pass_data形状为torch.Size([1, 72, 13]),future_data的形状为torch.Size([1, 24, 12],其中第一维度为batch_size,可以随意调整,第二维为时间,第三维为特征
    configs = Config()
    model = MixFormer.Model(configs)

    path = settings.model_path
    state_dict = torch.load(path)
    model.load_state_dict(state_dict)
    model.eval()

    pass_data, future_data = inputs
    outputs = model(pass_data, future_data)

    return outputs#形状为torch.Size([1, 24, 1]