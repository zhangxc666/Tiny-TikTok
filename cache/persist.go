package cache

import (
	"context"
	"douyin/dao"
	"douyin/utls"
	"strconv"
	"strings"
)

// PersistHistoryToDB 持久化聊天记录至DB中
func PersistHistoryToDB() error {
	rc := MakeRdbCache()
	messageKeys, err := rc.GetKeys(context.Background(), "chat::message::*")
	if err != nil {
		return err
	}
	for _, key := range messageKeys {
		messageList, err := GetMessageList(key, strconv.Itoa(0), "inf")
		if err != nil {
			return err
		}
		for _, message := range *messageList {
			err := dao.GetMessageInstance().AddMessage(&message)
			if err != nil {
				return err
			}
		}
		err = RemovePersistMessage(key, *messageList)
		if err != nil {
			return err
		}
	}
	return nil
}

// PersistUserCountToDB 持久化userCount至DB中
func PersistUserCountToDB() error {
	rc := MakeRdbCache()
	countKeys, err := rc.GetKeys(context.Background(), "user_count::*")
	if err != nil {
		return err
	}

	for _, key := range countKeys {
		count, err := GetUserCount(context.Background(), key)
		if err != nil {
			return err
		}
		err = dao.GetUser2Instance().UpdateCount(count)
		if err != nil {
			return err
		}
	}
	return nil
}

// PersistVideoInfoToDB 持久化video信息至DB中
func PersistVideoInfoToDB() error {
	rc := MakeRdbCache()
	videoInfoKeys, err := rc.GetKeys(context.Background(), "video_info::*")
	if err != nil {
		return err
	}
	for _, key := range videoInfoKeys {
		videoMap, err := rc.HGetAll(context.Background(), key)
		if err != nil {
			return err
		}
		videoInfo, err := utls.CreateVideoInfo(videoMap)
		if err != nil {
			return err
		}
		err = dao.GetVideoInstance().UpdateVideoInfo(videoInfo)
		if err != nil {
			return err
		}
	}
	return nil
}

// PersistFavorToDB 持久化点赞记录至DB中
func PersistFavorToDB() error {
	rc := MakeRdbCache()
	favorKeys, err := rc.GetKeys(context.Background(), "user_favorite::*")
	if err != nil {
		return err
	}
	for _, key := range favorKeys {
		userIdStr := strings.SplitN(key, "::", 2)[1]
		userID, _ := strconv.ParseInt(userIdStr, 10, 64)
		strList, err := rc.SGetAll(context.Background(), key)
		if err != nil {
			return err
		}
		for _, str := range strList {
			favorStr := strings.Split(str, "+")
			videoID, _ := strconv.ParseInt(favorStr[1], 10, 64)
			var err error
			if favorStr[0] == "1" {
				err = dao.GetLikeInstance().Update(&dao.Like{
					userID,
					videoID,
				})
			} else {
				err = dao.GetLikeInstance().DeleteLike(&dao.Like{
					userID,
					videoID,
				})
				if err != nil {
					return err
				}
				_, err = rc.SRem(context.Background(), key, str)
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}
