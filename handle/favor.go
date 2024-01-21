package handle

import (
	"douyin/common"
	"douyin/dao"
	"douyin/service"
	"douyin/utls"
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type FavoriteResponse struct {
	common.Response
}
type FavoriteListsResponse struct {
	common.Response
	VideoLists []dao.RetVideo `json:"video_list,omitempty"`
}

func FavoriteAction(c *gin.Context) {
	//解析得到id
	userid := c.MustGet("userid").(int64)
	action, _ := strconv.Atoi(c.Query("action_type"))
	videoId, err := strconv.Atoi(c.Query("video_id"))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, FavoriteResponse{
			Response: common.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}
	err = service.FavoriteAction(c, userid, int64(videoId), int(action))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, FavoriteResponse{
			Response: common.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, FeedResponse{
		Response: common.Response{StatusCode: 0, StatusMsg: "successful"},
	})
	return
}

func FavoriteList(c *gin.Context) {
	//已经鉴别完token
	userid, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, FavoriteListsResponse{
			Response: common.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}
	//得到喜欢的列表
	videoLists, err := service.FavoriteList(c, int64(userid))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, FavoriteListsResponse{
			Response: common.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, FavoriteListsResponse{
		Response:   common.Response{StatusCode: 0, StatusMsg: "successful"},
		VideoLists: utls.ChangeVideoToRetVideo(videoLists),
	})
	return
}
