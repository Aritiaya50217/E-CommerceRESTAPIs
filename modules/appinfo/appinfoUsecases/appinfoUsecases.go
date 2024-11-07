package appinfousecases

import (
	"github.com/Aritiaya50217/E-CommerceRESTAPIs/modules/appinfo"
	appinforepositories "github.com/Aritiaya50217/E-CommerceRESTAPIs/modules/appinfo/appinfoRepositories"
)

type IAppinfoUsecase interface {
	FindCategory(req *appinfo.CategoryFilter) ([]*appinfo.Category, error)
	InsertCategory(req []*appinfo.Category) error
	DeleteCategory(categoryId int) error
	UpdateCategory(req *appinfo.Category) error
}

type appinfoUsecase struct {
	appinfoRepository appinforepositories.IAppinfoRepository
}

func AppinfoUsecase(appinfoRepository appinforepositories.IAppinfoRepository) IAppinfoUsecase {
	return &appinfoUsecase{
		appinfoRepository: appinfoRepository,
	}
}

func (u *appinfoUsecase) FindCategory(req *appinfo.CategoryFilter) ([]*appinfo.Category, error) {
	category, err := u.appinfoRepository.FindCategory(req)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (u *appinfoUsecase) InsertCategory(req []*appinfo.Category) error {
	return u.appinfoRepository.InsertCategory(req)
}

func (u *appinfoUsecase) DeleteCategory(categoryId int) error {
	return u.appinfoRepository.DeleteCategory(categoryId)
}

func (u *appinfoUsecase) UpdateCategory(req *appinfo.Category) error {
	return u.appinfoRepository.UpdateCategory(req)
}
