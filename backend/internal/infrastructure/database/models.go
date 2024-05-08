package database

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

var allModels = []any{
	new(User),
	new(Tweet),
	new(Media),
	new(Like),
	new(Session),
	new(Trend),
}

type User struct {
	ID          uint64 `gorm:"primaryKey"`
	Username    string
	Password    []byte
	Name        string
	Website     string
	Description string
	Pic         string
	PicCover    string
	Followings  []*User `gorm:"many2many:user_relation;foreignKey:ID;joinForeignKey:Following;References:ID;joinReferences:Follower"`
	Followers   []*User `gorm:"many2many:user_relation;foreignKey:ID;joinForeignKey:Follower;References:ID;joinReferences:Following"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type Session struct {
	SessionID string `gorm:"primaryKey"`
	UserID    uint64
	User      User `gorm:"foreignKey:UserID;references:ID;"`
	ExpiredAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Tweet struct {
	ID         uint64 `gorm:"primaryKey"`
	AuthorID   uint64
	Author     *User `gorm:"foreignKey:AuthorID;references:ID;"`
	TweetData  string
	TweetMedia []Media
	Likes      []Like
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

// Cascade Media and Likes soft delete issue
func (t *Tweet) AfterDelete(tx *gorm.DB) error {
	if err := tx.First(&Tweet{ID: t.ID}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err = tx.Where("tweet_id = ?", t.ID).Delete(&Media{}).Error; err != nil {
				return err
			}
			if err = tx.Where("tweet_id = ?", t.ID).Delete(&Like{}).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

type Media struct {
	ID        uint64 `gorm:"primaryKey"`
	TweetID   *uint64
	Link      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Like struct {
	ID        uint64 `gorm:"primaryKey"`
	TweetID   uint64
	Tweet     Tweet `gorm:"foreignKey:TweetID;references:ID;"`
	UserID    uint64
	User      User `gorm:"foreignKey:UserID;references:ID;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Trend struct {
	ID         uint64 `gorm:"primaryKey"`
	Name       string
	TweetCount uint64
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}
