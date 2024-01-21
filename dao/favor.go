package dao

import (
	"sync"
)

type Like struct {
	UserId  int64 `gorm:"column:user2_id"`
	VideoId int64 `gorm:"column:video_id"`
}

func (Like) TableName() string {
	return "like"
}

type LikeDao struct {
}

var likeDao *LikeDao

var likeOnce sync.Once

func GetLikeInstance() *LikeDao {
	likeOnce.Do(func() {
		likeDao = &LikeDao{}
	})
	return likeDao
}

// AddLike 添加映射
func (LikeDao) Update(like *Like) error {
	var count int64
	db.Model(&Like{}).Where("user2_id = ? AND video_id = ?", like.UserId, like.VideoId).Count(&count)
	var err error
	if count == 0 {
		// 如果记录不存在，则插入
		err = db.Create(&like).Error
	} else {
		// 如果记录已存在，则更新（根据需要调整此部分）
		err = db.Model(&Like{}).Where("user2_id = ? AND video_id = ?", like.UserId, like.VideoId).Updates(like).Error
	}
	return err
}

// DeleteLike 删除映射
func (LikeDao) DeleteLike(like *Like) error {
	//会删除所有符合条件得到对象 但是无所谓 都行
	err := db.Where("user2_id = ? and video_id = ?", like.UserId, like.VideoId).Delete(like).Error
	if err != nil {
		return err
	}
	//要对视频点赞数进行更新
	return nil
}

// QueryLikeByUserid DeleteLike 查找映射 并且返回lists
func (LikeDao) QueryLikeByUserid(userid int64) ([]Video, error) {
	user := &User{}
	err := db.Preload("VideoLieLists").Where("user2_id = ?", userid).Preload("VideoLieLists.Author").Find(user).Error
	if err != nil {
		return nil, err
	}
	return user.VideoLieLists, nil
}

// QueryVideoIDsByUserID 根据用户ID查询点赞的视频ID
func (LikeDao) QueryVideoIDsByUserID(userID int64) ([]int64, error) {
	videoIDs := []int64{}
	err := db.Model(&Like{}).Where("user2_id = ?", userID).Select("video_id").Find(&videoIDs).Error
	return videoIDs, err
}
