package appinforepositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/Aritiaya50217/E-CommerceRESTAPIs/modules/appinfo"
	"github.com/jmoiron/sqlx"
)

type IAppinfoRepository interface {
	FindCategory(req *appinfo.CategoryFilter) ([]*appinfo.Category, error)
	InsertCategory(req []*appinfo.Category) error
	DeleteCategory(categoryId int) error
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

func (r *appinfoRepository) InsertCategory(req []*appinfo.Category) error {
	ctx := context.Background()

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	query := "insert into categories (title) values "

	for i, cat := range req {
		if i != len(req)-1 {
			query += fmt.Sprintf(`
		('%v'),`, cat.Title)
		} else {
			query += fmt.Sprintf(`
		('%v')`, cat.Title)
		}
	}
	rows, err := tx.QueryxContext(ctx, query)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("insert categories failed: %v", err.Error())
	}

	var i int
	for rows.Next() {
		if err := rows.Scan(&req[i].Id); err != nil {
			return fmt.Errorf("scan categories id failed: %v", err)
		}
		i++
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *appinfoRepository) DeleteCategory(categoryId int) error {
	ctx := context.Background()
	query := "delete from categories where id = ? "

	if _, err := r.db.DB.ExecContext(ctx, query, categoryId); err != nil {
		return fmt.Errorf("delete cateogry failed: %v", err)
	}
	return nil
}
