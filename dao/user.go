package dao

import (
	"gorm.io/gorm"
	"sync"
)

type User2 struct {
	ID              int64      `gorm:"column:user_id" json:"id,string"`
	Name            string     `gorm:"column:name"    json:"name"`
	Password        string     `gorm:"column:password"       json:"-"`
	IsFollow        bool       `gorm:"-"                     json:"-" `
	Avatar          string     `gorm:"column:avatar"         json:"avatar"`
	BackGroundImage string     `gorm:"column:background_image"        json:"background_image"`
	Signature       string     `gorm:"column:signature"               json:"signature"`
	Usercount       *UserCount `gorm:"foreignKey:UserID"`
	VideoLieLists   []Video    `gorm:"many2many:like;" json:"-"`
}

type User struct {
	ID              int64   `gorm:"column:user_id" json:"id"`
	Name            string  `gorm:"column:name"    json:"name"`
	FollowCount     int64   `gorm:"column:follow_count"   json:"follow_count"`
	FollowerCount   int64   `gorm:"column:follower_count" json:"follower_count"`
	Password        string  `gorm:"column:password"       json:"-"`
	IsFollow        bool    `gorm:"-"                     json:"is_follow" `
	Avatar          string  `gorm:"column:avatar"         json:"avatar"`
	BackGroundImage string  `gorm:"column:background_image"        json:"background_image"`
	Signature       string  `gorm:"column:signature"               json:"signature"`
	TotalFavorite   int64   `gorm:"-"               json:"total_favorited"`
	WorkCount       int64   `gorm:"-"               json:"work_count"`
	FavoriteCount   int64   `gorm:"-"               json:"favorite_count"`
	VideoLieLists   []Video `gorm:"many2many:like;" json:"-"`
}
type UserCount struct {
	gorm.Model
	UserID         int64 `json:"user_id,string"`
	FollowCount    int64 `gorm:"column:follow_count"   json:"follow_count,string"`
	FollowerCount  int64 `gorm:"column:follower_count" json:"follower_count,string"`
	TotalFavorited int64 `gorm:"total_favorited"       json:"total_favorited,string"`
	WorkCount      int64 `gorm:"word_count"            json:"work_count,string"`
	FavoriteCount  int64 `gorm:"favorite_count"        json:"favorite_count,string"`
}

func (User2) TableName() string {
	return "User2"
}

func (UserCount) TableName() string {
	return "user_count"
}

type UserDao2 struct {
}

var userDao2 *UserDao2
var user2Once sync.Once

func GetUser2Instance() *UserDao2 {
	user2Once.Do(func() {
		userDao2 = &UserDao2{}
	})
	return userDao2
}

func (UserDao2) AddUser(user *User2) error {
	res := db.Create(user)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (UserDao2) AddCount(count *UserCount) error {
	res := db.Create(count)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (UserDao2) UpdateCount(count *UserCount) error {
	res := db.Model(count).Where("user_id = ?", count.UserID).Updates(map[string]interface{}{"follow_count": count.FollowCount, "follower_count": count.FollowerCount})
	return res.Error
}

func (UserDao2) ExistUserByUsername(username string) (bool, error) {
	var count int64
	err := db.Model(&User2{}).Where("name = ?", username).Count(&count).Error
	return count > 0, err
}

func (UserDao2) ExistUserByUserID(userID int64) (bool, error) {
	var count int64
	err := db.Model(&User2{}).Where("user_id = ?", userID).Count(&count).Error
	return count > 0, err
}

func (UserDao2) QueryUserInfoByUsername(username string) (*User2, error) {
	var user User2
	err := db.Preload("Usercount").Preload("VideoLieLists").Model(&User2{}).Where("name = ?", username).Find(&user).Error
	return &user, err
}

func (UserDao2) QueryUserInfoByUserID(userID int64) (*User2, error) {
	var user User2
	err := db.Debug().Preload("Usercount").Preload("VideoLieLists").Model(&User2{}).Where("user_id = ?", userID).Find(&user).Error
	return &user, err
}

func (UserDao2) QueryUserCountByUserID(userID int64) (*UserCount, error) {
	var userCount UserCount
	err := db.Model(&UserCount{}).Where("user_id = ?", userID).Find(&userCount).Error
	return &userCount, err
}
