package redis

import "github.com/omniful/go_commons/redis"

type Redis struct{
	*redis.Client

}
var RD *Redis

func GetClient()*Redis{
	return RD

}
func SetClient(client *redis.Client){

	RD=&Redis{client}

}