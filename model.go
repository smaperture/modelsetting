package modelsetting

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Setting struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *gorm.DeletedAt
	Key       string
	Value     datatypes.JSONMap
}

func Init(db *gorm.DB) {
	DB = db
}

func (s *Setting) GetBool(property string) bool {
	value := s.Value[property]
	if value == nil {
		return false
	}

	return value.(bool)
}

func (s *Setting) Read() error {
	if err := DB.First(&s, s).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
	}

	if len(s.Value) == 0 {
		s.Value = make(map[string]interface{})
	}

	return nil
}

func (s *Setting) Upsert() error {
	if s.ID != 0 {
		return DB.Save(&s).Error
	}

	value := s.Value

	err := DB.FirstOrCreate(&s, Setting{
		Key: s.Key,
	}).Error
	if err != nil {
		return err
	}

	s.Value = value

	return DB.Save(&s).Error
}
