import torch
import worker


class DummyModel(torch.nn.Module):
    def __init__(self, configs):
        super().__init__()

    def load_state_dict(self, state):
        # Accept any state dict silently
        return

    def eval(self):
        return self

    def forward(self, pass_data, future_data):
        # Return ones with expected shape [batch, 24, 1]
        batch = pass_data.shape[0]
        device = pass_data.device
        return torch.ones((batch, 24, 1), device=device)


def run_smoke_test():
    # Monkeypatch the heavy model and torch.load to avoid needing real weights
    # worker.MixFormer.Model = DummyModel  # type: ignore
    # worker.torch.load = lambda path: {}  # type: ignore

    pass_tensor = torch.randn(1, 72, 13)
    future_tensor = torch.randn(1, 24, 12)

    outputs = worker.outputing([pass_tensor, future_tensor])
    print("output shape:", outputs.shape)
    print("output values:", outputs.squeeze(-1).tolist())


if __name__ == "__main__":
    run_smoke_test()
