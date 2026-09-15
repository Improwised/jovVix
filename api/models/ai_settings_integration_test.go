package models_test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/Improwised/jovvix/api/models"
	"github.com/Improwised/jovvix/api/services"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	_ "github.com/lib/pq"
)

func TestAISettingsPostgresRoundTrip(t *testing.T) {
	dsn := os.Getenv("JOVVIX_TEST_DSN")
	if dsn == "" {
		t.Skip("requires JOVVIX_TEST_DSN")
	}
	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer sqlDB.Close()
	db := goqu.New("postgres", sqlDB)
	var userID string
	found, err := db.From("users").Select("id").Limit(1).ScanVal(&userID)
	if err != nil || !found {
		t.Fatalf("user lookup: found=%v err=%v", found, err)
	}
	ciphertext, salt, nonce, err := services.EncryptAIKey(userID, []byte("vault-test-password"), []byte("test-provider-key"))
	if err != nil {
		t.Fatal(err)
	}
	m := models.InitAISettingsModel(db)
	if err := m.Save(models.AISettings{UserID: userID, Provider: "test", BaseURL: "https://example.com/v1", Model: "test", Ciphertext: ciphertext, Salt: salt, Nonce: nonce}); err != nil {
		t.Fatalf("save: %v", err)
	}
	defer m.Delete(userID)
	stored, err := m.Get(userID)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	plain, err := services.DecryptAIKey(userID, stored, []byte("vault-test-password"))
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	if string(plain) != "test-provider-key" {
		t.Fatal("unexpected plaintext")
	}
}
