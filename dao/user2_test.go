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
	user, _ := GetUser2Instance().QueryUserInfoByUserID(1)
	fmt.Println(*user)
}

//func TestUserDao2_QueryUserInfoByUserID1(t *testing.T) {
//	DBInit()
//	user,_:=GetUser2Instance().q
//	type args struct {
//		userID int64
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    *User2
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			us := UserDao2{}
//			got, err := us.QueryUserInfoByUserID(tt.args.userID)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("QueryUserInfoByUserID() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("QueryUserInfoByUserID() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
