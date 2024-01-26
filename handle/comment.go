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

type CommentResponse struct {
	common.Response
	Comment dao.Comment
}
type CommentListsResponse struct {
	common.Response
	CommentLists []dao.Comment `json:"comment_list"`
}

func CommentAction(c *gin.Context) {

	var (
		Comment   *dao.Comment
		videoID   int64
		commentID int64
		err       error
	)
	userID := c.MustGet("userid").(int64)
	action := c.Query("action_type")
	commentText := c.Query("comment_text")
	//得到video_id
	videoIDStr := c.Query("video_id")
	if videoIDStr != "" {
		videoID, err = strconv.ParseInt(videoIDStr, 10, 64)
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusOK, CommentResponse{
				Response: common.Response{StatusCode: 1, StatusMsg: err.Error()},
			})
			return
		}
	}
	//得到评论id
	commentIDStr := c.Query("comment_id")
	if commentIDStr != "" {
		commentID, err = strconv.ParseInt(commentIDStr, 10, 64)
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusOK, CommentResponse{
				Response: common.Response{StatusCode: 1, StatusMsg: err.Error()},
			})
			return
		}
	}
	ok, err := service.CheckComment(action, commentText)
	if err != nil || ok == false {
		log.Println(err.Error())
		c.JSON(http.StatusOK, CommentResponse{
			Response: common.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
	}
	//进行方法判断
	Comment, err = service.CommentAction(c, userID, videoID, commentID, action, commentText)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, CommentResponse{
			Response: common.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}
	if Comment == nil {
		c.JSON(http.StatusOK, CommentResponse{
			Response: common.Response{StatusCode: 0, StatusMsg: "successful"},
		})
		return
	}

	c.JSON(http.StatusOK, CommentResponse{
		Response: common.Response{StatusCode: 0, StatusMsg: "successful"},
		Comment:  *Comment,
	})
	return
}

func CommentList(c *gin.Context) {
	var commentLists []dao.Comment

	videoID, err := strconv.Atoi(c.Query("video_id"))
	if err != nil {
		c.JSON(http.StatusOK, CommentListsResponse{
			Response: common.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}
	token := c.Query("token")
	userID := int64(-1)
	if token != "" {
		t, claim, err := utls.ParseToken(token)
		if t.Valid == true && err == nil {
			userID = claim.UserId
		}
	}
	commentLists, err = service.GetCommentList(c, userID, int64(videoID))
	if err != nil {
		c.JSON(http.StatusOK, CommentListsResponse{
			Response: common.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, CommentListsResponse{
		Response:     common.Response{StatusCode: 0, StatusMsg: "successful"},
		CommentLists: commentLists,
	})
	return
}
