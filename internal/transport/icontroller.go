package transport

import "time"

type IController interface {
	IsStorageReady() bool
	PutObjectWithExpires(key, value string, lifetime time.Duration) bool
	PutObject(key, value string) bool
}
