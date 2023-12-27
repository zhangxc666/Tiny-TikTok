package service

import (
	"context"
	"douyin/cache"
	"douyin/dao"
	"douyin/utls"
)

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

func IsFollowTarget(ctx context.Context, userID int64, targetID int64) (bool, error) {
	if userID == targetID {
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

func AddFollowCount(ctx context.Context, userID, targetID int64) error {
	userCountKey := utls.CreateUserCountKey(userID)
	targetCountKey := utls.CreateUserCountKey(targetID)
	if err := cache.AddFollowCount(ctx, userCountKey); err != nil {
		return err
	}
	if err := cache.AddFollowerCount(ctx, targetCountKey); err != nil {
		return err
	}
	return nil
}

func SubFollowCount(ctx context.Context, userID, targetID int64) error {
	userCountKey := utls.CreateUserCountKey(userID)
	targetCountKey := utls.CreateUserCountKey(targetID)
	if err := cache.SubFollowCount(ctx, userCountKey); err != nil {
		return err
	}
	if err := cache.SubFollowerCount(ctx, targetCountKey); err != nil {
		return err
	}
	return nil
}
