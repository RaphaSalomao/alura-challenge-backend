package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	Id        uuid.UUID `sql:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	CreatedAt time.Time `sql:"type:timestamptz"`
	UpdatedAt time.Time `sql:"type:timestamptz"`
}

func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	b.Id = uuid.New()
	b.CreatedAt = time.Now()
	b.UpdatedAt = time.Now()
	return
}

func (b *Base) BeforeSave(tx *gorm.DB) (err error) {
	b.UpdatedAt = time.Now()
	return
}
