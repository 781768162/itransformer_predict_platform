class Config:
    def __init__(self):
        # Task configuration
        self.task_name = 'long_term_forecast'
        self.is_training = True
        self.model = 'iTransformer'

        # Data loader configuration
        self.data = 'data'
        self.root_path = 'data/thermal/data.csv'
        self.data_path = 'data.csv'

        # Forecasting task configuration
        self.seq_len = 72
        self.label_len = 1
        self.pred_len = 24
        self.seasonal_patterns = 'Monthly'
        self.inverse = False

        # Model define configuration
        self.expand = 2
        self.d_conv = 4
        self.top_k = 5
        self.num_kernels = 6
        self.enc_in = 13
        self.dec_in = 1
        # self.c_out = 26
        self.d_model = 12
        self.n_heads = 2
        self.e_layers = 4
        self.d_layers = 1
        self.d_ff = 1024
        self.moving_avg = 25
        self.factor = 1
        self.distil = True
        self.dropout = 0.3
        self.embed = 'timeF'
        self.activation = 'gelu'
        self.output_attention = False
        self.channel_independence = 1
        self.decomp_method = 'moving_avg'
        self.use_norm = 1
        self.down_sampling_layers = 0
        self.down_sampling_window = 1
        self.down_sampling_method = None
        self.seg_len = 48
        self.freq = 'h'

        # Optimization configuration
        self.num_workers = 10
        self.itr = 1
        self.train_epochs = 1000
        self.batch_size = 1
        self.patience = 50
        self.learning_rate = 1e-3
        self.des = 'test'
        self.loss = 'MSE'
        self.lradj = 'type1'
        self.use_amp = False

        # GPU configuration
        self.use_gpu = True
        self.gpu = 0
        self.use_multi_gpu = False
        self.devices = '0,1,2,3'

        # De-stationary projector configuration
        self.p_hidden_dims = [128, 128]
        self.p_hidden_layers = 2
