package data

import (
	"fmt"
	"github.com/cuihairu/salon/internal/model"
	"github.com/cuihairu/salon/internal/utils"
	"gorm.io/gorm"
)

type AdminRepository struct {
	utils.Repository[model.Admin]
}

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{
		Repository: utils.NewRepository[model.Admin](db),
	}
}

func (a *AdminRepository) Create(admin *model.Admin) error {
	if len(admin.Name) == 0 {
		return fmt.Errorf("name cannot be empty")
	}
	if len(admin.Password) == 0 {
		return fmt.Errorf("password cannot be empty")
	}
	if len(admin.Salt) == 0 {
		return fmt.Errorf("salt cannot be empty")
	}
	if len(admin.Title) == 0 {
		admin.Title = "admin"
	}
	if len(admin.Signature) == 0 {
		admin.Signature = "admin signature"
	}
	if len(admin.Address) == 0 {
		admin.Address = "浙江省杭州市滨江区"
	}
	if len(admin.Role) == 0 {
		admin.Role = "admin"
	}
	if len(admin.Group) == 0 {
		admin.Group = "admin"
	}
	if admin.Tags.IsNil() {
		tags := model.Tags{model.Tag{Key: "admin", Label: "admin"}}
		admin.Tags.SetData(&tags)
	}
	if admin.Geographic.IsNil() {
		geographic := &model.Geographic{
			Province: model.Tag{Key: "zhejiang", Label: "zhejiang"},
			City:     model.Tag{Key: "hangzhou", Label: "hangzhou"},
		}
		admin.Geographic.SetData(geographic)
	}
	if admin.Email == nil || !utils.IsEmail(*admin.Email) {
		email := fmt.Sprintf("%s@youngs.fun", admin.Name)
		admin.Email = &email
	}
	if !utils.IsURL(admin.Avatar) {
		admin.Avatar = "https://gw.alipayobjects.com/zos/rmsportal/BiazfanxmamNRoxxVxka.png"
	}
	if admin.Country == nil {
		country := "China"
		admin.Country = &country
	}
	return a.Repository.Create(admin)
}

func (a *AdminRepository) FindByName(name string) (*model.Admin, error) {
	if len(name) == 0 {
		return nil, fmt.Errorf("name cannot be empty")
	}
	admins, err := a.FindByField("name", name)
	if err != nil {
		return nil, err
	}
	if len(admins) == 0 {
		return nil, nil // NotFound
	}
	return &admins[0], nil
}

func (a *AdminRepository) FindByPhone(phone string) (*model.Admin, error) {
	if len(phone) == 0 {
		return nil, fmt.Errorf("phoneqq cannot be empty")
	}
	admins, err := a.FindByField("phone", phone)
	if err != nil {
		return nil, err
	}
	if len(admins) == 0 {
		return nil, nil // NotFound
	}
	return &admins[0], nil
}

func (a *AdminRepository) FindByEmail(email string) (*model.Admin, error) {
	if len(email) == 0 {
		return nil, fmt.Errorf("email cannot be empty")
	}
	admins, err := a.FindByField("email", email)
	if err != nil {
		return nil, err
	}
	if len(admins) == 0 {
		return nil, nil // NotFound
	}
	return nil, err
}
