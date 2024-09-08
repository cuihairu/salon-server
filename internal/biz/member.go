package biz

import (
	"github.com/cuihairu/salon/internal/data"
	"github.com/cuihairu/salon/internal/model"
	"go.uber.org/zap"
)

type MemberBiz struct {
	memberRepo *data.MemberRepository
	logger     *zap.Logger
}

func NewMemberBiz(dataStore *data.DataStore, logger *zap.Logger) *MemberBiz {
	return &MemberBiz{
		memberRepo: dataStore.MemberRepo,
		logger:     logger,
	}
}

func (biz *MemberBiz) GetAllMembers() ([]model.Member, error) {
	return biz.memberRepo.FindAll()
}

func (biz *MemberBiz) GetMemberByID(id uint) (*model.Member, error) {
	return biz.memberRepo.FindByID(id)
}

func (biz *MemberBiz) UpdateMember(member *model.Member) error {
	return biz.memberRepo.Update(member)
}

func (biz *MemberBiz) DeleteMember(id uint) error {
	return biz.memberRepo.Delete(id)
}
