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

// CheckProject check project exists
func (db *BoltDB) CheckProject(project string) bool {

	err := db.instance.Batch(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(project))
		if b == nil {
			return errors.New("not found")
		}
		return nil
	})
	return err == nil
}

// Save insert or update document
func (db *BoltDB) Save(project, kind, key string, data interface{}) error {
	err := db.instance.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(project))
		if err != nil {
			return err
		}
		var documents map[string]interface{}
		err = json.Unmarshal(b.Get([]byte(kind)), &documents)
		if err != nil {
			documents = make(map[string]interface{})
		}
		documents[kind] = data

		if byteData, err := json.Marshal(documents); err == nil {
			err = b.Put([]byte(kind), byteData)
			if err != nil {
				return err
			}
		}

		return nil
	})
	return err
}

// Find document by key
func (db *BoltDB) Find(project, kind, key string) (data []byte, err error) {

	err = db.instance.Batch(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(project))
		if err != nil {
			return err
		}
		var documents map[string]interface{}
		err = json.Unmarshal(b.Get([]byte(kind)), &documents)
		if err != nil {
			return err
		}

		if documents[key] != nil {
			data, err = json.Marshal(documents[key])
		} else {
			err = fmt.Errorf("data [%s] not exists", key)
		}

		return err
	})

	return
}

// FindAll document in project
func (db *BoltDB) FindAll(project, kind string, callback func(key string, data []byte)) (err error) {
	err = db.instance.Batch(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(project))
		if err != nil {
			return err
		}

		var documents map[string]interface{}
		err = json.Unmarshal(b.Get([]byte(kind)), &documents)
		if err != nil {
			return err
		}

		for k, v := range documents {
			if byteData, err := json.Marshal(v); err == nil {
				callback(k, byteData)
			}
		}

		return nil
	})
	return
}

// Delete document
func (db *BoltDB) Delete(project, kind, key string) (err error) {
	err = db.instance.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(project))
		if err != nil {
			return err
		}

		var documents map[string]interface{}
		err = json.Unmarshal(b.Get([]byte(kind)), &documents)
		if err != nil {
			return err
		}

		// remove doc
		delete(documents, key)

		// save
		if byteData, err := json.Marshal(documents); err == nil {
			err = b.Put([]byte(kind), byteData)
			if err != nil {
				return err
			}
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
