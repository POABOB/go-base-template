package repository

import (
	"github.com/PPOABOB/go-base-template/domain/model"

	"github.com/jinzhu/gorm"
)

// 建立需要實現的接口
type IBaseRepository interface {
	// 初始化 Table
	InitTable() error
	// 根據 ID 來搜尋
	FindBaseByID(int64) (*model.Base, error)
	// 新增一筆 base 資料
	CreateBase(*model.Base) (int64, error)
	// 根據 ID 刪除一筆 base 資料
	DeleteBaseByID(int64) error
	// 更新資料
	UpdateBase(*model.Base) error
	// 搜尋所有 base 資料
	FindAll() ([]model.Base, error)
}

// 建立baseRepository
func NewBaseRepository(db *gorm.DB) IBaseRepository {
	return &BaseRepository{mysqlDb: db}
}

type BaseRepository struct {
	mysqlDb *gorm.DB
}

// 初始化Table
func (u *BaseRepository) InitTable() error {
	return u.mysqlDb.CreateTable(&model.Base{}).Error
}

// 根據 ID 來搜尋
func (u *BaseRepository) FindBaseByID(baseID int64) (base *model.Base, err error) {
	base = &model.Base{}
	return base, u.mysqlDb.First(base, baseID).Error
}

// 新增一筆 base 資料
func (u *BaseRepository) CreateBase(base *model.Base) (int64, error) {
	return base.ID, u.mysqlDb.Create(base).Error
}

// 根據 ID 刪除一筆 base 資料
func (u *BaseRepository) DeleteBaseByID(baseID int64) error {
	return u.mysqlDb.Where("id = ?", baseID).Delete(&model.Base{}).Error
}

// 更新資料
func (u *BaseRepository) UpdateBase(base *model.Base) error {
	return u.mysqlDb.Model(base).Update(base).Error
}

// 搜尋所有 base 資料
func (u *BaseRepository) FindAll() (baseAll []model.Base, err error) {
	return baseAll, u.mysqlDb.Find(&baseAll).Error
}
