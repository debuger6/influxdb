package http

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestBackupMetaService(t *testing.T) {
	b := &BackupBackend{}
	h := NewBackupHandler(b)
	h.BackupService = &FakeBackupService{}
	h.SqlBackupService = &FakeSqlBackupService{}

	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	h.handleBackupMetadata(rr, r)

	rs := rr.Result()

	require.Equal(t, rs.StatusCode, http.StatusOK)

	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
	require.NoError(t, err)

	require.Contains(t, string(body), "kv store lol")
	require.Contains(t, string(body), "sql store lol")
}

type FakeBackupService struct{}

func (f *FakeBackupService) BackupKVStore(ctx context.Context, w io.Writer) error {
	io.Copy(w, strings.NewReader("kv store lol"))
	return nil
}

func (f *FakeBackupService) BackupShard(ctx context.Context, w io.Writer, shardID uint64, since time.Time) error {
	return nil
}

type FakeSqlBackupService struct{}

func (f *FakeSqlBackupService) BackupSqlStore(ctx context.Context, w io.Writer) error {
	io.Copy(w, strings.NewReader("sql store lol"))
	return nil
}
