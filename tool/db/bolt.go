package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"go.etcd.io/bbolt"
)

// BoltDB driver using go.etcd.io/bbolt
type BoltDB struct {
	info     *Connection
	instance *bbolt.DB
}

// CheckTable check table exists
func (db *BoltDB) CheckTable(table string) bool {

	err := db.instance.Batch(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(table))
		if b == nil {
			return errors.New("not found")
		}
		return nil
	})
	return err == nil
}

// Save insert or update document
func (db *BoltDB) Save(table string, key string, data interface{}) error {
	err := db.instance.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(table))
		if err != nil {
			return err
		}

		if byteData, err := json.Marshal(data); err == nil {
			err = b.Put([]byte(key), byteData)
			if err != nil {
				return err
			}
		}

		return nil
	})
	return err
}

// Find document by key
func (db *BoltDB) Find(table string, key string) (data []byte, err error) {

	err = db.instance.Batch(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(table))
		if err != nil {
			return err
		}
		data = b.Get([]byte(key))
		if data == nil {
			err = fmt.Errorf("data [%v] not found", key)
		}
		return err
	})

	return
}

// FindAll document in table
func (db *BoltDB) FindAll(table string, callback func(key string, data []byte)) (err error) {
	err = db.instance.Batch(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(table))
		if err != nil {
			return err
		}
		err = b.ForEach(func(k, v []byte) (err error) {
			callback(string(k), v)
			return
		})

		return err
	})
	return
}

// Delete document
func (db *BoltDB) Delete(table string, key string) (err error) {
	err = db.instance.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(table))
		if err != nil {
			return err
		}
		err = b.Delete([]byte(key))
		if err != nil {
			return err
		}
		return nil
	})
	return
}

// Connect db
func (db *BoltDB) Connect(conn *Connection, args ...interface{}) (err error) {
	db.info = conn
	// check directory
	path := conn.Host
	if strings.LastIndex(path, "/") > -1 {
		path = path[0:strings.LastIndex(path, "/")]
		os.MkdirAll(path, 0755)
	}

	db.instance, err = bbolt.Open(conn.Host, 0644, nil)
	return
}

// Close db
func (db *BoltDB) Close() (err error) {
	return db.instance.Close()
}
