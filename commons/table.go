package commons

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PrimaryKeyResetter interface {
	ResetPrimaryKey()
}

type TableBase struct {
	ID        string    `gorm:"type:char(36);primary_key;" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (b *TableBase) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}

	return nil
}

func (b *TableBase) ResetPrimaryKey() {
	b.ID = ""
}
