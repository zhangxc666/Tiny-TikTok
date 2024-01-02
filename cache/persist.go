package cache

import (
	"context"
	"douyin/dao"
	"fmt"
	"strconv"
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

func PersistCountToDB() error {
	rc := MakeRdbCache()
	countKeys, err := rc.GetKeys(context.Background(), "user_count::*")
	if err != nil {
		return err
	}

	for _, key := range countKeys {
		count, err := GetUserCount(context.Background(), key)
		fmt.Println(count, key)
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
