package cache

import (
	"douyin/config"
	"fmt"
	"testing"
)

func TestSetUserCount(t *testing.T) {
	err := config.ConfInit()
	if err != nil {
		panic(err)
		return
	}
	err = RedisPoolInit()
	if err != nil {
		panic(err)
		return
	}
	if err != nil {
	}
	j := IsUserRelation(1, 2)
	fmt.Println(j)
	err = IncrByUserTotalFavorite(8888)

}
