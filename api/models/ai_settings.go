package models

import (
	"database/sql"
	"github.com/doug-martin/goqu/v9"
)

const AISettingsTable = "ai_settings"

type AISettings struct {
	UserID     string `db:"user_id"`
	Provider   string `db:"provider"`
	BaseURL    string `db:"base_url"`
	Model      string `db:"model"`
	Ciphertext []byte `db:"encrypted_api_key"`
	Salt       []byte `db:"salt"`
	Nonce      []byte `db:"nonce"`
}
type AISettingsModel struct{ db *goqu.Database }

func InitAISettingsModel(db *goqu.Database) *AISettingsModel { return &AISettingsModel{db} }
func (m *AISettingsModel) Get(id string) (AISettings, error) {
	var s AISettings
	ok, e := m.db.From(AISettingsTable).Where(goqu.Ex{"user_id": id}).ScanStruct(&s)
	if e != nil {
		return s, e
	}
	if !ok {
		return s, sql.ErrNoRows
	}
	return s, nil
}
func (m *AISettingsModel) Save(s AISettings) error {
	r := goqu.Record{"user_id": s.UserID, "provider": s.Provider, "base_url": s.BaseURL, "model": s.Model, "encrypted_api_key": s.Ciphertext, "salt": s.Salt, "nonce": s.Nonce, "updated_at": goqu.L("now()")}
	_, e := m.db.Insert(AISettingsTable).Rows(r).OnConflict(goqu.DoUpdate("user_id", goqu.Record{"provider": s.Provider, "base_url": s.BaseURL, "model": s.Model, "encrypted_api_key": s.Ciphertext, "salt": s.Salt, "nonce": s.Nonce, "updated_at": goqu.L("now()")})).Prepared(true).Executor().Exec()
	return e
}
func (m *AISettingsModel) Delete(id string) error {
	_, e := m.db.Delete(AISettingsTable).Where(goqu.Ex{"user_id": id}).Executor().Exec()
	return e
}
