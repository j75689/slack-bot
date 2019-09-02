package db

import "errors"

// Driver name
type Driver string

const (
	mongodb Driver = "MongoDB"
	bolt    Driver = "Bolt"
)

// Connection struct
type Connection struct {
	DBName string
	Host   string
	Port   string
	User   string
	Pass   string
}

// Storage definition abstract method
type Storage interface {
	CheckProject(project string) bool
	Save(project, kind, key string, data interface{}) error
	Find(project, kind, key string) ([]byte, error)
	FindAll(callback func(project, kind, key string, data []byte)) error
	Delete(project, kind, key string) error
	Connect(conn *Connection, args ...interface{}) error
	Close() error
}

var supported = map[Driver]Storage{
	mongodb: new(MongoDB),
	bolt:    new(BoltDB),
}

// New db driver
func New(driver Driver, conn *Connection, args ...interface{}) (Storage, error) {
	if constructor := supported[driver]; constructor != nil {
		return constructor, constructor.Connect(conn, args...)
	}
	return nil, errors.New("not supported driver [" + string(driver) + "]")
}
