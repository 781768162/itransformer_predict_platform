from dataclasses import dataclass
from typing import List


@dataclass
class Settings:
	task_input_topic: str = "task_input"
	task_result_topic: str = "task_result"
	kafka_brokers: List[str] = ("kafka:9092",)
	kafka_group_id: str = "python-algo"
	model_path: str = "save_model/best_Mix_24"


settings = Settings()
