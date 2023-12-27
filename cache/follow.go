package cache

import (
	"context"
	"strconv"
	"time"
)

func KeyExists(c context.Context, key string) (bool, error) {
	rc := MakeRdbCache()
	cnt, err := rc.Exists(c, key)
	if err != nil {
		return false, err
	}
	if cnt != 1 {
		return false, nil
	}
	return true, nil
}

// AddManyFollows 添加关注的ID至缓存中
func AddManyFollows(c context.Context, key string, IDs []int64) error {
	Str := make([]string, len(IDs))
	for i, v := range IDs {
		Str[i] = strconv.FormatInt(v, 10)
	}
	rc := MakeRdbCache()
	_, err := rc.SAdd(c, key, Str)
	if err != nil {
		return err
	}
	_, err = rc.Expire(c, key, time.Hour*48)
	return err
}

// AddManyFans 添加自己的粉丝至缓存中
func AddManyFans(c context.Context, key string, IDs []int64) error {
	Str := make([]string, len(IDs))
	for i, v := range IDs {
		Str[i] = strconv.FormatInt(v, 10)
	}
	rc := MakeRdbCache()
	_, err := rc.SAdd(c, key, Str)
	if err != nil {
		return err
	}
	_, err = rc.Expire(c, key, time.Hour*48)
	return err
}

// AddManyFriends 添加自己的好友至缓存中
func AddManyFriends(c context.Context, key string, IDs []int64) error {
	Str := make([]string, len(IDs))
	for i, v := range IDs {
		Str[i] = strconv.FormatInt(v, 10)
	}
	rc := MakeRdbCache()
	_, err := rc.SAdd(c, key, Str)
	if err != nil {
		return err
	}
	_, err = rc.Expire(c, key, time.Hour*48)
	return err
}

func AddFollow(c context.Context, key string, targetID []int64) error {
	rc := MakeRdbCache()
	_, err := rc.SAdd(c, key, targetID)
	if err != nil {
		return err
	}
	_, err = rc.Expire(c, key, time.Hour*48)
	return err
}

func AddFollowCount(c context.Context, userCountKey string) error {
	rc := MakeRdbCache()
	_, err := rc.IncrHMCount(c, userCountKey, "follow_count", 1)
	return err
}

func AddFollowerCount(c context.Context, userCountKey string) error {
	rc := MakeRdbCache()
	_, err := rc.IncrHMCount(c, userCountKey, "follower_count", 1)
	return err
}

func SubFollowCount(c context.Context, userCountKey string) error {
	rc := MakeRdbCache()
	_, err := rc.IncrHMCount(c, userCountKey, "follow_count", -1)
	return err
}

func SubFollowerCount(c context.Context, userCountKey string) error {
	rc := MakeRdbCache()
	_, err := rc.IncrHMCount(c, userCountKey, "follower_count", -1)
	return err
}
func DelFollow(c context.Context, key string) error {
	rc := MakeRdbCache()
	return rc.Del(c, key)
}

func GetAllMembersByKey(c context.Context, key string) ([]int64, error) {
	rc := MakeRdbCache()
	ids, err := rc.SGetAll(c, key)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, nil
	}
	idList := make([]int64, len(ids))
	for i := range idList {
		idList[i], _ = strconv.ParseInt(ids[i], 10, 64)
	}
	return idList, nil
}

func IsFollow(c context.Context, key string, targetID int64) (bool, error) {
	rc := MakeRdbCache()
	return rc.SIsExist(c, key, targetID)
}
