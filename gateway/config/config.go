package config

import "time"

type DBConfig struct {
	User string
	Pass string
	Host string
	Port string
	Name string
}

type KafkaConfig struct {
	Brokers       []string
	ProducerTopic string
	ConsumerTopic string
	GroupID       string
}

type ServerConfig struct {
	Addr      string
	UploadDir string
}

type JWTConfig struct {
	Secret         string
	ExpireDuration time.Duration
}

var Settings = struct {
	DB     DBConfig
	Kafka  KafkaConfig
	Server ServerConfig
	JWT    JWTConfig
}{
	DB: DBConfig{
		User: "root",
		Pass: "",
		Host: "127.0.0.1",
		Port: "3306",
		Name: "",
	},
	Kafka: KafkaConfig{
		Brokers:       []string{"127.0.0.1:9092"},
		ProducerTopic: "task_input",
		ConsumerTopic: "task_result",
		GroupID:       "go-gateway",
	},
	Server: ServerConfig{
		Addr:      "localhost:8080",
		UploadDir: "/tmp/uploads",
	},
	JWT: JWTConfig{
		Secret:         "testSecretKey012345",
		ExpireDuration: 24 * time.Hour,
	},
}
