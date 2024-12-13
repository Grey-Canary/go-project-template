package testhelper

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func NewTestPgxConn(t *testing.T) *pgx.Conn {
	t.Helper()

	ctx := context.Background()

	envFile, _ := godotenv.Read("../../.env")

	connString := envFile["TEST_DATABASE_URL"]
	if connString == "" {
		t.Skipf("skipping due to missing environment variable %v", "TEST_DATABASE_URL")
	}

	config, err := pgx.ParseConfig(connString)
	require.NoError(t, err)

	conn, err := pgx.ConnectConfig(ctx, config)
	require.NoError(t, err)

	t.Cleanup(func() {
		conn.Close(ctx)
	})

	return conn
}
