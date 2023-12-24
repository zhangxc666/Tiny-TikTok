package utls

import (
	"douyin/dao"
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

func CreateMessageContent(userID, targetID int64, content string) string {
	return fmt.Sprintf("%d-%d-%s", userID, targetID, content)
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
