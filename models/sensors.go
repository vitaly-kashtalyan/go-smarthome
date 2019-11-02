package models

import (
	"database/sql"
	"github.com/jinzhu/gorm"
	"github.com/vitaly-kashtalyan/go-smarthome/db"
	"log"
	"time"
)

type Sensors struct {
	Pin         int       `gorm:"default:NULL" sql:"type:tinyint(3);" binding:"required" json:"pin"`
	DecSensor   string    `gorm:"default:NULL" sql:"type:char(16);" binding:"min=16,max=16" json:"dec"`
	Temperature float32   `gorm:"default:NULL" sql:"type:decimal(4,2);" json:"temperature"`
	Humidity    float32   `gorm:"default:NULL" sql:"type:decimal(4,2);" json:"humidity"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"update_at"`
}

type SensorsHistory struct {
	ID          uint      `gorm:"primary_key" json:"id"`
	Pin         int       `gorm:"default:NULL" sql:"type:tinyint(3);index:idx_pin" binding:"required" json:"pin"`
	DecSensor   string    `gorm:"default:NULL" sql:"type:char(16);" binding:"min=16,max=16" json:"dec"`
	Temperature float32   `gorm:"default:NULL" sql:"type:decimal(4,2);" json:"temperature"`
	Humidity    float32   `gorm:"default:NULL" sql:"type:decimal(4,2);" json:"humidity"`
	CreatedAt   time.Time `sql:"index" json:"date"`
}

type RelayStateHistory struct {
	ID        uint          `gorm:"primary_key" json:"id"`
	RelayId   sql.NullInt32 `gorm:"default:NULL" sql:"type:tinyint(3);" binding:"required" json:"relay_id"`
	State     sql.NullInt32 `gorm:"default:NULL" sql:"type:tinyint(3);" binding:"required" json:"state"`
	CreatedAt time.Time     `gorm:"default:NULL" json:"create_at"`
	UpdatedAt sql.NullTime  `gorm:"default:NULL" json:"update_at"`
}

func (s *Sensors) AfterSave(scope *gorm.Scope) (err error) {
	sensorsHistory := SensorsHistory{}
	db.GetDB().Where(SensorsHistory{Pin: s.Pin, DecSensor: s.DecSensor}).
		Order("created_at desc").
		Limit(1).Find(&sensorsHistory)

	if sensorsHistory.ID == 0 || sensorsHistory.ID > 0 && sensorsHistory.Temperature != s.Temperature {
		var newRecord = SensorsHistory{
			Pin:         s.Pin,
			DecSensor:   s.DecSensor,
			Temperature: s.Temperature,
			Humidity:    s.Humidity}
		if err := db.GetDB().Create(&newRecord).Error; err != nil {
			log.Println("error creating sensor history record: ", err)
		}
	}
	return
}
