package models_test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/Improwised/jovvix/api/models"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A profile name change must keep traits.email (GitHub #65). Before, the update wrote
// {"name":...,"email":""} over the identity's traits. Needs JOVVIX_TEST_DSN; everything rolls back.
func TestUpdateKratosIdentityTraitsKeepsEmail(t *testing.T) {
	dsn := os.Getenv("JOVVIX_TEST_DSN")
	if dsn == "" {
		t.Skip("requires JOVVIX_TEST_DSN")
	}
	sqlDB, err := sql.Open("postgres", dsn)
	require.NoError(t, err)
	t.Cleanup(func() { sqlDB.Close() })

	tx, err := goqu.New("postgres", sqlDB).Begin()
	require.NoError(t, err)
	t.Cleanup(func() { _ = tx.Rollback() })

	var hasKratos bool
	require.NoError(t, tx.QueryRow(`SELECT to_regclass('kratos.identities') IS NOT NULL`).Scan(&hasKratos))
	if hasKratos {
		t.Skip("kratos.identities exists with Kratos's own columns; use a database without Kratos migrations")
	}
	_, err = tx.Exec(`CREATE SCHEMA IF NOT EXISTS kratos; CREATE TABLE kratos.identities (id uuid PRIMARY KEY, traits json NOT NULL)`)
	require.NoError(t, err)

	id := uuid.NewString()
	_, err = tx.Exec(`INSERT INTO kratos.identities VALUES ($1, '{"email":"a@example.com","name":{"first":"Old","last":"Name"}}')`, id)
	require.NoError(t, err)

	require.NoError(t, models.UpdateKratosIdentityTraits(tx, id, []byte(`{"first":"New","last":"Person"}`)))

	var email, first, last string
	require.NoError(t, tx.QueryRow(`SELECT traits->>'email', traits->'name'->>'first', traits->'name'->>'last' FROM kratos.identities WHERE id=$1`, id).Scan(&email, &first, &last))
	assert.Equal(t, "a@example.com", email)
	assert.Equal(t, "New", first)
	assert.Equal(t, "Person", last)
}
