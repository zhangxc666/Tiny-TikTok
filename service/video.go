package service

import (
	"context"
	"douyin/cache"
	"douyin/config"
	"douyin/dao"
	"douyin/utls"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"path/filepath"
	"time"
)

func Feed(c *gin.Context, token string, lastTime int64) ([]dao.Video, int64, error) {
	feedKey := utls.CreatePublishKey()
	nextTime, videoIDs, err := cache.GetPublishVideoIDs(c, feedKey, lastTime)
	if err != nil {
		return nil, -1, err
	}
	if len(videoIDs) == 0 {
		var timeStamps []int64
		videoIDs, timeStamps, err = dao.GetVideoInstance().QueryAllPublishVideoID()
		if err != nil {
			return nil, -1, nil
		}
		if len(videoIDs) == 0 {
			return nil, -1, errors.New("没有视频")
		}
		err = cache.SetPublishVideoIDs(c, feedKey, timeStamps, videoIDs)
		if err != nil {
			return nil, -1, err
		}
	}
	var videoInfos []dao.Video
	// 判断token是否有效
	userID := int64(-1)
	if token != "" {
		t, claim, err := utls.ParseToken(token)
		if t.Valid == true && err == nil {
			userID = claim.UserId
		}
	}
	videoInfos, err = GetManyVideoInfos(c, videoIDs, userID)
	if err != nil {
		return nil, nextTime, err
	}
	return videoInfos, nextTime, nil
}

func Publish(c *gin.Context, title string, data *multipart.FileHeader, userid int64) error {
	fileName := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%d_%s", userid, fileName)
	saveVideoFilePath := filepath.Join("./public/", finalName)
	savePageFilePath := "./public/" + finalName
	if err := c.SaveUploadedFile(data, saveVideoFilePath); err != nil {
		return err
	}
	if err := utls.GenerateSnapshot(saveVideoFilePath, savePageFilePath, 3); err != nil {
		return err
	}
	videoUrl := "http://" + config.C.Resouece.Ipaddress + ":" + config.C.Resouece.Port + "/" + "static/" + finalName
	coverUrl := "http://" + config.C.Resouece.Ipaddress + ":" + config.C.Resouece.Port + "/" + "static/" + finalName + ".png"
	// 添加视频至数据库中
	videoInfo := &dao.Video{
		UserId:        userid,
		PlayUrl:       videoUrl,
		CoverUrl:      coverUrl,
		FavoriteCount: 0,
		CommentCount:  0,
		Title:         title,
		TimeStamp:     time.Now().UnixNano() / int64(time.Millisecond),
	}
	if err := dao.GetVideoInstance().AddVideo(videoInfo); err != nil {
		return err
	}
	// 将videoInfo存入redis中，key:videoID，value:videoInfo
	videoKey := utls.CreateVideoKey(videoInfo.ID)
	videoMap := utls.CreateMapVideoInfo(videoInfo)
	if err := cache.SetVideoInfo(c, videoKey, videoMap); err != nil {
		return err
	}
	// userid的workCount+1
	if err := AddWorkCount(c, userid); err != nil {
		return err
	}
	// 添加新发布的视频至ZSET中
	if err := cache.AddPublishVideo(c, utls.CreatePublishKey(), float64(videoInfo.TimeStamp), videoInfo.ID); err != nil {
		return err
	}
	// 删除userid的发布视频列表缓存
	if err := cache.DelUserVideo(c, utls.CreateUserVideoKey(userid)); err != nil {
		return err
	}
	return nil
}

func PublishList(ctx context.Context, userID int64) ([]dao.Video, error) {
	// 获取user的信息
	userIndex, err := GetUserIndex(ctx, userID, userID)
	if err != nil {
		return nil, err
	}
	userVideoKey := utls.CreateUserVideoKey(userID)
	// 获取user发布的所有video的IDs
	videoIDs, err := cache.GetUserVideoIDs(ctx, userVideoKey)
	if err != nil {
		fmt.Println("videoIDs1: ", videoIDs)
		return nil, err
	}
	if videoIDs == nil {
		videoIDs, err = dao.GetVideoInstance().QueryVideoIDsByUserId(userID)
		fmt.Println("videoIDs2: ", videoIDs)
		if err != nil {
			return nil, err
		}
		if len(videoIDs) == 0 {
			return nil, nil
		}
		if err = cache.SetUserVideoIDs(ctx, userVideoKey, videoIDs); err != nil {
			return nil, err
		}
	}
	// 根据video的IDS获取video的信息
	isFavorList, err := IsManyFavorVideos(ctx, userID, videoIDs)
	if err != nil {
		return nil, err
	}
	videoList := make([]dao.Video, len(videoIDs))
	for i := range videoList {
		videoInfo, err := GetVideoInfoSelf(ctx, videoIDs[i])
		if err != nil {
			return nil, err
		}
		videoInfo.Author = *userIndex
		videoList[i] = *videoInfo
		videoList[i].ID = videoIDs[i]
		videoList[i].IsFavorite = isFavorList[i]
	}
	return videoList, nil
}

// GetVideoInfoSelf 根据videoid获取视频信息（不包含user相关信息）,从redis查，查不到从数据库查）
func GetVideoInfoSelf(ctx context.Context, videoID int64) (*dao.Video, error) {
	videoKey := utls.CreateVideoKey(videoID)
	videoInfo, err := cache.GetVideoInfo(ctx, videoKey)
	if err != nil {
		return nil, err
	}
	if videoInfo == nil {
		videoInfo, err = dao.GetVideoInstance().QueryVideoInfoByVideoID(videoID)
		if err != nil {
			return nil, err
		}
		if videoInfo == nil {
			return nil, nil
		}
		if err := cache.SetVideoInfo(ctx, videoKey, utls.CreateMapVideoInfo(videoInfo)); err != nil {
			return nil, err
		}
	}
	return videoInfo, err
}

// GetVideoInfo 获取videoInfo，包括作者信息
func GetVideoInfo(ctx context.Context, videoID, userID int64) (*dao.Video, error) {
	videoInfo, err := GetVideoInfoSelf(ctx, videoID)
	if err != nil {
		return nil, err
	}
	// 获取作者信息
	userIndex, err := GetUserIndex(ctx, userID, videoInfo.UserId)
	if err != nil {
		return nil, err
	}
	videoInfo.Author = *userIndex
	videoInfo.IsFavorite, err = IsFavorVideo(ctx, userID, videoID)
	if err != nil {
		return nil, err
	}
	return videoInfo, nil
}

// GetManyVideoInfos 批量获取VideoInfo
func GetManyVideoInfos(ctx context.Context, videoIDs []int64, userID int64) ([]dao.Video, error) {
	var (
		isFavor []bool
		err     error
	)
	if userID != -1 {
		isFavor, err = IsManyFavorVideos(ctx, userID, videoIDs)
		if err != nil {
			return nil, err
		}
	}
	videoInfos := make([]dao.Video, len(videoIDs))
	authorIDs := make([]int64, len(videoIDs))
	for i := range videoInfos {
		videoID := videoIDs[i]
		videoInfo, err := GetVideoInfoSelf(ctx, videoID)
		if err != nil {
			return nil, err
		}
		videoInfos[i] = *videoInfo
		authorIDs[i] = videoInfo.UserId
		if userID != -1 {
			videoInfos[i].IsFavorite = isFavor[i]
		}
	}
	authorInfos, err := GetUserList(ctx, userID, authorIDs)
	if err != nil {
		return nil, err
	}
	for i := range videoIDs {
		videoInfos[i].Author = *(utls.ChangeUserToUser2(&authorInfos[i]))
		fmt.Printf("%+v\n", videoInfos[i])
	}

	return videoInfos, err
}

// AddVideoFavorCount 添加视频的喜欢数
func AddVideoFavorCount(ctx context.Context, videoID int64) error {
	if err := cache.AddVideoFavorCount(ctx, videoID); err != nil {
		return err
	}
	return nil
}

// SubVideoFavorCount 减少视频的喜欢数
func SubVideoFavorCount(ctx context.Context, videoID int64) error {
	if err := cache.SubVideoFavorCount(ctx, videoID); err != nil {
		return err
	}
	return nil
}
