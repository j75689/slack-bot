package db

import (
	"fmt"
	"os"
	"strings"

	"github.com/j75689/slack-bot/kind"
	"gopkg.in/yaml.v2"

	"go.etcd.io/bbolt"
)

// BoltDB driver using go.etcd.io/bbolt
type BoltDB struct {
	info     *Connection
	instance *bbolt.DB
}

// CheckProject check project exists
func (db *BoltDB) CheckProject(project string) bool {
	data, err := db.Find(project, kind.Project, project)
	return err == nil && data != nil
}

// Save insert or update document
func (db *BoltDB) Save(project, kind, key string, data interface{}) error {
	err := db.instance.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(project))
		if err != nil {
			return err
		}

		var documents = make(map[string]interface{})
		err = yaml.Unmarshal(b.Get([]byte(kind)), &documents)

		documents[key] = data

		byteData, err := yaml.Marshal(documents)
		if err != nil {
			return err
		}

		err = b.Put([]byte(kind), byteData)
		if err != nil {
			return err
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
		err = yaml.Unmarshal(b.Get([]byte(kind)), &documents)
		if err != nil {
			return err
		}

		if documents[key] != nil {
			data, err = yaml.Marshal(documents[key])
		} else {
			err = fmt.Errorf("data [%s] not exists", key)
		}

		return err
	})

	return
}

// FindAll document
func (db *BoltDB) FindAll(callback func(project, kind, key string, data []byte)) (err error) {

	type Temp struct {
		Project string
		Kind    string
		Key     string
		Data    []byte
	}
	var tempdata = []*Temp{}

	err = db.instance.Batch(func(tx *bbolt.Tx) error {
		// loop bucket
		err = tx.ForEach(func(project []byte, b *bbolt.Bucket) error {
			var documents map[string]interface{}
			// loop kind
			err = b.ForEach(func(kind, data []byte) error {
				documents = make(map[string]interface{})
				err = yaml.Unmarshal(b.Get(kind), &documents)
				if err != nil {
					return err
				}
				// loop document
				for k, v := range documents {
					if byteData, err := yaml.Marshal(v); err == nil {
						tempdata = append(tempdata, &Temp{
							Project: string(project),
							Kind:    string(kind),
							Key:     k,
							Data:    byteData,
						})
					}
				}
				return err
			})
			return err
		})
		return err
	})

	for _, temp := range tempdata {
		callback(temp.Project, temp.Kind, temp.Key, temp.Data)
	}

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
		err = yaml.Unmarshal(b.Get([]byte(kind)), &documents)
		if err != nil {
			return err
		}

		// remove doc
		delete(documents, key)

		// save
		if byteData, err := yaml.Marshal(documents); err == nil {
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
