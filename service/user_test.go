package service

import (
	"context"
	"douyin/cache"
	"douyin/dao"
	"fmt"
	"testing"
)

func TestRegister(t *testing.T) {
	err := dao.DbInit()
	if err != nil {
		return
	}
	err = cache.RedisPoolInit()
	if err != nil {
		t.Error(err.Error())
		return
	}
	type args struct {
		username string
		password string
	}
	tests := []struct {
		name string
		args
	}{
		{
			"测试小王",
			args{
				"____________",
				"123456mksjxnjancanjskandndjasnjdkasn",
			},
		},
		{
			"测试小贺",
			args{
				"小王",
				"123456",
			},
		},
		{
			"测试小陈",
			args{
				"",
				"123456",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			userinfo, err := Register(context.Background(), test.username, test.password)
			if err != nil {
				t.Errorf("UserRegister ERROR is %v", err)
				return
			}
			fmt.Println(userinfo)
		})
	}
}

func TestLogin(t *testing.T) {
	err := dao.DbInit()
	if err != nil {
		t.Error(err.Error())
		return
	}
	type args struct {
		username string
		password string
	}
	tests := []struct {
		name string
		args
	}{
		{
			"测试1",
			args{
				"小王",
				"123456",
			},
		},
		{
			"测试2",
			args{
				"_____",
				"123456",
			},
		},
		{
			"测试3",
			args{
				"hhhy",
				"123456",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			userinfo, err := Login(context.Background(), test.username, test.password)
			if err != nil {
				t.Errorf("UserRegister ERROR is %v", err)
				return
			}
			fmt.Println(userinfo)
		})
	}
}
