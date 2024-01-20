package cache

import (
	"context"
	"douyin/config"
	"douyin/global"
	"douyin/utls"
	"fmt"
	"strconv"
	"testing"
)

func TestAddPublishVideo(t *testing.T) {
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
	str, err := rc.ZGetRevRangeByScoreWithScores(context.Background(), utls.CreatePublishKey(), "-inf", strconv.FormatInt(1705759420416, 10), 0, 30)
	if err != nil {
		panic(err)
	}
	t.Log(str)
}
