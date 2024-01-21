package cache

import (
	"context"
	"douyin/utls"
	"strconv"
	"strings"
	"time"
)

func SetUserFavorVideoIDs(ctx context.Context, key string, videoIDs []int64) error {
	Str := make([]string, len(videoIDs))
	for i, v := range videoIDs {
		Str[i] = strconv.FormatInt(v, 10)
	}
	rc := MakeRdbCache()
	_, err := rc.SAdd(ctx, key, Str)
	if err != nil {
		return err
	}
	_, err = rc.Expire(ctx, key, time.Hour*48)
	return err
}

func ExistFavorRecord(ctx context.Context, key string, videoID int64) (bool, error) {
	rc := MakeRdbCache()
	member := utls.CreateFavorMember(videoID, 1)
	// 存储记录，原因是当前缓存是写回
	// 与follow不同的是，关注采用的策略是延迟双删
	return rc.SIsExist(ctx, key, member)
}

func AddFavorRecord(ctx context.Context, key string, videoID int64) error {
	rc := MakeRdbCache()
	cancelStr := utls.CreateFavorMember(videoID, 2)
	_, err := rc.SRem(ctx, key, cancelStr)
	if err != nil {
		return err
	}
	favorStr := utls.CreateFavorMember(videoID, 1)
	_, err = rc.SAdd(ctx, key, favorStr)
	return err
}

func AddCancelFavorRecord(ctx context.Context, key string, videoID int64) error {
	rc := MakeRdbCache()
	favorStr := utls.CreateFavorMember(videoID, 1)
	_, err := rc.SRem(ctx, key, favorStr)
	if err != nil {
		return err
	}
	cancelFavorStr := utls.CreateFavorMember(videoID, 2)
	_, err = rc.SAdd(ctx, key, cancelFavorStr)
	return err
}

func GetAllFavoriteVideosIDs(ctx context.Context, userID int64) ([]int64, error) {
	rc := MakeRdbCache()
	videoStr, err := rc.SGetAll(ctx, utls.CreateFavorKey(userID))
	if err != nil {
		return nil, err
	}
	if len(videoStr) == 0 {
		return nil, nil
	}
	videoIDs := make([]int64, 0)
	for _, v := range videoStr {
		str := strings.Split(v, "+")
		flag := str[0]
		videoID, _ := strconv.ParseInt(str[1], 10, 64)
		if flag == "2" {
			continue
		}
		videoIDs = append(videoIDs, videoID)
	}
	return videoIDs, nil
}
