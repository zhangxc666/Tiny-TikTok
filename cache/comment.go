package cache

import (
	"context"
	"douyin/dao"
	"douyin/utls"
	"strconv"
	"strings"
	"time"
)

func DelCommentZSet(ctx context.Context, videoID int64) error {
	rc := MakeRdbCache()
	commentKey := utls.CreateCommentKey(videoID)
	err := rc.Del(ctx, commentKey)
	return err
}

func ExistVideoComment(ctx context.Context, videoID int64) (bool, error) {
	rc := MakeRdbCache()
	commentKey := utls.CreateCommentKey(videoID)
	exist, err := rc.Exists(ctx, commentKey)
	return exist == 1, err
}

func AddComments(ctx context.Context, videoID int64, timeStamps []float64, commentIDs []int64, userIDs []int64, contents []string) error {
	rc := MakeRdbCache()
	commentMembers := make([]any, len(timeStamps))
	commentKey := utls.CreateCommentKey(videoID)
	for i, _ := range timeStamps {
		commentMembers[i] = utls.CreateCommentMember(commentIDs[i], userIDs[i], contents[i])
	}
	_, err := rc.ZAdd(ctx, commentKey, timeStamps, commentMembers)
	return err
}

func GetCommentList(ctx context.Context, videoID int64) ([]dao.Comment, []int64, error) {
	rc := MakeRdbCache()
	commentKey := utls.CreateCommentKey(videoID)
	redisZ, err := rc.ZGetRevRangeByScoreWithScores(ctx, commentKey, "-inf", "+inf", 0, -1)
	if err != nil {
		return nil, nil, err
	}
	userIDs := make([]int64, len(redisZ))
	commentList := make([]dao.Comment, len(redisZ))
	for i := range commentList {
		timestamp := int64(redisZ[i].Score)
		seconds := timestamp / 1000
		nanoseconds := (timestamp % 1000) * 1000000
		commentList[i].CreatedAt = time.Unix(seconds, nanoseconds)
		member, _ := redisZ[i].Member.(string)
		strs := strings.SplitN(member, "-", 3)
		commentID, _ := strconv.ParseInt(strs[0], 10, 64)
		userID, _ := strconv.ParseInt(strs[1], 10, 64)
		comment := strs[2]
		commentList[i].ID = uint(commentID)
		commentList[i].UserId = userID
		commentList[i].CommentText = comment
		userIDs[i] = userID
	}
	return commentList, userIDs, nil
}
