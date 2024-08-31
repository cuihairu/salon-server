package data

import (
	"github.com/cuihairu/salon/internal/model"
	"github.com/cuihairu/salon/internal/utils"
	"gorm.io/gorm"
)

type MemberRepository struct {
	utils.Repository[model.Member]
}

func NewMemberRepository(db *gorm.DB) *MemberRepository {
	return &MemberRepository{
		Repository: utils.NewRepository[model.Member](db),
	}
}
