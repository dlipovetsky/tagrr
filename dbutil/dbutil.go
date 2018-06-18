package dbutil

import (
	"bytes"

	bolt "github.com/coreos/bbolt"
)

func GetKeys(b *bolt.Bucket, keys []string) map[string]string {
	result := make(map[string]string)
	for _, k := range keys {
		if v := b.Get([]byte(k)); v != nil {
			result[k] = string(v)
		}
	}
	return result
}

func GetPrefix(b *bolt.Bucket, prefix string) map[string]string {
	result := make(map[string]string)
	c := b.Cursor()
	for k, v := c.Seek([]byte(prefix)); k != nil && bytes.HasPrefix(k, []byte(prefix)); k, v = c.Next() {
		result[string(k)] = string(v)
	}
	return result
}

func GetAll(b *bolt.Bucket) map[string]string {
	result := make(map[string]string)
	b.ForEach(func(k, v []byte) error {
		result[string(k)] = string(v)
		return nil
	})
	return result
}

func Put(b *bolt.Bucket, key, value string) error {
	return b.Put([]byte(key), []byte(value))
}

func Delete(b *bolt.Bucket, key string) error {
	return b.Delete([]byte(key))
}
