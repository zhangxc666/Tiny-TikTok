package utls

import (
	"douyin/dao"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
	"strings"
	"time"
)

func CreateChatHistoryKey(sendID, targetID int64) string {
	if sendID > targetID {
		sendID, targetID = targetID, sendID
	}
	return "chat::history::" + strconv.FormatInt(sendID, 10) + "-" + strconv.FormatInt(targetID, 10)
}

func CreateChatPersistKey(sendID, targetID int64) string {
	if sendID > targetID {
		sendID, targetID = targetID, sendID
	}
	return "chat::message::" + strconv.FormatInt(sendID, 10) + "-" + strconv.FormatInt(targetID, 10)
}

func CreateFollowKey(userID int64) string {
	return "follow::" + strconv.FormatInt(userID, 10)
}

func CreateFanKey(userID int64) string {
	return "fan::" + strconv.FormatInt(userID, 10)
}

func CreateFriendKey(userID int64) string {
	return "friend::" + strconv.FormatInt(userID, 10)
}

func CreatePublishKey() string {
	return "video_publish"
}

func CreateUserVideoKey(userID int64) string {
	return "user_video::" + strconv.Itoa(int(userID))
}

func CreateUserInfoKey(userID int64) string {
	return "user_info::" + strconv.FormatInt(userID, 10)
}

func CreateUserCountKey(userID int64) string {
	return "user_count::" + strconv.FormatInt(userID, 10)
}

func CreateVideoKey(videoID int64) string { return "video_info::" + strconv.Itoa(int(videoID)) }

func CreateFavorKey(userID int64) string {
	return "user_favorite::" + strconv.FormatInt(userID, 10)
}

func CreateFavorMember(videoID int64, actionType int) string {
	return strconv.Itoa(actionType) + "+" + strconv.FormatInt(videoID, 10)
}

func CreateCommentKey(videoID int64) string {
	return "comment::" + strconv.FormatInt(videoID, 10)
}

func CreateCommentMember(commendID int64, userID int64, content string) string {
	return strconv.FormatInt(commendID, 10) + "-" + strconv.FormatInt(userID, 10) + "-" + content
}

func CreateMapUserInfo(userInfo *dao.User2) map[string]interface{} {
	userStr, _ := json.Marshal(userInfo)
	userMap := make(map[string]interface{})
	_ = json.Unmarshal(userStr, &userMap)
	delete(userMap, "Usercount")
	fmt.Println("userInfo: ", userMap)
	return userMap
}

func CreateMapVideoInfo(videoInfo *dao.Video) map[string]interface{} {
	videoStr, _ := json.Marshal(videoInfo)
	videoMap := make(map[string]interface{})
	_ = json.Unmarshal(videoStr, &videoMap)
	delete(videoMap, "CreatedAt")
	delete(videoMap, "DeletedAt")
	delete(videoMap, "UpdatedAt")
	delete(videoMap, "author")
	delete(videoMap, "is_favorite")
	fmt.Println("userMap:", videoMap)
	return videoMap
}

func CreateMapUserCount(userCount *dao.UserCount) map[string]interface{} {
	userStr, _ := json.Marshal(userCount)
	userMap := make(map[string]interface{})
	_ = json.Unmarshal(userStr, &userMap)
	delete(userMap, "CreatedAt")
	delete(userMap, "DeletedAt")
	delete(userMap, "UpdatedAt")
	delete(userMap, "ID")
	fmt.Println("userCount: ", userMap)
	return userMap
}

func CreateMessageContent(userID, targetID int64, content string) string {
	return fmt.Sprintf("%d-%d-%s", userID, targetID, content)
}

func CreateVideoInfo(videoMap map[string]string) (*dao.Video, error) {
	videoStr, _ := json.Marshal(videoMap)
	videoInfo := new(dao.Video)
	err := json.Unmarshal(videoStr, videoInfo)
	if err != nil {
		return nil, err
	}
	return videoInfo, nil
}
func CreateUserInfo(userMap map[string]string) (*dao.User2, error) {
	userStr, _ := json.Marshal(userMap)
	userInfo := new(dao.User2)
	err := json.Unmarshal(userStr, userInfo)
	if err != nil {
		return nil, err
	}
	return userInfo, err
}

func CreateUserCount(userMap map[string]string) (*dao.UserCount, error) {
	userStr, _ := json.Marshal(userMap)
	userCount := new(dao.UserCount)
	err := json.Unmarshal(userStr, userCount)
	if err != nil {
		return nil, err
	}
	return userCount, err
}
func StringToMessage(member, score string) dao.Message {
	str := strings.SplitN(member, "-", 3)
	userid, _ := strconv.ParseInt(str[0], 10, 64)
	targetid, _ := strconv.ParseInt(str[1], 10, 64)
	scoreNum, _ := strconv.ParseInt(score, 10, 64)
	content := str[2]
	return dao.Message{
		UserId:     userid,
		ToUserId:   targetid,
		Content:    content,
		CreateTime: scoreNum,
	}
}
func ConvertToMessage(z []redis.Z) *[]dao.Message {
	messageList := make([]dao.Message, len(z))
	for i := range messageList {
		m, f := z[i].Member, z[i].Score
		messageList[i] = StringToMessage(m.(string), strconv.FormatFloat(f, 'f', -1, 64))
	}
	return &messageList
}
func ConvertAllStringToMessage(messages []string) *[]dao.Message {
	MessageList := make([]dao.Message, len(messages)>>1)
	for i := 0; i < len(messages); i += 2 {
		MessageList[(i+1)>>1] = StringToMessage(messages[i], messages[i+1])
	}
	return &MessageList
}

func ExecuteTimedTask(interval time.Duration, f func() error) {
	timer := time.NewTimer(interval)
	for {
		timer.Reset(interval) // 这里复用了 timer
		select {
		case <-timer.C:
			err := f()

			if err != nil {
				panic(err)
				log.Println(err)
				return
			}
		}
	}
}

func ChangeUser2ToUser(user *dao.User2) *dao.User {
	if user.Usercount == nil {
		return &dao.User{
			ID:              user.ID,
			Name:            user.Name,
			IsFollow:        user.IsFollow,
			Avatar:          user.Avatar,
			BackGroundImage: user.BackGroundImage,
			Signature:       user.Signature,
			VideoLieLists:   user.VideoLieLists,
		}
	}
	return &dao.User{
		ID:              user.ID,
		Name:            user.Name,
		FollowCount:     user.Usercount.FollowCount,
		FollowerCount:   user.Usercount.FollowerCount,
		IsFollow:        user.IsFollow,
		Avatar:          user.Avatar,
		BackGroundImage: user.BackGroundImage,
		Signature:       user.Signature,
		TotalFavorite:   user.Usercount.TotalFavorited,
		WorkCount:       user.Usercount.WorkCount,
		FavoriteCount:   user.Usercount.FavoriteCount,
		VideoLieLists:   user.VideoLieLists,
	}
}

func ChangeUserToUser2(user *dao.User) *dao.User2 {
	return &dao.User2{
		ID:              user.ID,
		Name:            user.Name,
		IsFollow:        user.IsFollow,
		Avatar:          user.Avatar,
		BackGroundImage: user.BackGroundImage,
		Signature:       user.Signature,
		VideoLieLists:   user.VideoLieLists,
		Usercount: &dao.UserCount{
			UserID:         user.ID,
			FollowCount:    user.FollowCount,
			FollowerCount:  user.FollowerCount,
			TotalFavorited: user.TotalFavorite,
			WorkCount:      user.WorkCount,
			FavoriteCount:  user.FavoriteCount,
		},
	}
}

func ChangeVideoToRetVideo(videoList []dao.Video) []dao.RetVideo {
	retVideoList := make([]dao.RetVideo, len(videoList))
	for i, v := range videoList {
		retVideoList[i].Author = *ChangeUser2ToUser(&v.Author)
		retVideoList[i].ID = v.ID
		retVideoList[i].UserId = v.UserId
		retVideoList[i].PlayUrl = v.PlayUrl
		retVideoList[i].CoverUrl = v.CoverUrl
		retVideoList[i].FavoriteCount = v.FavoriteCount
		retVideoList[i].CommentCount = v.CommentCount
		retVideoList[i].IsFavorite = v.IsFavorite
		retVideoList[i].TimeStamp = v.TimeStamp
		retVideoList[i].Title = v.Title
	}
	return retVideoList
}
