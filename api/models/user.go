package models

import (
	"database/sql"

	goqu "github.com/doug-martin/goqu/v9"
	"github.com/rs/xid"
)

// This boilerplate we are storing password in plan format!

// UserTable represent table name
const UserTable = "users"

// User model
type User struct {
	ID        string         `json:"id"`
	KratosID  string         `json:"kratos_id" db:"kratos_id"`
	FirstName string         `json:"first_name" db:"first_name" validate:"required"`
	LastName  string         `json:"last_name" db:"last_name" validate:"required"`
	Email     string         `json:"email" db:"email" validate:"required"`
	Username  string         `json:"username" db:"username" validate:"required"`
	Password  sql.NullString `json:"-" db:"password"`
	Roles     string         `json:"roles,omitempty" db:"roles" validate:"required"`
	CreatedAt string         `json:"created_at,omitempty" db:"created_at,omitempty"`
	UpdatedAt string         `json:"updated_at,omitempty" db:"updated_at,omitempty"`
}

// UserModel implements user related database operations
type UserModel struct {
	db *goqu.Database
}

// InitUserModel Init model
func InitUserModel(goqu *goqu.Database) (UserModel, error) {
	return UserModel{
		db: goqu,
	}, nil
}

// GetUsers list all users
func (model *UserModel) GetUsers() ([]User, error) {
	var users []User
	if err := model.db.From(UserTable).ScanStructs(&users); err != nil {
		return nil, err
	}
	return users, nil
}

// GetUser get user by id
func (model *UserModel) GetById(id string) (User, error) {
	user := User{}
	found, err := model.db.From(UserTable).Where(goqu.Ex{
		"id": id,
	}).Select(
		"id",
		"first_name",
		"last_name",
		"email",
		"username",
		"roles",
	).ScanStruct(&user)

	if err != nil {
		return user, err
	}

	if !found {
		return user, sql.ErrNoRows
	}

	return user, err
}

// InsertUser retrieve user
func (model *UserModel) InsertUser(user User) (User, error) {
	user.ID = xid.New().String()

	_, err := model.db.Insert(UserTable).Rows(
		goqu.Record{
			"id":         user.ID,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"email":      user.Email,
			"password":   user.Password,
			"roles":      user.Roles,
			"username":   user.Username,
		},
	).Executor().Exec()
	if err != nil {
		return user, err
	}

	user, err = model.GetById(user.ID)
	return user, err
}

func (model *UserModel) InsertKratosUser(user User) error {
	rows, err := model.db.Select(goqu.L("EXISTS ?", model.db.Select().From(UserTable).Where(goqu.L("kratos_id = ?", user.KratosID)))).Executor().Query()
	if err != nil {
		return err
	}

	var exists bool
	for rows.Next() {
		err := rows.Scan(&exists)
		if err != nil {
			return err
		}
	}

	if !exists {
		user.ID = xid.New().String()
		_, err := model.db.Insert(UserTable).Rows(
			goqu.Record{
				"id":         user.ID,
				"kratos_id":  user.KratosID,
				"first_name": user.FirstName,
				"last_name":  user.LastName,
				"email":      user.Email,
				"created_at": user.CreatedAt,
				"updated_at": user.UpdatedAt,
				"username":   user.Username,
				"roles":      user.Roles,
			},
		).Executor().Exec()
		if err != nil {
			return err
		}
	}
	return nil
}

func (model *UserModel) GetUserByEmailAndPassword(email string, password string) (User, error) {
	user := User{}
	found, err := model.db.From(UserTable).Where(goqu.Ex{
		"email":    email,
		"password": password,
	}).Select(
		"id",
		"first_name",
		"last_name",
		"email",
	).ScanStruct(&user)

	if err != nil {
		return user, err
	}

	if !found {
		return user, sql.ErrNoRows
	}

	return user, err
}

func (model *UserModel) CountUsers() (int64, error) {
	return model.db.From(UserTable).Count()
}

func (model *UserModel) IsUniqueEmail(email string) (bool, error) {
	query := model.db.From("users").Select(goqu.I("id")).Where(goqu.Ex{"email": email}).Limit(1)

	// Execute the query
	rows, err := query.Executor().Query()
	if err != nil {
		return false, err
	}
	defer rows.Close()

	// Check if any rows were returned
	return !rows.Next(), err
}

func (model *UserModel) GetUserRole(userID string) (string, error) {
	var role string = "not found"

	found, err := model.db.From(UserTable).Where(goqu.Ex{
		"id": userID,
	}).Select(
		"roles",
	).ScanStruct(&role)

	if err != nil {
		return role, err
	}

	if !found {
		return role, sql.ErrNoRows
	}

	return role, err
}

func (model *UserModel) GetUserByKratosID(kratosID string) (User, error) {
	user := User{}

	ok, err := model.db.From(UserTable).Where(goqu.Ex{
		"kratos_id": kratosID,
	}).ScanStruct(&user)

	if err != nil {
		return User{}, err
	}

	if !ok {
		return User{}, sql.ErrNoRows
	} else {
		return user, nil
	}
}
