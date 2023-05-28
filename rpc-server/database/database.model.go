package database

type Database interface {
	Set(key string, value string) (string, error)
	Get(key string) (string, error)
	Delete(key string) (string, error)
}

func Factory(databaseName string) (Database, error) {
	switch databaseName {
	case "redis":
		return createRedisDatabase()
	default:
		return nil, &NotImplementedDatabaseError{databaseName}
	}
}
