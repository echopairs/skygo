package model

import (
	"crypto/md5"
	"fmt"

	"github.com/satori/go.uuid"
)

// user_id -> role_id -> (access ids) -> (access names)

type User struct {
	ID       int      `db:"id"`
	UserName string   `db:"name"`
	Password string   `json:"-" db:"password"`
	Salt     string   `json:"-" db:"salt"`
	Roles    []string `json:"-" db:"-"`
	Email    string   `json:"email" db:"email"`
	IsAdmin  bool     `json:"is_admin" db:"is_admin"`
}

type Role struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type UserRole struct {
	ID     int `db:"id"`
	UserID int `db:"user_id"`
	RoleID int `db:"role_id"`
}

type Access struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type RoleAccess struct {
	ID       int `db:"id"`
	RoleID   int `db:"role_id"`
	AccessID int `db:"access_id"`
}

func (user *User) VerifyPassword(password string) bool {
	p := fmt.Sprintf("%x", md5.Sum([]byte(password+user.Salt)))
	return p == user.Password
}

func RandomString() string {
	return fmt.Sprintf("%x", md5.Sum(uuid.NewV4().Bytes()))
}
