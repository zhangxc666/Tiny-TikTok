package dao

import (
	"fmt"
	"testing"
)

func TestUserDao2_AddUser(t *testing.T) {
	DBInit()
	GetUser2Instance().AddUser(&User2{
		Name:     "zxc2",
		Password: "123456",
		Avatar:   "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg",
		Usercount: &UserCount{
			FollowCount:    1,
			FollowerCount:  1,
			TotalFavorited: 1,
			WorkCount:      2,
			FavoriteCount:  2,
		},
	})
}
func TestUserDao2_QueryUserInfoByUserID(t *testing.T) {
	DBInit()
	user, _ := GetUser2Instance().QueryUserInfoByUserID(10)
	fmt.Println(user.Usercount)
}
