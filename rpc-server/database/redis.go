package database

import (
	"encoding/json"
	"sort"

	"github.com/ernst12/Backend_Server-TikTok_Tech_Immersion-Assignment/rpc-server/kitex_gen/rpc"
	"github.com/go-redis/redis/v7"
)

type redisDatabase struct {
	client *redis.Client
}

func createRedisDatabase() (Database, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis-client:6379",
		Password: "", // no password
		DB:       0,
	})
	_, err := client.Ping().Result() // makes sure db is connected
	if err != nil {
		return nil, &CreateDatabaseError{}
	}
	return &redisDatabase{client: client}, nil
}

func (r *redisDatabase) Append(key string, value *rpc.Message) error {
	oldValue, err := r.client.Get(key).Bytes()
	if err == redis.Nil {
		oldValue = nil
	} else if err != nil {
		return err
	}

	var messsageArr []rpc.Message

	if oldValue != nil && len(oldValue) != 0 {
		// append to existing chats instead
		//tempArr := JsonType{}
		var tempArr []rpc.Message

		if jsonErr := json.Unmarshal(oldValue, &tempArr); jsonErr != nil {
			return jsonErr
		}

		messsageArr = append(tempArr, *value)
	} else {
		messsageArr = append(messsageArr, *value)
	}

	// sort sendTime in ascending order
	sort.Slice(messsageArr, func(i, j int) bool {
		return messsageArr[i].SendTime < messsageArr[j].SendTime
	})

	serialized, jsonErr := json.Marshal(messsageArr)
	if jsonErr != nil {
		return jsonErr
	}

	err = r.client.Set(key, serialized, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *redisDatabase) Get(key string) ([]*rpc.Message, error) {
	value, err := r.client.Get(key).Bytes()
	if err == redis.Nil {
		return nil, nil // key not found
	} else if err != nil {
		return nil, err
	}

	var messageArr []*rpc.Message
	if jsonErr := json.Unmarshal(value, &messageArr); jsonErr != nil {
		return nil, jsonErr
	}

	return messageArr, nil
}

func (r *redisDatabase) Delete(key string) error {
	_, err := r.client.Del(key).Result()
	if err != nil {
		return err
	}
	return nil
}
