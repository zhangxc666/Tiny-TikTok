package handle

import (
	"douyin/common"
	"douyin/dao"
	"douyin/service"
	"douyin/utls"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type UserRegisterResponse struct {
	common.Response
	Token  string `json:"token"`   // 用户鉴权token
	UserId int64  `json:"user_id"` // 用户id
}

type UserResponse struct {
	common.Response
	User dao.User `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	if password == "" {
		c.JSON(http.StatusOK, UserRegisterResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  "密码不能为空",
			},
		})
		return
	}
	info, err := service.Register(c, username, password)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, UserRegisterResponse{
			Response: common.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	} else {
		c.JSON(http.StatusOK, UserRegisterResponse{
			Response: common.Response{StatusCode: 0},
			Token:    info.Token,
			UserId:   info.UserID,
		})
		return
	}
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	info, err := service.Login(c, username, password)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, UserRegisterResponse{
			Response: common.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	} else {
		c.JSON(http.StatusOK, UserRegisterResponse{
			Response: common.Response{StatusCode: 0},
			Token:    info.Token,
			UserId:   info.UserID,
		})
		return
	}
}

func UserInfo(c *gin.Context) {
	var err error
	var user *dao.User

	targetID, err := strconv.Atoi(c.Query("user_id"))
	userID, _ := c.MustGet("userid").(int64)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, UserResponse{
			Response: common.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	user2, err := service.GetUserIndex(c, userID, int64(targetID))
	user = utls.ChangeUser2ToUser(user2)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, UserResponse{
			Response: common.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}
	//返回成功
	c.JSON(http.StatusOK, UserResponse{
		Response: common.Response{StatusCode: 0, StatusMsg: "successful"},
		User:     *user,
	})
	return
}
