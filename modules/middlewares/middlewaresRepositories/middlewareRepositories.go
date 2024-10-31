package middlewaresrepositories

import (
	"fmt"

	"github.com/Aritiaya50217/E-CommerceRESTAPIs/modules/middlewares"
	"github.com/jmoiron/sqlx"
)

type IMiddlewaresRepository interface {
	FindAccessToken(userId, accessToken string) bool
	FindRole() ([]*middlewares.Role, error)
}

type middlewareRepository struct {
	db *sqlx.DB
}

func MiddlewaresRepository(db *sqlx.DB) IMiddlewaresRepository {
	return &middlewareRepository{
		db: db,
	}
}

func (r *middlewareRepository) FindAccessToken(userId, accessToken string) bool {
	query := `
	SELECT
		(case when count(*) = 1 then true else false end) 
	FROM "oauth"
	WHERE "user_id" = $1
	AND "access_token" = $2;`
	var check bool
	if err := r.db.Get(&check, query, userId, accessToken); err != nil {
		return false
	}
	return true

}

func (r *middlewareRepository) FindRole() ([]*middlewares.Role, error) {
	query := `
	select 
	"id","title"
	from "roles"
	order by "id" desc;
	`
	roles := make([]*middlewares.Role, 0)
	if err := r.db.Select(&roles, query); err != nil {
		return nil, fmt.Errorf("roles are empty")
	}
	return roles, nil
}
