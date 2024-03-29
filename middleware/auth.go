package middleware

import (
	"douyin/cache"
	"douyin/common"
	"douyin/dao"
	"douyin/handle"
	"douyin/utls"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserAuth() gin.HandlerFunc {

	return func(c *gin.Context) {
		//得到token字段
		//1.get请求
		token := c.Query("token")
		fmt.Println("token", token)
		if token == "" {
			//2.post请求
			token = c.PostForm("token")
		}

		// 两种情况下来，判断是否有token
		if token == "" {
			c.JSON(http.StatusOK, handle.UserResponse{
				Response: common.Response{StatusCode: 1, StatusMsg: "token 不存在"},
			})
			//终止
			c.Abort()
			return
		}
		//解析
		t, claim, err := utls.ParseToken(token)

		//判断是否有效
		if !t.Valid || err != nil {
			c.JSON(http.StatusOK, handle.UserResponse{
				Response: common.Response{StatusCode: 1, StatusMsg: "token有效期过了或者" + err.Error()},
			})
			c.Abort()
			return
		}

		//1.首先到redis中查找，没有的话去mysql中查找
		//2.mysql中没有说明token失败
		var isExists = true
		userKey := utls.CreateUserInfoKey(claim.UserId)
		isExists, err = cache.ExistUserInfoKey(c, userKey)
		if err != nil {
			fmt.Println(err)
			//在redis中不存在
			isExists = false
		}
		if !isExists {
			//进行db查找
			exist, err := dao.GetUser2Instance().ExistUserByUserID(claim.UserId)
			if err != nil {
				c.JSON(http.StatusOK, handle.UserResponse{
					Response: common.Response{StatusCode: 1, StatusMsg: "token find failed"},
				})
				c.Abort()
				return
			}
			if exist == false {
				c.JSON(http.StatusOK, handle.UserResponse{
					Response: common.Response{StatusCode: 1, StatusMsg: "id is not exists"},
				})
				c.Abort()
				return
			}
		}

		//传给handle层
		c.Set("userid", claim.UserId)
		//挂起
		c.Next()
	}
}
