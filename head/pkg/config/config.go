package config

import "os"

type Config struct{}

func (c Config) GetKey(key string) string {
	return os.Getenv(key)
}

func (c Config) GetRedisHost() string {
	return c.GetKey("REDIS_HOST")
}

func (c Config) GetKafkaHost() string {
	return c.GetKey("KAFKA_HOST")
}

func (c Config) GetClickhouseHost() string {
	return c.GetKey("CLICK_HOST")
}

func (c Config) GetKafkaTopic() string {
	return c.GetKey("KAFKA_TOPIC")
}
