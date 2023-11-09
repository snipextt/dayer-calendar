package utils

import (
	"log"
	"time"

	"github.com/snipextt/dayer/cmd"
	"github.com/snipextt/dayer/storage"
)

func SetOAuthState(uid string) (s string, err error) {
	ctx, cancel := NewContext()
	defer cancel()
	s = cmd.GetRandomString(32)
	client := storage.GetRedisInstance()
	err = client.Set(ctx, uid, s, time.Minute*5).Err()
	if err != nil {
		return "", err
	}
	return s, nil
}

func GetOAuthState(uid string) (s string, err error) {
	ctx, cancel := NewContext()
	defer cancel()
	client := storage.GetRedisInstance()
	s, err = client.Get(ctx, uid).Result()
	log.Println(s)
	return
}
