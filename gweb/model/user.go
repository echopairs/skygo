package model

// user_id -> role_id -> (access ids) -> (access names)

type User struct {
	ID       int      `db:"id"`
	UserName string   `db:"name"`
	Password string   `json:"-" db:"password"`
	Salt     string   `json:"-" db:"salt"`
	Roles    []string `db:"-"`
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
