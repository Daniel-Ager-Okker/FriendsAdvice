package transport

import "time"

type IController interface {
	IsStorageReady() bool
	PutObjectWithExpires(key int, value string, lifetime time.Duration) bool
	PutObject(key int, value string) bool
	GetObject(key int) (string, bool)
}
