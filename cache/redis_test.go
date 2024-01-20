package cache

import (
	"context"
	"douyin/config"
	"douyin/global"
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

func Test(t *testing.T) {
	err := config.ConfInit()
	if err != nil {
		panic(err)
		return
	}
	global.Group.Rdb, err = InitRedis()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	rc := MakeRdbCache()
	fmt.Println(rc)
	scorelist := []float64{1111, 2222, 5555, 6666, 3333}
	msgList := []interface{}{"111", "2222", "33333", "44444", "9797"}
	_, err = rc.ZAdd(context.Background(), "zxc66", scorelist, msgList)
	if err != nil {
		panic(err)
		return
	}
	res, err := rc.ZGetRevRangeByScoreWithScores(context.Background(), "zxc66", "-inf", "100000", 0, 2)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
	fmt.Println(res[len(res)-1].Score)
}
