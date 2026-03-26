package models

import (
	"database/sql"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// NullTime custom type for nullable time with JSON support
type NullTime struct {
	sql.NullTime
}

// MarshalJSON implements json.Marshaler
func (nt NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(nt.Time.Format("2006-01-02"))
}

// UnmarshalJSON implements json.Unmarshaler
func (nt *NullTime) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == nil || *s == "" {
		nt.Valid = false
		return nil
	}
	t, err := time.Parse("2006-01-02", *s)
	if err != nil {
		return err
	}
	nt.Time = t
	nt.Valid = true
	return nil
}

// User 用户模型
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"size:50;unique;not null"`
	Password  string         `json:"-" gorm:"size:100;not null"`
	Role      string         `json:"role" gorm:"size:20;not null"` // super, user
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Product 产品模型
type Product struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:255;not null"`
	Description string         `json:"description" gorm:"type:text"`
	Category    string         `json:"category" gorm:"size:100;not null"`
	Standard    string         `json:"standard" gorm:"size:100"`
	Material    string         `json:"material" gorm:"size:100"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Images      []ProductImage `json:"images" gorm:"foreignKey:ProductID"`
}

// ProductImage 产品图片模型
type ProductImage struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	ProductID uint           `json:"product_id" gorm:"not null"`
	ImageURL  string         `json:"image_url" gorm:"size:255;not null"`
	Order     int            `json:"order" gorm:"default:0"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Feedback 反馈模型
type Feedback struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"size:255;not null"`
	Email     string         `json:"email" gorm:"size:255;not null"`
	Phone     string         `json:"phone" gorm:"size:100"`
	Company   string         `json:"company" gorm:"size:255"`
	Product   string         `json:"product" gorm:"size:255"`
	Message   string         `json:"message" gorm:"type:text;not null"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// News 新闻模型
type News struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"size:255;not null"`
	CoverImage  string         `json:"cover_image" gorm:"size:255"`
	PublishDate NullTime       `json:"publish_date" gorm:"type:date"`
	Summary     string         `json:"summary" gorm:"size:500"`
	Content     string         `json:"content" gorm:"type:text"`
	Status      int            `json:"status" gorm:"default:1"` // 1: 已发布, 0: 草稿
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
