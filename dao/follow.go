package dao

import (
	"fmt"
	"sync"
)

type Follow struct {
	//关注者
	FollowId int64 `gorm:"column:follow_id;primary_key"`
	//
	FollowedId int64 `gorm:"column:followed_id;primary_key"`
}

func (Follow) TableName() string {
	return "follow"
}

type FollowDao struct {
}

var followDao *FollowDao

var FollowOnce sync.Once

func GetFollowInstance() *FollowDao {
	FollowOnce.Do(func() {
		followDao = &FollowDao{}
	})
	return followDao
}

func (FollowDao) QueryAllFollow() ([]Follow, error) {

	var FollowLists []Follow
	err := db.Find(&FollowLists).Error
	if err != nil {
		return nil, err
	}
	return FollowLists, err
}

// AddFollow 添加关注映射
func (FollowDao) AddFollow(follow *Follow) error {
	err := db.Create(follow).Error
	if err != nil {
		return err
	}
	//进行关注数量的更新
	return nil
}

// DeleteFollow 删除关注映射
func (FollowDao) DeleteFollow(follow *Follow) error {
	err := db.Where("follow_id = ? and followed_id = ?", follow.FollowId, follow.FollowedId).Delete(follow).Error
	if err != nil {
		return err
	}
	//进行关注数量的更新
	return nil
}

func (FollowDao) QueryFollowLists(userid int64) ([]User, error) {
	var userLists []User
	err := db.Raw("SELECT * FROM  `user` WHERE user.user_id IN ( SELECT follow.followed_id FROM follow WHERE follow.follow_id = ? )", userid).Scan(&userLists).Error
	if err != nil {
		return nil, err
	}
	fmt.Printf("%#v", userLists)
	return userLists, nil
}

func (FollowDao) QueryFollowerLists(userid int64) ([]User, error) {
	var userLists []User
	err := db.Raw("SELECT * FROM  `user` WHERE user.user_id IN ( SELECT follow.follow_id FROM follow WHERE follow.followed_id = ? )", userid).Scan(&userLists).Error
	if err != nil {
		return nil, err
	}
	fmt.Printf("%#v", userLists)
	return userLists, nil
}

// QueryAllFollowID 查询所有关注的ID
func (FollowDao) QueryAllFollowID(userID int64) ([]int64, error) {
	var followIDs []int64
	err := db.Raw("SELECT follow.followed_id FROM `follow` WHERE follow.follow_id = ?", userID).Scan(&followIDs).Error
	if err != nil {
		return nil, err
	}

	fmt.Printf("%#v", followIDs)
	return followIDs, nil
}

// QueryAllFanID 查询所有粉丝的ID
func (FollowDao) QueryAllFanID(userID int64) ([]int64, error) {
	var followIDs []int64
	err := db.Raw("SELECT follow.follow_id FROM `follow` WHERE follow.followed_id = ?", userID).Scan(&followIDs).Error
	if err != nil {
		return nil, err
	}
	fmt.Printf("%#v", followIDs)
	return followIDs, nil
}

// QueryAllFriendID 查询所有好友ID
func (FollowDao) QueryAllFriendID(userID int64) ([]int64, error) {
	var friendIDs []int64
	err := db.Raw("SELECT f1.follow_id\nFROM Follow f1\nJOIN Follow f2 ON f1.followed_id = f2.follow_id AND f1.follow_id = f2.followed_id\nWHERE f1.followed_id = ?;", userID).Scan(&friendIDs).Error
	if err != nil {
		return nil, nil
	}
	return friendIDs, nil
}

func (FollowDao) IsFollow(userID, targetID int64) (bool, error) {
	var count int64
	err := db.Model(&Follow{}).Where("follow_id = ? and followed_id = ?", userID, targetID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (FollowDao) QueryEachFollow(userid int64) ([]User, error) {
	var userLists []User
	err := db.Raw("SELECT * FROM  `user` WHERE user.user_id != ? and user.user_id in \n(SELECT DISTINCT follow.followed_id FROM  follow \njoin\n(SELECT follow.follow_id FROM follow WHERE follow.followed_id = ?) a \non\na.follow_id = follow.followed_id) ", userid, userid).Scan(&userLists).Error
	if err != nil {
		return nil, err
	}
	fmt.Printf("%#v", userLists)
	return userLists, nil
}

func (FollowDao) IsBothFollow(userid, friendID int64) bool {
	cnt := int64(-1)
	db.Debug().Model(&Follow{}).Where("( follow_id = ? and followed_id = ?  ) or ( follow_id = ?  and followed_id = ? )", userid, friendID, friendID, userid).Count(&cnt)
	return cnt == 2
}
