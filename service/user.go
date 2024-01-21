package service

import (
	"context"
	"douyin/cache"
	"douyin/dao"
	"douyin/utls"
	"errors"
	"fmt"
)

type UserRegisInfo struct {
	Token  string `json:"token"`   // 用户鉴权token
	UserID int64  `json:"user_id"` // 用户id
}

func GetUserList(ctx context.Context, userID int64, targetIDs []int64) ([]dao.User, error) {
	userList := make([]dao.User, len(targetIDs))
	isFollowList, err := IsFollowManyTargets(ctx, userID, targetIDs)
	if err != nil {
		return nil, err
	}
	for i, targetID := range targetIDs {
		userInfoKey, userCountKey := utls.CreateUserInfoKey(targetID), utls.CreateUserCountKey(targetID)
		userInfo, err := cache.GetUserInfo(ctx, userInfoKey)
		if err != nil {
			return nil, err
		}
		if userInfo == nil {
			userInfo, err = dao.GetUser2Instance().QueryUserInfoByUserID(targetID)
			if err != nil {
				return nil, err
			}
			err = cache.SetUserInfo(ctx, userInfoKey, utls.CreateMapUserInfo(userInfo))
			if err != nil {
				return nil, err
			}
			fmt.Println("userInfo:", userInfo)
		} else {
			if err := SetUserCountToCache(ctx, targetID); err != nil {
				return nil, err
			}
			userCount, err := cache.GetUserCount(ctx, userCountKey)
			if err != nil {
				return nil, err
			}
			fmt.Println(userCount, "usercount:")
			userInfo.Usercount = userCount
		}
		userInfo.IsFollow = isFollowList[i]
		userList[i] = *utls.ChangeUser2ToUser(userInfo)
	}
	return userList, nil
}
func Register(ctx context.Context, username, password string) (*UserRegisInfo, error) {
	exist, err := dao.GetUser2Instance().ExistUserByUsername(username)
	if err != nil || exist {
		if exist {
			return nil, errors.New("user existed")
		}
		return nil, err
	}
	hashPassword := utls.Md5Encryption(password)
	emptyUserCount := &dao.UserCount{
		FollowCount:    0,
		FollowerCount:  0,
		TotalFavorited: 0,
		FavoriteCount:  0,
		WorkCount:      0,
	}
	user := dao.User2{
		Name:            username,
		Password:        hashPassword,
		Avatar:          "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg",
		BackGroundImage: "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg",
		Signature:       "zxc666",
		Usercount:       emptyUserCount,
	}
	err = dao.GetUser2Instance().AddUser(&user)
	if err != nil {
		return nil, err
	}
	userInfoKey := utls.CreateUserInfoKey(user.ID)
	userCountKey := utls.CreateUserCountKey(user.ID)
	userInfoMap, userCountMap := utls.CreateMapUserInfo(&user), utls.CreateMapUserCount(emptyUserCount)
	err = cache.SetUserInfo(ctx, userInfoKey, userInfoMap) // 设置个人信息缓存
	if err != nil {
		return nil, err
	}
	err = cache.SetUserCount(ctx, userCountKey, userCountMap) // 设置计数缓存
	if err != nil {
		return nil, err
	}
	var info UserRegisInfo
	token, err := utls.GenerateToken(username, user.ID)
	if err != nil {
		return nil, err
	}
	info.Token = token
	info.UserID = user.ID
	return &info, nil
}

func Login(ctx context.Context, username, password string) (*UserRegisInfo, error) {
	userInfo, err := dao.GetUser2Instance().QueryUserInfoByUsername(username)
	if err != nil {
		return nil, err
	}
	fmt.Println(userInfo.Name)
	if utls.CheckPassword(password, userInfo.Password) == false {
		return nil, errors.New("用户名或密码错误")
	}
	var info UserRegisInfo
	token, err := utls.GenerateToken(username, userInfo.ID)
	if err != nil {
		return nil, err
	}
	info.Token = token
	info.UserID = userInfo.ID
	return &info, nil
}

func GetUserIndex(ctx context.Context, userID, targetID int64) (*dao.User2, error) {
	var isFollow bool
	var err error
	if userID == targetID {
		isFollow = false
	} else {
		isFollow, err = dao.GetFollowInstance().IsFollow(userID, targetID)
		if err != nil {
			return nil, err
		}
	}
	userInfoKey, userCountKey := utls.CreateUserInfoKey(userID), utls.CreateUserCountKey(userID)
	userInfo, err := cache.GetUserInfo(ctx, userInfoKey)
	if err != nil {
		return nil, err
	}
	if err := SetUserCountToCache(ctx, userID); err != nil {
		return nil, err
	}
	userCount, err := cache.GetUserCount(ctx, userCountKey)
	if err != nil {
		return nil, err
	}
	if userInfo != nil {
		userInfo.Usercount = userCount
		userInfo.IsFollow = isFollow
		return userInfo, nil
	}
	user, err := dao.GetUser2Instance().QueryUserInfoByUserID(userID)
	if err != nil {
		return nil, err
	}
	if err = cache.SetUserInfo(ctx, userInfoKey, utls.CreateMapUserInfo(user)); err != nil {
		return nil, err
	}
	if err = cache.SetUserCount(ctx, userCountKey, utls.CreateMapUserCount(user.Usercount)); err != nil {
		return nil, err
	}
	return user, nil
}

func SetUserCountToCache(ctx context.Context, userID int64) error {
	userKey := utls.CreateUserCountKey(userID)
	if exist, err := cache.KeyExists(ctx, userKey); err != nil {
		return err
	} else if exist == false {
		userCount, err := dao.GetUser2Instance().QueryUserCountByUserID(userID)
		if err != nil {
			return err
		}
		if err := cache.SetUserCount(ctx, userKey, utls.CreateMapUserCount(userCount)); err != nil {
			return err
		}
	}
	return nil
}

func AddWorkCount(ctx context.Context, userID int64) error {
	userCountKey := utls.CreateUserCountKey(userID)
	if err := SetUserCountToCache(ctx, userID); err != nil {
		return err
	}
	if err := cache.AddWorkCount(ctx, userCountKey); err != nil {
		return err
	}
	return nil
}

func SubWorkCount(ctx context.Context, userID int64) error {
	userCountKey := utls.CreateUserCountKey(userID)
	if err := SetUserCountToCache(ctx, userID); err != nil {
		return err
	}
	if err := cache.SubWorkCount(ctx, userCountKey); err != nil {
		return err
	}
	return nil
}

func AddUserFavorCount(ctx context.Context, userID, targetID int64) error {
	userCountKey := utls.CreateUserCountKey(userID)
	targetCountKey := utls.CreateUserCountKey(targetID)
	if err := cache.AddFavorCount(ctx, userCountKey); err != nil {
		return err
	}
	if err := cache.AddFavoritedCount(ctx, targetCountKey); err != nil {
		return err
	}
	return nil
}

func SubUserFavorCount(ctx context.Context, userID, targetID int64) error {
	userCountKey, targetCountKey := utls.CreateUserCountKey(userID), utls.CreateUserCountKey(targetID)
	if err := cache.SubFavorCount(ctx, userCountKey); err != nil {
		return err
	}
	if err := cache.SubFavoritedCount(ctx, targetCountKey); err != nil {
		return err
	}
	return nil
}
