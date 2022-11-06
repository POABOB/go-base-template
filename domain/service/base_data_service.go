package service

import (
	"github.com/POABOB/go-base-template/domain/model"
	"github.com/POABOB/go-base-template/domain/repository"
)

// 接口類型
type IBaseDataService interface {
	AddBase(*model.Base) (int64, error)
	DeleteBase(int64) error
	UpdateBase(*model.Base) error
	FindBaseByID(int64) (*model.Base, error)
	FindAllBase() ([]model.Base, error)
}

// Create
// 注意：return IBaseDataService 接口類型
func NewBaseDataService(baseRepository repository.IBaseRepository) IBaseDataService {
	return &BaseDataService{baseRepository}
}

type BaseDataService struct {
	// 注意：這裡是 IBaseRepository 類型
	BaseRepository repository.IBaseRepository
}

// Insert
func (u *BaseDataService) AddBase(base *model.Base) (int64, error) {
	return u.BaseRepository.CreateBase(base)
}

// Delete
func (u *BaseDataService) DeleteBase(baseID int64) error {
	return u.BaseRepository.DeleteBaseByID(baseID)
}

// Update
func (u *BaseDataService) UpdateBase(base *model.Base) error {
	return u.BaseRepository.UpdateBase(base)
}

// Read By Id
func (u *BaseDataService) FindBaseByID(baseID int64) (*model.Base, error) {
	return u.BaseRepository.FindBaseByID(baseID)
}

// Read All
func (u *BaseDataService) FindAllBase() ([]model.Base, error) {
	return u.BaseRepository.FindAll()
}
