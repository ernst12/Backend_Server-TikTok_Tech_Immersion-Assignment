package database

import (
	"encoding/json"

	"github.com/go-redis/redis/v7"
)

type RedisClient struct {
	cli *redis.Client
}

func (c *RedisClient) InitClient(address, password string) error {
	r := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password, // no password set
		DB:       0,        // use default DB
	})

	// test connection
	if err := r.Ping().Err(); err != nil {
		return err
	}

	c.cli = r
	return nil
}

/*func createRedisDatabase() (Database, error) {
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
}*/

func (c *RedisClient) SaveMessage(roomID string, message *Message) error {
	// Marshal the Go struct into JSON bytes
	text, err := json.Marshal(message)
	if err != nil {
		return err
	}

	member := &redis.Z{
		Score:  float64(message.Timestamp), // The sort key
		Member: text,                       // Data
	}

	_, err = c.cli.ZAdd(roomID, member).Result()
	if err != nil {
		return err
	}

	return nil
}

/*func (r *RedisClient) Append(key string, value *rpc.Message) error {
	oldValue, err := r.cli.Get(key).Bytes()
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

	err = r.cli.Set(key, serialized, 0).Err()
	if err != nil {
		return err
	}

	return nil
}*/

func (c *RedisClient) GetMessagesByRoomID(roomID string, start, end int64, reverse bool) ([]*Message, error) {
	var (
		rawMessages []string
		messages    []*Message
		err         error
	)

	if reverse {
		// Desc order with time -> first message is the latest message
		rawMessages, err = c.cli.ZRevRange(roomID, start, end).Result()
		if err != nil {
			return nil, err
		}
	} else {
		// Asc order with time -> first message is the earliest message
		rawMessages, err = c.cli.ZRange(roomID, start, end).Result()
		if err != nil {
			return nil, err
		}
	}

	for _, msg := range rawMessages {
		temp := &Message{}
		err := json.Unmarshal([]byte(msg), temp)
		if err != nil {
			return nil, err
		}
		messages = append(messages, temp)
	}

	return messages, nil
}

/*func (r *RedisClient) Get(key string) ([]*rpc.Message, error) {
	value, err := r.cli.Get(key).Bytes()
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
}*/

func (r *RedisClient) Delete(key string) error {
	_, err := r.cli.Del(key).Result()
	if err != nil {
		return err
	}
	return nil
}
