package handle

import (
	"douyin/common"
	"douyin/dao"
	"douyin/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type MessageResponse struct {
	common.Response
	MessageLists []dao.Message `json:"message_list"`
}

func ChatAction(c *gin.Context) {
	// 通过中间件验证传userid，不会不合法
	userID, _ := c.MustGet("userid").(int64)
	targetID, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{StatusCode: 0, StatusMsg: err.Error()})
	}
	_ = c.Query("action_type")
	content := c.Query("content")
	// 双方不是朋友
	err = service.SendMessage(userID, targetID, content)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, common.Response{StatusCode: 0, StatusMsg: "successful"})
	return
}

func GetMessageList(c *gin.Context) {
	userid, _ := c.MustGet("userid").(int64)
	targetID, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, MessageResponse{Response: common.Response{StatusCode: 1, StatusMsg: err.Error()}})
		return
	}
	preMsgTime, err := strconv.ParseInt(c.Query("pre_msg_time"), 10, 64)
	preMsgTime += 1
	if err != nil {
		c.JSON(http.StatusOK, MessageResponse{Response: common.Response{StatusCode: 1, StatusMsg: err.Error()}})
		return
	}
	msgList, err := service.GetMessage(userid, targetID, preMsgTime)
	if err != nil {
		c.JSON(http.StatusOK, MessageResponse{Response: common.Response{StatusCode: 1, StatusMsg: err.Error()}})
		return
	}
	c.JSON(http.StatusOK, MessageResponse{Response: common.Response{StatusCode: 0, StatusMsg: "successful"}, MessageLists: *msgList})

}
