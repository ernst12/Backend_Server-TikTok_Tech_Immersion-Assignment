package database

import (
	"github.com/go-redis/redis/v7"
	"encoding/json"
)

type redisDatabase struct {
	client *redis.Client
}

func createRedisDatabase() (Database, error) {
	client := redis.NewClient(&redis.Options{
		Addr: "redis-client:6379",
		Password: "", // no password
		DB: 0,
	})
	_, err := client.Ping().Result() // makes sure db is connected
	if err != nil {
		return nil, &CreateDatabaseError{}
	}
	return &redisDatabase{client: client}, nil
}

func (r *redisDatabase) Set(key string, value *Chat) (error) {
	serialized, jsonErr := json.Marshal(value)
    if jsonErr != nil {
       return jsonErr
    }
	
	_, err := r.client.Set(key, serialized, 0).Result()
	if err != nil {
		return generateError("set", err)
	}
	
	return nil
}

func (r *redisDatabase) Get(key string) (*Chat, error) {
	value, err := r.client.Get(key).Bytes()
	if err != nil {
		return &Chat{}, generateError("get", err)
	}

	var chat *Chat
	if jsonErr := json.Unmarshal(value, &chat); jsonErr != nil {
		return &Chat{}, jsonErr
	}

	return chat, nil
}

func (r *redisDatabase) Delete(key string) (error) {
	_, err := r.client.Del(key).Result()
	if err != nil {
		return generateError("delete", err)
	}
	return nil
}

func generateError(operation string, err error) (error) {
	if err == redis.Nil {
		return &OperationError{operation}
	}
	return &DownError{}
}
