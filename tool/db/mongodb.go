package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB driver using go.mongodb.org/mongo-driver/mongo
type MongoDB struct {
	info     *Connection
	instance *mongo.Client
}

// CheckProject check project exists
func (db *MongoDB) CheckProject(project string) bool {
	return false
}

// Save insert or update document
func (db *MongoDB) Save(project, kind, key string, data interface{}) error {
	return nil
}

// Find document by key
func (db *MongoDB) Find(project, kind, key string) ([]byte, error) {
	return nil, nil
}

// FindAll document
func (db *MongoDB) FindAll(callback func(project, kind, key string, data []byte)) error {
	return nil
}

// Delete document
func (db *MongoDB) Delete(project, kind, key string) (err error) {
	return
}

// Connect db
func (db *MongoDB) Connect(conn *Connection, args ...interface{}) (err error) {
	db.info = conn

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	url := db.info.Host
	if db.info.Port != "" {
		url = fmt.Sprintf("%s:%s", db.info.Host, db.info.Port)
	}
	opts := options.Client().ApplyURI(url)

	if db.info.User != "" && db.info.Pass != "" {
		opts = opts.SetAuth(options.Credential{
			Username: db.info.User,
			Password: db.info.Pass,
		})
	}

	db.instance, err = mongo.Connect(ctx, opts)

	return
}

// Close db
func (db *MongoDB) Close() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return db.instance.Disconnect(ctx)
}
