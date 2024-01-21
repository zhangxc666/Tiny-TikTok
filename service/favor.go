package service

import (
	"context"
	"douyin/cache"
	"douyin/dao"
	"douyin/utls"
	"errors"
)

func FavoriteAction(ctx context.Context, userID int64, videoID int64, actionType int) error {
	favorKey := utls.CreateFavorKey(userID)
	exist, err := cache.ExistFavorRecord(ctx, favorKey, videoID)
	if err != nil {
		return err
	}
	videoInfo, err := GetVideoInfo(ctx, userID, videoID)
	// 保证了videoInfo和UserInfo强制在缓存中
	if err != nil {
		return err
	}
	switch actionType {
	case 1:
		if exist {
			return errors.New("已点赞过，不可重复点赞")
		}
		// 添加点赞关系至redis中
		if err := cache.AddFavorRecord(ctx, favorKey, videoID); err != nil {
			return err
		}
		// 修改videoInfo的favorCount
		if err := AddVideoFavorCount(ctx, videoID); err != nil {
			return err
		}
		// 增加用户的喜欢数
		if err := AddUserFavorCount(ctx, userID, videoInfo.UserId); err != nil {
			return err
		}
	case 2:
		if !exist {
			return errors.New("未点赞，不可取消")
		}
		if err := cache.AddCancelFavorRecord(ctx, favorKey, videoID); err != nil {
			return err
		}
		if err := SubVideoFavorCount(ctx, videoID); err != nil {
			return err
		}
		if err := SubUserFavorCount(ctx, userID, videoInfo.UserId); err != nil {
			return err
		}
	default:
		return errors.New("未知操作")
	}
	return nil
}

func FavoriteList(ctx context.Context, userID int64) ([]dao.Video, error) {
	videoIDs, err := cache.GetAllFavoriteVideosIDs(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(videoIDs) == 0 {
		videoIDs, err = dao.GetLikeInstance().QueryVideoIDsByUserID(userID)
		if err != nil {
			return nil, err
		}
		if len(videoIDs) == 0 {
			return nil, nil
		}
		if err = cache.SetUserFavorVideoIDs(ctx, utls.CreateFavorKey(userID), videoIDs); err != nil {
			return nil, err
		}
	}
	videoInfos, err := GetManyVideoInfos(ctx, videoIDs, userID)
	if err != nil {
		return nil, err
	}
	return videoInfos, nil
}

func IsManyFavorVideos(ctx context.Context, userID int64, videoIDs []int64) ([]bool, error) {
	isFavor := make([]bool, len(videoIDs))
	for i, videoID := range videoIDs {
		exist, err := IsFavorVideo(ctx, userID, videoID)
		if err != nil {
			return nil, err
		}
		isFavor[i] = exist
	}
	return isFavor, nil
}

func IsFavorVideo(ctx context.Context, userID, videoID int64) (bool, error) {
	favorKey := utls.CreateFavorKey(userID)
	return cache.ExistFavorRecord(ctx, favorKey, videoID)
}
