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

type PublishedResponse struct {
	common.Response
	VideoLists []dao.RetVideo `json:"video_list,omitempty"`
}
type FeedResponse struct {
	common.Response
	VideoLists []dao.RetVideo `json:"video_list,omitempty"`
	NextTime   int64          `json:"next_time,omitempty"`
}

func Feed(c *gin.Context) {
	token := c.Query("token")
	lastTime, _ := strconv.ParseInt(c.Query("latest_time"), 10, 64)
	//返回所有视频信息
	videoLists, timestamp, err := service.Feed(c, token, lastTime)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, FeedResponse{
			Response: common.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	} else {
		c.JSON(http.StatusOK, FeedResponse{
			Response:   common.Response{StatusCode: 0, StatusMsg: "successful"},
			VideoLists: utls.ChangeVideoToRetVideo(videoLists),
			NextTime:   timestamp,
		})
		return
	}
}

func PublishList(c *gin.Context) {
	//返回所有视频信息
	userid, _ := strconv.Atoi(c.Query("user_id"))

	videoLists, err := service.PublishList(c, int64(userid))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, FeedResponse{
			Response: common.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	} else {
		c.JSON(http.StatusOK, FeedResponse{
			Response:   common.Response{StatusCode: 0, StatusMsg: "successful"},
			VideoLists: utls.ChangeVideoToRetVideo(videoLists),
		})
		return
	}
}

func Publish(c *gin.Context) {
	//返回所有视频信息
	title := c.PostForm("title")
	//得到token 获取userid
	userid := c.MustGet("userid").(int64)
	//获取文件
	file, err := c.FormFile("data")

	if err != nil {
		//得到的文件错误
		log.Println(err.Error())
		c.JSON(http.StatusOK, FeedResponse{
			Response: common.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}
	if err := service.Publish(c, title, file, userid); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, FeedResponse{Response: common.Response{StatusCode: 1, StatusMsg: err.Error()}})
		return
	}
	c.JSON(http.StatusOK, FeedResponse{Response: common.Response{StatusCode: 0, StatusMsg: "success!"}})
	return
}
