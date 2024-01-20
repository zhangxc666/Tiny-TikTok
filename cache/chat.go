package cache

import (
	"context"
	"douyin/dao"
	"douyin/utls"
	"time"
)

// PushManyHistoryToRedis 把聊天记录推送至缓存中
func PushManyHistoryToRedis(key string, msgList []dao.Message) error {
	scoreList, contentList := []float64{}, []interface{}{}
	rc := MakeRdbCache()
	for _, msg := range msgList {
		scoreList = append(scoreList, float64(msg.CreateTime))
		contentList = append(contentList, utls.CreateMessageContent(msg.UserId, msg.ToUserId, msg.Content))
	}
	_, err := rc.ZAdd(context.Background(), key, scoreList, contentList)
	if err != nil {
		return err
	}
	_, err = rc.Expire(context.Background(), key, time.Hour*24)
	return err
}

// PushMessage 把聊天记录推送至缓存中
func PushMessage(key string, member string, score int64) error {
	rc := MakeRdbCache()
	_, err := rc.ZAdd(context.Background(), key, []float64{float64(score)}, []interface{}{member})
	return err
}

// GetMessageList 获取聊天记录
func GetMessageList(key, low, high string) (*[]dao.Message, error) {
	rc := MakeRdbCache()
	s, err := rc.ZGetRangeByScoreWithScores(context.Background(), key, low, high, 0, 0)
	if err != nil {
		return nil, err
	}
	return utls.ConvertToMessage(s), nil
}

// RemovePersistMessage 移除持久化缓存
func RemovePersistMessage(key string, msgList []dao.Message) error { //都是同一个key中的member
	members := make([]interface{}, len(msgList))
	if len(msgList) == 0 {
		return nil
	}
	for i, message := range msgList {
		members[i] = utls.CreateMessageContent(message.UserId, message.ToUserId, message.Content)
	}
	rc := MakeRdbCache()
	_, err := rc.ZRem(context.Background(), key, members)

	return err
}

// IsMessageKeyExist 判断messagekey是否存在
func IsMessageKeyExist(key string) (bool, error) {
	rc := MakeRdbCache()
	cnt, err := rc.Exists(context.Background(), key)
	return cnt != 0, err
}
