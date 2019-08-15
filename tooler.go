// tooler
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/boltdb/bolt"
)

func init() {
	log.Println("init tooler!")
}

/**for boltdb
 */
type SetKeyValue func(k, v []byte)
type GetKeyValue func() ([]byte, []byte)

func boltEnumKeyValue(fileDb string, bucket string, kv SetKeyValue) {
	db, err := bolt.Open(fileDb, 0600, nil)
	if err != nil {
		log.Println("[Err]", err.Error())
		return
	}
	defer db.Close()

	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return nil
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			kv(k, v)
		}

		return nil
	})
}

func boltBatchWriteKeyValue(fileDb string, bucket string, kv GetKeyValue) {
	db, err := bolt.Open(fileDb, 0600, nil)
	if err != nil {
		log.Println("[Err]", err.Error())
		return
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		for {
			k, v := kv()
			if k == nil || v == nil {
				break
			}

			err = b.Put(k, v)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func boltWriteKeyValue(fileDb string, bucket string, kv GetKeyValue) {
	db, err := bolt.Open(fileDb, 0600, nil)
	if err != nil {
		log.Println("[Err] boltWriteKeyValue", err.Error())
		return
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		k, v := kv()
		if k == nil || v == nil {
			return nil
		}

		err = b.Put(k, v)
		if err != nil {
			return err
		}
		return nil
	})
}

func boltDeleteBucket(fileDb string, bucket string) {
	db, err := bolt.Open(fileDb, 0600, nil)
	if err != nil {
		log.Println("[Err]", err.Error())
		return
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(bucket))
	})
}

/**for excel
 */

/**for key
 */
func PhoneProperty2Key(pp *PhoneProperty) string {
	return fmt.Sprintf("%s_%s_%s_%s", pp.productor, pp.isp, pp.province, pp.area)
}

func Key2PhoneProperty(key string) *PhoneProperty { //productor, isp, province, area
	p := strings.Split(key, "_")
	return &PhoneProperty{p[0], p[1], p[2], p[3]}
}

/**for httpClient
 */
type JsonParser func([]byte)

func HttpPost(client *http.Client, uri string, data string, parseJson JsonParser) bool {
	request, err := http.NewRequest("POST", uri, strings.NewReader(data))
	if err != nil {
		log.Println("[Err] new request.", err.Error())
		return false
	}
	request.Header.Set("Content-type", "application/json")
	request.Header.Set("charset", "utf-8")

	response, err := (*client).Do(request)
	if err != nil {
		log.Println("[Err] Do request.", err.Error())
		return false
	}
	if response.StatusCode != 200 {
		log.Println("[Err]", "reponse err.", response.StatusCode)
		return false
	}

	if parseJson != nil {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Println("[Err]", err.Error())
			return false
		}
		parseJson(body)
	}

	return true
}
