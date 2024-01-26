package service

import (
	"context"
	"douyin/cache"
	"douyin/dao"
	"douyin/utls"
	"errors"
	"time"
)

// CommentAction 评论或删除评论操作
func CommentAction(ctx context.Context, userID, videoID, commentID int64, actionType, commentText string) (*dao.Comment, error) {
	switch actionType {
	case "1": // 创建评论
		comment := &dao.Comment{
			UserId:      userID,
			VideoId:     videoID,
			CommentText: commentText,
		}
		// 第一次删除缓存
		if err := cache.DelCommentZSet(ctx, videoID); err != nil {
			panic(err)
			return nil, err
		}

		// 添加到数据库中，延迟双删
		if err := dao.GetCommentInstance().CreateComment(comment); err != nil {
			panic(err)
			return nil, err
		}

		// 延迟
		time.Sleep(time.Millisecond * 50)

		// 第二次删除
		if err := cache.DelCommentZSet(ctx, videoID); err != nil {
			panic(err)
			return nil, err
		}
		// 修改comment—count缓存
		if err := cache.AddVideoCommentCount(ctx, videoID); err != nil {
			panic(err)
			return nil, err
		}
		var err error
		comment.User, err = GetUserIndex(ctx, userID, comment.UserId)
		if err != nil {
			panic(err)
			return nil, err
		}
		return comment, nil
	case "2": // 删除评论
		// 第一次删除缓存
		if err := cache.DelCommentZSet(ctx, videoID); err != nil {
			panic(err)
			return nil, err
		}

		// 操作数据库
		if err := dao.GetCommentInstance().DeleteCommentById(commentID); err != nil {
			panic(err)
			return nil, err
		}

		// 延迟
		time.Sleep(time.Millisecond * 50)

		// 第二次删除
		if err := cache.DelCommentZSet(ctx, videoID); err != nil {
			panic(err)
			return nil, err
		}

		if err := cache.SubVideoCommentCount(ctx, videoID); err != nil {
			panic(err)
			return nil, err
		}
		return nil, nil
	default:
		return nil, errors.New("未知操作")
	}
	return nil, nil
}

func GetCommentList(ctx context.Context, userID int64, videoID int64) ([]dao.Comment, error) {
	var commentList []dao.Comment
	if exist, err := cache.ExistVideoComment(ctx, videoID); err != nil {
		panic(err)
		return nil, err
	} else if exist == false {
		// 当前缓存查不到，查数据库
		commentList, err = dao.GetCommentInstance().QueryCommentByVideoId(videoID)
		if err != nil {
			panic(err)
			return nil, err
		}
		if len(commentList) == 0 {
			return nil, nil
		}
		timeStamps, commentIDs, contents, userIDs := make([]float64, len(commentList)), make([]int64, len(commentList)), make([]string, len(commentList)), make([]int64, len(commentList))
		for i := range commentList {
			comment := commentList[i]
			timeStamps[i] = float64(comment.CreatedAt.UnixMilli())
			commentIDs[i] = int64(comment.ID)
			contents[i] = comment.CommentText
			userIDs[i] = comment.UserId
		}
		// 存数据库
		if err := cache.AddComments(ctx, videoID, timeStamps, commentIDs, userIDs, contents); err != nil {
			panic(err)
			return nil, err
		}
	} else {
		// 缓存查到了
		var (
			userIDs []int64
			err     error
		)
		// 查缓存，获取commentList和userIDs
		commentList, userIDs, err = cache.GetCommentList(ctx, videoID)
		if err != nil {
			panic(err)
			return nil, err
		}
		// 查userList
		userList, err := GetUserList(ctx, userID, userIDs)
		if err != nil {
			panic(err)
			return nil, err
		}
		// 存commentList
		for i := range commentList {
			commentList[i].User = utls.ChangeUserToUser2(&userList[i])
		}
	}
	return commentList, nil
}
func CheckComment(action, commentText string) (bool, error) {
	if commentText == "" && action == "1" {
		return false, errors.New("评论字符串不能为空")
	}
	return true, nil
	// TODO 不良评论检测
}
