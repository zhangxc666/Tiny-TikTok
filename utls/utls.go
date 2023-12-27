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

func CreateUserInfoKey(userID int64) string {
	return "user_info::" + strconv.FormatInt(userID, 10)
}

func CreateUserCountKey(userID int64) string {
	return "user_count::" + strconv.FormatInt(userID, 10)
}

func CreateMapUserInfo(userInfo *dao.User2) map[string]interface{} {
	userStr, _ := json.Marshal(userInfo)
	userMap := make(map[string]interface{})
	_ = json.Unmarshal(userStr, &userMap)
	delete(userMap, "Usercount")
	fmt.Println("userInfo: ", userMap)
	return userMap
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
				log.Println(err)
				return
			}
		}
	}
}

func ChangeUser2ToUser(user *dao.User2) *dao.User {
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
	}
}
