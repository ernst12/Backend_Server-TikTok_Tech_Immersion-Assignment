package database

import "github.com/go-redis/redis/v7"

type redisDatabase struct {
	client *redis.Client
}

func createRedisDatabase() (Database, error) {
	client := redis.NewClient(&redis.Options{
		Addr: "redis-client",
		Password: "", // no password
		DB: 0,
	})
	_, err := client.Ping().Result() // makes sure db is connected
	if err != nil {
		return nil, &CreateDatabaseError{}
	}
	return &redisDatabase{client: client}, nil
}

func (r *redisDatabase) Set(key string, value string) (string, error) {
	_, err := r.client.Set(key, value, 0).Result()
	if err != nil {
		return generateError("set", err)
	}
	return key, nil
}

func (r *redisDatabase) Get(key string) (string, error) {
	value, err := r.client.Get(key).Result()
	if err != nil {
		return generateError("get", err)
	}
	return value, nil
}

func (r *redisDatabase) Delete(key string) (string, error) {
	_, err := r.client.Del(key).Result()
	if err != nil {
		return generateError("delete", err)
	}
	return key, nil
}

func generateError(operation string, err error) (string, error) {
	if err == redis.Nil {
		return "", &OperationError{operation}
	}
	return "", &DownError{}
}
