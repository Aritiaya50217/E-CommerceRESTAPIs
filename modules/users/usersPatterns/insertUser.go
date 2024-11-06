package userpatterns

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Aritiaya50217/E-CommerceRESTAPIs/modules/users"
	"github.com/jmoiron/sqlx"
)

type IInsertUser interface {
	Customer() (IInsertUser, error)
	Admin() (IInsertUser, error)
	Result() (*users.UserPassport, error)
}

type userReq struct {
	id  string
	req *users.UserRegisterReq
	db  *sqlx.DB
}

type customer struct {
	*userReq
}

type admin struct {
	*userReq
}

func InsertUser(db *sqlx.DB, req *users.UserRegisterReq, isAdmin bool) IInsertUser {
	if isAdmin {
		return newAdmin(db, req)
	}
	return newCustomer(db, req)
}

func GetUserByEamil(db *sqlx.DB, email string) (res IInsertUser, u bool) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	word := "%" + email + "%"
	query := "select * from users " +
		"where email like " + word
	if err := db.QueryRowContext(ctx, query); err != nil {
		return nil, false
	}
	return res, true
}

func newCustomer(db *sqlx.DB, req *users.UserRegisterReq) IInsertUser {
	return &customer{
		userReq: &userReq{
			req: req,
			db:  db,
		},
	}
}

func newAdmin(db *sqlx.DB, req *users.UserRegisterReq) IInsertUser {
	return &admin{
		userReq: &userReq{
			req: req,
			db:  db,
		},
	}
}

func (f *userReq) Customer() (IInsertUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	query := " INSERT INTO users (username,password,email,role_id) " +
		"VALUES ( " + "'" + f.req.Username + "'" + "," + "'" +
		f.req.Password + "'" + "," + "'" + f.req.Email + "'" + "," + "1 ); "

	if err := f.db.QueryRowContext(
		ctx,
		query,
	).Scan(&f.id); err != nil {
		return nil, err
	}
	return f, nil
}

func (f *userReq) Admin() (IInsertUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := " INSERT INTO users (username,password,email,role_id) " +
		"VALUES ( " + "'" + f.req.Username + "'" + "," + "'" +
		f.req.Password + "'" + "," + "'" + f.req.Email + "'" + "," + "2 ); "

	if err := f.db.QueryRowContext(
		ctx,
		query,
	).Scan(); err != nil {
		return nil, err
	}
	return f, nil
}

func (f *userReq) Result() (*users.UserPassport, error) {
	query := `
	SELECT
		json_build_object(
			'user', "t",
			'token', NULL
		)
	FROM (
		SELECT
			"u"."id",
			"u"."email",
			"u"."username",
			"u"."role_id"
		FROM "users" "u"
		WHERE "u"."id" = $1
	) AS "t"`

	data := make([]byte, 0)
	if err := f.db.Get(&data, query, f.id); err != nil {
		return nil, fmt.Errorf("get user failed: %v", err)
	}

	user := new(users.UserPassport)
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, fmt.Errorf("unmarshal user failed: %v", err)
	}
	return user, nil
}
