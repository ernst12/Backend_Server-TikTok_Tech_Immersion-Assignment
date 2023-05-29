package database

type Database interface {
	Set(key string, value *Chat) (error)
	Get(key string) (*Chat, error)
	Delete(key string) (error)
}

type Chat struct {
	Chat string
	Text string
	Sender string
	Send_time int64
}

func Factory(databaseName string) (Database, error) {
	switch databaseName {
	case "redis":
		return createRedisDatabase()
	default:
		return nil, &NotImplementedDatabaseError{databaseName}
	}
}
