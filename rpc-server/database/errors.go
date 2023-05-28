package database

type OperationError struct {
	operation string
}

func (err *OperationError) Error() string {
	return "Could not perform the " + err.operation + " operation"
}

// DownError when its not a redis.Nil response, in this case the database is down
type DownError struct{}

func (downErr *DownError) Error() string {
	return "Database is down"
}

// CreateDatabaseError when can't perform set on database
type CreateDatabaseError struct{}

func (err *CreateDatabaseError) Error() string {
	return "Could not create Database"
}

// NotImplementedDatabaseError when user tries to create a not implemented database
type NotImplementedDatabaseError struct {
	database string
}

func (err *NotImplementedDatabaseError) Error() string {
	return err.database + " not implemented"
}
