package model

import (
	"encoding/json"
	"fmt"
	"github.com/cuihairu/salon/internal/utils"
	"gorm.io/gorm"
)

type Tag struct {
	Key   string `json:"key"`
	Label string `json:"label"`
}

type Tags []Tag

type Geographic struct {
	Province Tag `json:"province"`
	City     Tag `json:"city"`
}

func (g *Geographic) Scan(value interface{}) error {
	if b, ok := value.([]byte); ok {
		return json.Unmarshal(b, g)
	}
	return fmt.Errorf("failed to scan geographic: %v", value)
}

func (t *Tags) Scan(value interface{}) error {
	if b, ok := value.([]byte); ok {
		return json.Unmarshal(b, t)
	}
	return fmt.Errorf("failed to scan geographic: %v", value)
}

type Admin struct {
	gorm.Model
	Name       string                      `gorm:"type:varchar(255);unique;index:name_idx" json:"name"`
	Phone      *string                     `gorm:"type:varchar(255);unique;index:phone_idx" json:"phone,omitempty"`
	Email      *string                     `gorm:"type:varchar(255);unique;index:email_idx" json:"email,omitempty"`
	Country    *string                     `gorm:"type:varchar(255)" json:"country,omitempty"`
	Avatar     string                      `gorm:"type:text;not null" json:"avatar"`
	Signature  string                      `json:"signature"`
	Title      string                      `json:"title"`
	Group      string                      `json:"group"`
	Address    string                      `json:"address"`
	Geographic utils.JsonField[Geographic] `gorm:"type:jsonb" json:"geographic"`
	Role       string                      `json:"role" json:"role"`
	Tags       utils.JsonField[Tags]       `gorm:"type:jsonb" json:"tags"`
	Password   []byte                      `gorm:"type:bytea;not null" json:"-"`
	Salt       []byte                      `gorm:"type:bytea;not null" json:"-"`
}
