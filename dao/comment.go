package dao

import (
	"encoding/json"
	"gorm.io/gorm"
	"sync"
)

type Comment struct {
	gorm.Model
	User        *User2 `gorm:"foreignKey:UserId" json:"user"`
	UserId      int64  `gorm:"column:user_id"    json:"-"`
	VideoId     int64  `gorm:"column:video_id"   json:"-"`
	CommentText string `gorm:"column:comment_text"   json:"content,omitempty"`
}

func (c *Comment) MarshalJSON() ([]byte, error) {
	type RetComment struct {
		ID          int64  `json:"comment_id"`
		ContentText string `json:"content"`
		CreateDate  string `json:"create_date"`
		*User       `json:"user"`
	}
	return json.Marshal(&RetComment{
		ID:          int64(c.ID),
		ContentText: c.CommentText,
		CreateDate:  c.CreatedAt.Format("01-02"), // 评论发布日期，格式 mm-dd
		User:        changeUser2ToUser(c.User),
	})
}

func changeUser2ToUser(user *User2) *User {
	if user.Usercount == nil {
		return &User{
			ID:              user.ID,
			Name:            user.Name,
			IsFollow:        user.IsFollow,
			Avatar:          user.Avatar,
			BackGroundImage: user.BackGroundImage,
			Signature:       user.Signature,
			VideoLieLists:   user.VideoLieLists,
		}
	}
	return &User{
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

func (Comment) TableName() string {
	return "comment"
}

type CommentDao struct {
}

var commentDao *CommentDao

var commentOnce sync.Once

func GetCommentInstance() *CommentDao {
	commentOnce.Do(func() {
		commentDao = &CommentDao{}
	})
	return commentDao
}

func (CommentDao) CreateComment(comment *Comment) error {
	err := db.Create(comment).Error
	if err != nil {
		return err
	}
	return nil
}
func (CommentDao) DeleteCommentById(commentId int64) error {
	//根据主键删除评论
	res := db.Model(&Comment{}).Where("ID = ?", commentId).Delete(&Comment{})
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (CommentDao) QueryCommentByVideoId(videoId int64) ([]Comment, error) {
	var commentLists []Comment
	//按时间的倒叙排序
	err := db.Model(&Comment{}).Preload("User").Where("video_id =?", videoId).Order("created_at desc").Find(&commentLists).Error
	if err != nil {
		return nil, err
	}
	return commentLists, nil
}
