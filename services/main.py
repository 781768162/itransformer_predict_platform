import logging

import mq.mq as mq


def main():
	logging.basicConfig(
		level=logging.INFO,
		format="%(asctime)s %(levelname)s %(message)s",
		filename="log/service.log",   # 设置日志文件路径
        filemode="a",                 # 追加写入
	)
	
	mq.run_forever()
	
if __name__ == "__main__":
	main()