package sqlite

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFlush(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store, clean := NewTestStore(t)
	defer clean(t)

	err := store.execTrans(ctx, `CREATE TABLE test_table_1 (id TEXT NOT NULL PRIMARY KEY)`)
	require.NoError(t, err)

	err = store.execTrans(ctx, `INSERT INTO test_table_1 (id) VALUES ("one"), ("two"), ("three")`)
	require.NoError(t, err)

	vals, err := store.queryToStrings(`SELECT * FROM test_table_1`)
	require.NoError(t, err)
	require.Equal(t, 3, len(vals))

	store.Flush(context.Background())

	vals, err = store.queryToStrings(`SELECT * FROM test_table_1`)
	require.NoError(t, err)
	require.Equal(t, 0, len(vals))
}

func TestBackupSqlStore(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store, clean := NewTestStore(t)
	defer clean(t)

	err := store.execTrans(ctx, `CREATE TABLE test_table_1 (id TEXT NOT NULL PRIMARY KEY)`)
	require.NoError(t, err)

	err = store.execTrans(ctx, `INSERT INTO test_table_1 (id) VALUES ("one"), ("two"), ("three")`)
	require.NoError(t, err)

	tempDir, err := ioutil.TempDir("", "")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	backupPath := tempDir + "/db.sqlite"
	dest, err := os.Create(backupPath)
	require.NoError(t, err)

	store.BackupSqlStore(ctx, dest)

	b1, err := ioutil.ReadFile(store.path)
	require.NoError(t, err)
	b2, err := ioutil.ReadFile(backupPath)
	require.NoError(t, err)

	require.True(t, bytes.Equal(b1, b2))
}

func TestUserVersion(t *testing.T) {
	t.Parallel()

	store, clean := NewTestStore(t)
	defer clean(t)
	ctx := context.Background()

	err := store.execTrans(ctx, `PRAGMA user_version=12`)
	require.NoError(t, err)

	got, err := store.userVersion()
	require.NoError(t, err)
	require.Equal(t, 12, got)
}

func TestTableNames(t *testing.T) {
	t.Parallel()

	store, clean := NewTestStore(t)
	defer clean(t)
	ctx := context.Background()

	err := store.execTrans(ctx, `CREATE TABLE test_table_1 (id TEXT NOT NULL PRIMARY KEY);
	CREATE TABLE test_table_3 (id TEXT NOT NULL PRIMARY KEY);
	CREATE TABLE test_table_2 (id TEXT NOT NULL PRIMARY KEY);`)
	require.NoError(t, err)

	got, err := store.tableNames()
	require.NoError(t, err)
	require.Equal(t, []string{"test_table_1", "test_table_3", "test_table_2"}, got)
}
