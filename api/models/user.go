package models

import (
	"database/sql"
	"strings"

	goqu "github.com/doug-martin/goqu/v9"
	"github.com/rs/xid"
	"go.uber.org/zap"
)

// This boilerplate we are storing password in plan format!

// UserTable represent table name
const UserTable = "users"
const KratosIdentityTable = "kratos.identities"

// User model
type User struct {
	ID        string         `json:"id"`
	KratosID  sql.NullString `json:"kratos_id" db:"kratos_id"`
	FirstName string         `json:"first_name" db:"first_name" validate:"required"`
	LastName  string         `json:"last_name" db:"last_name" validate:"required"`
	Email     string         `json:"email" db:"email" validate:"required"`
	Username  string         `json:"username" db:"username" validate:"required"`
	Password  sql.NullString `json:"-" db:"password"`
	Roles     string         `json:"roles,omitempty" db:"roles" validate:"required"`
	ImageKey  string         `json:"img_key,omitempty" db:"img_key"`
	CreatedAt string         `json:"created_at,omitempty" db:"created_at,omitempty"`
	UpdatedAt string         `json:"updated_at,omitempty" db:"updated_at,omitempty"`
}

// UserModel implements user related database operations
type UserModel struct {
	db     *goqu.Database
	logger *zap.Logger
}

// InitUserModel Init model
func InitUserModel(goqu *goqu.Database, logger *zap.Logger) (UserModel, error) {
	return UserModel{
		db: goqu,
	}, nil
}

// GetUser get user by id
func (model *UserModel) GetById(id string) (User, error) {
	user := User{}
	found, err := model.db.From(UserTable).Where(goqu.Ex{
		"id": id,
	}).Select(
		"id",
		"kratos_id",
		"first_name",
		"last_name",
		"email",
		"username",
		"roles",
		"img_key",
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
			"img_key":    user.ImageKey,
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
				"kratos_id":  user.KratosID.String,
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

func (model *UserModel) CountUsers() (int64, error) {
	return model.db.From(UserTable).Count()
}

func (model *UserModel) IsUniqueEmailExceptId(userId, email string) (bool, error) {
	query := model.db.From("users").Select(goqu.I("id")).Where(
		goqu.Ex{"email": email},
		goqu.C("id").Neq(userId),
	).Limit(1)

	// Execute the query
	rows, err := query.Executor().Query()
	if err != nil {
		return false, err
	}
	defer rows.Close()

	// Check if any rows were returned
	return !rows.Next(), err
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

func (model *UserModel) UpdateKratosUserDetails(reqUser User, userMetadata []byte) error {
	isOk := false
	transaction, err := model.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if isOk {
			err := transaction.Commit()
			if err != nil {
				model.logger.Error("error during commit in register question", zap.Error(err))
			}
		} else {
			err := transaction.Rollback()
			if err != nil {
				model.logger.Error("error during rollback in register question", zap.Error(err))
			}
		}
	}()

	user, err := model.GetById(reqUser.ID)
	if err != nil {
		return err
	}

	err = UpdateKratosIdentifiers(transaction, user.Email, reqUser.Email)
	if err != nil {
		model.logger.Debug("error in UpdateKratosIdentifiers", zap.Error(err))
		return err
	}

	err = UpdateKratosIdentityTraits(transaction, strings.TrimSpace(user.KratosID.String), userMetadata)
	if err != nil {
		model.logger.Debug("error in UpdateKratosIdentityTraits", zap.Error(err))
		return err
	}

	err = UpdateUserMetadata(transaction, reqUser)
	if err != nil {
		model.logger.Debug("error in UpdateKratosIdentityTraits", zap.Error(err))
		return err
	}

	isOk = true
	return nil
}

func UpdateKratosIdentifiers(transaction *goqu.TxDatabase, oldEmail, newEmail string) error {

	_, err := transaction.Update("kratos.identity_credential_identifiers").Set(goqu.Record{
		"identifier": newEmail,
	}).Where(goqu.Ex{
		"identifier": oldEmail,
	}).Executor().Exec()

	return err
}

func UpdateKratosIdentityTraits(transaction *goqu.TxDatabase, kratosId string, data []byte) error {

	record := goqu.Record{
		"traits": data,
	}

	_, err := transaction.Update(KratosIdentityTable).Set(record).Where(goqu.Ex{
		"id": kratosId,
	}).Executor().Exec()

	return err
}

func UpdateUserMetadata(transaction *goqu.TxDatabase, user User) error {

	record := goqu.Record{
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
	}

	_, err := transaction.Update(UserTable).Set(record).Where(goqu.Ex{
		"id": user.ID,
	}).Executor().Exec()

	return err
}

// Delete user from user table and also returns the Kratos ID associated with the user
func (model *UserModel) DeleteUserById(transaction *goqu.TxDatabase, id string) (string, error) {

	var kratosId string
	_, err := transaction.Delete(UserTable).Where(goqu.Ex{"id": id}).Returning("kratos_id").Executor().ScanVal(&kratosId)
	if err != nil {
		return kratosId, err
	}

	return kratosId, nil
}
