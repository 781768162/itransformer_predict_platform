import json
import logging
import time

import torch
from kafka import KafkaConsumer, KafkaProducer

from config import settings
from worker import outputing


def make_consumer():
	return KafkaConsumer(
		settings.task_input_topic,
		bootstrap_servers=settings.kafka_brokers,
		group_id=settings.kafka_group_id,
		enable_auto_commit=False,
		value_deserializer=lambda v: json.loads(v.decode("utf-8")),
		auto_offset_reset="earliest",
	)


def make_producer():
	return KafkaProducer(
		bootstrap_servers=settings.kafka_brokers,
		value_serializer=lambda v: json.dumps(v).encode("utf-8"),
	)


def run_loop():
	consumer = make_consumer()
	producer = make_producer()
	# 生成失败消息并提交 offset
	try:
		for msg in consumer:
			
			fail_msg = {
				"task_id": 0,
                "status": "failed",
                "result": []
            }
			
			logging.info("consume task_input offset=%s partition=%s", msg.offset, msg.partition)
			payload = msg.value
			
			try:
				task_id = payload["task_id"]
				pass_data = payload["pass_data"]
				future_data = payload["future_data"]
				
			except Exception:
				logging.warning("invalid payload, skip")
				
				fail_msg["task_id"] = task_id if task_id else 0
				future = producer.send(settings.task_result_topic, value=fail_msg, key=str(task_id).encode())
				future.get(timeout=10)
				
				consumer.commit()
				continue

			# 输入已包含 batch 维度，[1,72,13] / [1,24,12]
			pass_tensor = torch.tensor(pass_data, dtype=torch.float32)
			future_tensor = torch.tensor(future_data, dtype=torch.float32)

			try:
				outputs = outputing([pass_tensor, future_tensor])
				outputs_list = outputs.squeeze(-1).squeeze(0).tolist()
				
			except Exception as e:
				logging.error("infer failed task_id=%s err=%s", task_id, e)
				logging.error(f"error: {str(e)}")
				
				fail_msg["task_id"] = task_id if task_id else 0
				future = producer.send(settings.task_result_topic, value=fail_msg, key=str(task_id).encode())
				future.get(timeout=10)
				
				consumer.commit()
				continue

			result_msg = {
				"task_id": task_id,
				"status": "success",
				"result": outputs_list,
			}

			future = producer.send(settings.task_result_topic, value=result_msg, key=str(task_id).encode())
			future.get(timeout=10)
			logging.info("produce task_result task_id=%s", task_id)
			consumer.commit()
			
	finally:
		consumer.close()
		producer.close()


def run_forever():
	while True:
		try:
			run_loop()
		except Exception:
			time.sleep(1)
			continue
