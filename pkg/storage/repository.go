package storage

import "google.golang.org/genproto/googleapis/storage/v1"

type Bucket string

const (
	AccessTokens Bucket = "my_AccessToken"
	RequestToken Bucket = "my_RequestToken"
)

type TokenStorage interface {
	Save(chatID int64, token string, bucket storage.Bucket) error
	Get(chatID int64, token string, bucket storage.Bucket) error
}
