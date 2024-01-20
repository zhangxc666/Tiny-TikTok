package service

import (
	"context"
	"douyin/cache"
	"douyin/dao"
	"douyin/utls"
	"errors"
	"time"
)

func FollowOrCancel(ctx context.Context, userID int64, targetID int64, action string) error {
	if userID == targetID {
		return errors.New("无法自己关注自己")
	}
	exist, err := dao.GetFollowInstance().IsFollow(userID, targetID)
	if err != nil {
		return err
	}
	switch action {
	case "1": // 关注
		if exist {
			return errors.New("无法关注已关注的用户")
		}
		// 延迟双删-第一次删除
		var isFollow bool
		if isFollow, err = IsFollowTarget(ctx, targetID, userID); err != nil {
			return err
		} else if isFollow {
			friendKey1, friendKey2 := utls.CreateFriendKey(userID), utls.CreateFriendKey(targetID)
			if err := cache.DelFriend(ctx, friendKey1); err != nil {
				return err
			}
			if err := cache.DelFriend(ctx, friendKey2); err != nil {
				return err
			}
		}
		followKey, fanKey := utls.CreateFollowKey(userID), utls.CreateFanKey(targetID)
		if err := cache.DelFollow(ctx, followKey); err != nil {
			return err
		}
		if err := cache.DelFan(ctx, fanKey); err != nil {
			return err
		}
		// 判断是否对方关注了自己，如果关注则双删好友缓存

		// 操作数据库
		if err := dao.GetFollowInstance().AddFollow(&dao.Follow{FollowId: userID, FollowedId: targetID}); err != nil {
			return err
		}
		time.Sleep(time.Millisecond * 100)
		// 第二次删除
		if err := cache.DelFollow(ctx, followKey); err != nil {
			return err
		}
		if err := cache.DelFan(ctx, fanKey); err != nil {
			return err
		}
		if isFollow {
			friendKey1, friendKey2 := utls.CreateFriendKey(userID), utls.CreateFriendKey(targetID)
			if err := cache.DelFriend(ctx, friendKey1); err != nil {
				return err
			}
			if err := cache.DelFriend(ctx, friendKey2); err != nil {
				return err
			}
		}
		// 增加count
		if err := AddFollowCount(ctx, userID, targetID); err != nil {
			return err
		}
	case "2":
		if !exist {
			return errors.New("当前用户未关注，无法取消关注")
		}
		followKey, fanKey := utls.CreateFollowKey(userID), utls.CreateFanKey(targetID)
		// 第一次删除
		var isFollow bool
		if isFollow, err = IsFollowTarget(ctx, userID, targetID); err != nil {
			return err
		} else if isFollow {
			friendKey1, friendKey2 := utls.CreateFriendKey(userID), utls.CreateFriendKey(targetID)
			if err := cache.DelFriend(ctx, friendKey1); err != nil {
				return err
			}
			if err := cache.DelFriend(ctx, friendKey2); err != nil {
				return err
			}
		}
		if err := cache.DelFollow(ctx, followKey); err != nil {
			return err
		}
		if err := cache.DelFan(ctx, fanKey); err != nil {
			return err
		}
		// 操作数据库
		if err := dao.GetFollowInstance().DeleteFollow(&dao.Follow{FollowId: userID, FollowedId: targetID}); err != nil {
			return err
		}
		time.Sleep(time.Millisecond * 100)
		// 第二次删除
		if err := cache.DelFollow(ctx, followKey); err != nil {
			return err
		}
		if err := cache.DelFan(ctx, fanKey); err != nil {
			return err
		}
		if isFollow {
			friendKey1, friendKey2 := utls.CreateFriendKey(userID), utls.CreateFriendKey(targetID)
			if err := cache.DelFriend(ctx, friendKey1); err != nil {
				return err
			}
			if err := cache.DelFriend(ctx, friendKey2); err != nil {
				return err
			}
		}
		if err := SubFollowCount(ctx, userID, targetID); err != nil {
			return err
		}
	default:
		return errors.New("未知操作")
	}
	return nil
}

// GetFollowList 获取所有关注的人的信息
func GetFollowList(ctx context.Context, userID int64) ([]dao.User, error) {
	followKey := utls.CreateFollowKey(userID)
	followIDs, err := cache.GetAllMembersByKey(ctx, followKey)
	if err != nil {
		return nil, err
	}
	if followIDs == nil {
		// 查询数据库中的ID
		followIDs, err = dao.GetFollowInstance().QueryAllFollowID(userID)
		if err != nil {
			return nil, err
		}
		// 如果没有关注的
		if len(followIDs) == 0 {
			return nil, nil
		}
		// 添加至cache
		err = cache.AddManyFollows(ctx, followKey, followIDs)
		if err != nil {
			return nil, err
		}
	}
	// 查数据库中的user信息
	userList, err := GetUserList(ctx, userID, followIDs)
	if err != nil {
		return nil, err
	}
	return userList, nil
}

// GetFanList 获取粉丝列表
func GetFanList(ctx context.Context, userID int64) ([]dao.User, error) {
	fanKey := utls.CreateFanKey(userID)
	fanIDs, err := cache.GetAllMembersByKey(ctx, fanKey)
	if err != nil {
		return nil, err
	}
	if fanIDs == nil {
		fanIDs, err = dao.GetFollowInstance().QueryAllFanID(userID)
		if err != nil {
			return nil, err
		}
		if len(fanIDs) == 0 {
			return nil, nil
		}
		err = cache.AddManyFans(ctx, fanKey, fanIDs)
		if err != nil {
			return nil, err
		}
	}
	// 查数据库中的user信息
	userList, err := GetUserList(ctx, userID, fanIDs)
	if err != nil {
		return nil, err
	}
	return userList, nil
}

// GetFriendList 获取朋友列表
func GetFriendList(ctx context.Context, userID int64) ([]dao.User, error) {
	friendKey := utls.CreateFriendKey(userID)
	friendIDs, err := cache.GetAllMembersByKey(ctx, friendKey)
	if err != nil {
		return nil, err
	}
	if friendIDs == nil {
		friendIDs, err = dao.GetFollowInstance().QueryAllFriendID(userID)
		if err != nil {
			return nil, err
		}
		if len(friendIDs) == 0 {
			return nil, nil
		}
		err = cache.AddManyFriends(ctx, friendKey, friendIDs)
		if err != nil {
			return nil, err
		}
	}
	userList, err := GetUserList(ctx, userID, friendIDs)
	if err != nil {
		return nil, err
	}
	return userList, nil
}

// IsFollowTarget 判断是否关注了对方
func IsFollowTarget(ctx context.Context, userID int64, targetID int64) (bool, error) {

	if userID == targetID || userID == -1 {
		return false, nil
	}
	followKey := utls.CreateFollowKey(userID)
	exist, err := cache.IsFollow(ctx, followKey, targetID)
	if err != nil {
		return false, err
	}
	if exist {
		return true, err
	}
	keyExist, err := cache.KeyExists(ctx, followKey)
	if err != nil {
		return false, err
	}
	if keyExist { // 缓存未过期，未找到对应的targetID
		return false, nil
	}
	// 去数据库中查找对应的关系
	exist, err = dao.GetFollowInstance().IsFollow(userID, targetID) // 是否存在
	if err != nil {
		return false, err
	}
	// 获取所有关注人ID
	targetIDs, err := dao.GetFollowInstance().QueryAllFollowID(userID)
	if err != nil {
		return false, err
	}
	if targetIDs == nil {
		return false, nil
	}
	// 添加关注
	err = cache.AddManyFollows(ctx, followKey, targetIDs)
	if err != nil {
		return exist, err
	}
	return exist, nil
}

// IsFollowManyTargets 判断是否关注了多人
func IsFollowManyTargets(ctx context.Context, userID int64, targetIDs []int64) ([]bool, error) {
	isFollowList := make([]bool, len(targetIDs))
	for i := range targetIDs {
		target, err := IsFollowTarget(ctx, userID, targetIDs[i])
		if err != nil {
			return nil, err
		}
		isFollowList[i] = target
	}
	return isFollowList, nil
}

// AddFollowCount 增加关注数（写回）
func AddFollowCount(ctx context.Context, userID, targetID int64) error {
	userCountKey := utls.CreateUserCountKey(userID)
	targetCountKey := utls.CreateUserCountKey(targetID)
	if err := SetUserCountToCache(ctx, userID); err != nil {
		return err
	}
	if err := SetUserCountToCache(ctx, targetID); err != nil {
		return err
	}
	if err := cache.AddFollowCount(ctx, userCountKey); err != nil {
		return err
	}
	if err := cache.AddFollowerCount(ctx, targetCountKey); err != nil {
		return err
	}
	return nil
}

// SubFollowCount 减少关注数（写回）
func SubFollowCount(ctx context.Context, userID, targetID int64) error {
	userCountKey := utls.CreateUserCountKey(userID)
	targetCountKey := utls.CreateUserCountKey(targetID)
	if err := SetUserCountToCache(ctx, userID); err != nil {
		return err
	}
	if err := SetUserCountToCache(ctx, targetID); err != nil {
		return err
	}
	if err := cache.SubFollowCount(ctx, userCountKey); err != nil {
		return err
	}
	if err := cache.SubFollowerCount(ctx, targetCountKey); err != nil {
		return err
	}
	return nil
}
