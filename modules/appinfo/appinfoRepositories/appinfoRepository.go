package appinforepositories

import (
	"fmt"
	"strings"

	"github.com/Aritiaya50217/E-CommerceRESTAPIs/modules/appinfo"
	"github.com/jmoiron/sqlx"
)

type IAppinfoRepository interface {
	FindCategory(req *appinfo.CategoryFilter) ([]*appinfo.Category, error)
}

type appinfoRepository struct {
	db *sqlx.DB
}

func AppinfoRepository(db *sqlx.DB) IAppinfoRepository {
	return &appinfoRepository{db: db}
}

func (r *appinfoRepository) FindCategory(req *appinfo.CategoryFilter) ([]*appinfo.Category, error) {
	query := "select id,title from categories"
	filterValues := make([]any, 0)
	if req.Title != "" {
		title := "%" + strings.ToLower(req.Title) + "%"
		query += "where lower(title) like ? "
		filterValues = append(filterValues, title)
	}

	category := make([]*appinfo.Category, 0)
	if err := r.db.Select(&category, query, filterValues...); err != nil {
		return nil, fmt.Errorf("select categories failed: %v", err)
	}
	return category, nil
}
