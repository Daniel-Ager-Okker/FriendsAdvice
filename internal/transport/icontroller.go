package transport

import "time"

type IController interface {
	IsStorageReady() bool
	PutObjectWithExpires(key uint64, value []byte, lifetime time.Duration) (bool, error)
	PutObject(key uint64, value []byte) (bool, error)
	GetObject(key uint64) ([]byte, bool)
}
