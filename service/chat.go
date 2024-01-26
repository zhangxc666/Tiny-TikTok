package service

import (
	"douyin/cache"
	"douyin/dao"
	"douyin/utls"
	"errors"
	"strconv"
	"time"
)

func SendMessage(userID, targetID int64, content string) error {
	if dao.GetFollowInstance().IsBothFollow(userID, targetID) == false {
		return errors.New("双方不是好友")
	}
	historyKey, persistKey := utls.CreateChatHistoryKey(userID, targetID), utls.CreateChatPersistKey(userID, targetID)
	contentMsg := utls.CreateMessageContent(userID, targetID, content)
	timeStamp := time.Now().UnixNano() / int64(time.Millisecond)
	// 先推到历史聊天记录缓存
	err := cache.PushMessage(historyKey, contentMsg, timeStamp)
	if err != nil {
		return err
	}
	// 再推到持久化数据库缓存
	err = cache.PushMessage(persistKey, contentMsg, timeStamp)
	if err != nil {
		return err
	}
	return nil
}

func GetMessage(userID, targetID, preMsgTime int64) (*[]dao.Message, error) {
	messageKey := utls.CreateChatHistoryKey(userID, targetID)
	// 先查redis
	exist, err := cache.IsMessageKeyExist(messageKey)
	if err != nil {
		return nil, err
	}
	// 当前key存在，表示聊天记录存在
	var MessageList *[]dao.Message
	if exist {
		MessageList, err = cache.GetMessageList(messageKey, strconv.FormatInt(preMsgTime, 10), "inf")
		if err != nil {
			return nil, err
		}
		if preMsgTime > 1 {
			temp := []dao.Message{}
			for _, message := range *MessageList {
				if message.UserId == userID {
					continue
				}
				temp = append(temp, message)
			}
			MessageList = &temp
		}
	} else { // 如果不存在再查mysql
		MessageList, err = dao.GetMessageInstance().QueryMessageLists(userID, targetID, preMsgTime)
		if err != nil {
			return nil, err
		}
		// 再把message推到redis
		err = cache.PushManyHistoryToRedis(messageKey, *MessageList)
		if err != nil {
			return MessageList, err
		}
	}
	return MessageList, nil
}
