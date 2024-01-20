package cache

import (
	"context"
	"douyin/dao"
	"douyin/utls"
	"time"
)

func SetUserInfo(c context.Context, key string, value map[string]interface{}) error {
	rc := MakeRdbCache()
	err := rc.HSet(c, key, value)
	if err != nil {
		return err
	}
	_, err = rc.Expire(c, key, time.Hour*168)
	return err
}

func SetUserCount(c context.Context, key string, value map[string]interface{}) error {
	rc := MakeRdbCache()
	err := rc.HSet(c, key, value)
	if err != nil {
		return err
	}
	_, err = rc.Expire(c, key, time.Hour*20)
	return err
}

func ExistUserInfoKey(c context.Context, key string) (bool, error) {
	rc := MakeRdbCache()
	exist, err := rc.Exists(c, key)
	if err != nil {
		return false, err
	}
	return exist > 0, nil
}
func GetUserInfo(c context.Context, key string) (*dao.User2, error) {
	rc := MakeRdbCache()
	userMap, err := rc.HGetAll(c, key)
	if err != nil {
		return nil, err
	}
	if len(userMap) == 0 {
		return nil, nil
	}
	return utls.CreateUserInfo(userMap)
}

func GetUserCount(c context.Context, key string) (*dao.UserCount, error) {
	rc := MakeRdbCache()
	userMap, err := rc.HGetAll(c, key)
	if err != nil {
		return nil, err
	}
	if len(userMap) == 0 {
		return nil, nil
	}
	return utls.CreateUserCount(userMap)
}

func AddWorkCount(c context.Context, key string) error {
	rc := MakeRdbCache()
	_, err := rc.IncrHMCount(c, key, "work_count", 1)
	return err
}

func SubWorkCount(c context.Context, key string) error {
	rc := MakeRdbCache()
	_, err := rc.IncrHMCount(c, key, "work_count", -1)
	return err
}
func AddFavorCount(ctx context.Context, key string) error {
	rc := MakeRdbCache()
	_, err := rc.IncrHMCount(ctx, key, "favorite_count", 1)
	return err
}

func SubFavorCount(ctx context.Context, key string) error {
	rc := MakeRdbCache()
	_, err := rc.IncrHMCount(ctx, key, "favorite_count", -1)
	return err
}

func AddFavoritedCount(ctx context.Context, key string) error {
	rc := MakeRdbCache()
	_, err := rc.IncrHMCount(ctx, key, "total_favorited", 1)
	return err
}

func SubFavoritedCount(ctx context.Context, key string) error {
	rc := MakeRdbCache()
	_, err := rc.IncrHMCount(ctx, key, "total_favorited", -1)
	return err
}
