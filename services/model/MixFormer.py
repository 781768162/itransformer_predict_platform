import torch
import torch.nn as nn
import torch.nn.functional as F
from model import iTransformer,TimeMixer

class Model(nn.Module):

    def __init__(self, configs):
        super(Model, self).__init__()
        self.past_layer = iTransformer.Model(configs)
        configs.enc_in = 12
        configs.seq_len = 24
        # self.future_layer = iTransformer.Model(configs)
        self.future_layer = TimeMixer.Model(seq_len=24, enc_in=12, pred_len=24, 
                                            d_model=12, dropout=0.3, 
                                            e_layers=4, dec_in=1)
        # self.linear_1= nn.Linear(configs.pred_len * 2, configs.pred_len)
        self.linear_1= nn.Linear(2, 1)
    


    def forward(self, x_past, x_known):
        y_1 = self.past_layer(x_past)
        y_2 = self.future_layer(x_known)
        y = torch.cat((y_1, y_2), dim=2)
        # y_t = torch.transpose(y, 1, 2)
        dec_out =  F.relu(self.linear_1(y))
        # dec_out = torch.transpose(dec_out, 1, 2)
        # dec_out = y_2

        return dec_out