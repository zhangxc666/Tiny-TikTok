package cache

import (
	"context"
	"douyin/dao"
	"douyin/utls"
	"fmt"
	"strconv"
	"time"
)

func SetVideoInfo(c context.Context, key string, value map[string]any) error {
	rc := MakeRdbCache()
	if err := rc.HSet(c, key, value); err != nil {
		return err
	}
	_, err := rc.Expire(c, key, time.Hour*48)
	return err
}

func SetUserVideoIDs(c context.Context, key string, IDs []int64) error {
	rc := MakeRdbCache()
	str := make([]string, len(IDs))
	for i := range str {
		str[i] = strconv.FormatInt(IDs[i], 10)
	}
	_, err := rc.SAdd(c, key, str)
	if err != nil {
		return err
	}
	_, err = rc.Expire(c, key, time.Hour*48)
	if err != nil {
		return err
	}
	return nil
}

func GetVideoInfo(c context.Context, key string) (*dao.Video, error) {
	rc := MakeRdbCache()
	videoMapInfo, err := rc.HGetAll(c, key)
	if err != nil {
		return nil, err
	}
	if len(videoMapInfo) == 0 {
		return nil, nil
	}
	videoInfo, err := utls.CreateVideoInfo(videoMapInfo)
	if err != nil {
		fmt.Println(videoMapInfo)
		panic(err)
		return nil, err
	}
	return videoInfo, nil
}

func AddPublishVideo(c context.Context, key string, timeStamp float64, videoID int64) error {
	rc := MakeRdbCache()
	if _, err := rc.ZAdd(c, key, []float64{timeStamp}, []any{videoID}); err != nil {
		return err
	}
	return nil
}

func DelUserVideo(c context.Context, key string) error {
	rc := MakeRdbCache()
	if err := rc.Del(c, key); err != nil {
		return err
	}
	return nil
}

func GetUserVideoIDs(c context.Context, key string) ([]int64, error) {
	rc := MakeRdbCache()
	str, err := rc.SGetAll(c, key)
	if err != nil {
		return nil, err
	}
	if len(str) == 0 {
		return nil, nil
	}
	IDs := make([]int64, len(str))
	for i, v := range str {
		IDs[i], err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
	}
	return IDs, nil
}

func GetPublishVideoIDs(c context.Context, key string, lastTime int64) (int64, []int64, error) {
	rc := MakeRdbCache()
	var lastTimeStr string
	lastTimeStr = strconv.FormatInt(lastTime, 10)
	videoStr, err := rc.ZGetRevRangeByScoreWithScores(c, key, "-inf", lastTimeStr, 0, 30)
	if err != nil {
		return -1, nil, err
	}
	if len(videoStr) == 0 {
		return lastTime, nil, nil
	}
	videoIDs := make([]int64, len(videoStr))
	for i, v := range videoStr {
		videoID, _ := strconv.ParseInt(v.Member.(string), 10, 64)
		videoIDs[i] = videoID
	}
	return int64(videoStr[0].Score), videoIDs, nil

}
