package othooks

import (
	"context"
	"database/sql"
	"testing"

	sqlite3 "github.com/mattn/go-sqlite3"
	"github.com/qustavo/sqlhooks/v2"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var (
	tracer trace.Tracer
)

func init() {
	tracer = otel.Tracer("mock")
	driver := sqlhooks.Wrap(&sqlite3.SQLiteDriver{}, New(tracer))
	sql.Register("ot", driver)
}

func TestSpansAreRecorded(t *testing.T) {
	db, err := sql.Open("ot", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	ctx, span := tracer.Start(context.Background(), "sql")

	{
		rows, err := db.QueryContext(ctx, "SELECT 1+?", "1")
		require.NoError(t, err)
		rows.Close()
	}

	{
		rows, err := db.QueryContext(ctx, "SELECT 1+?", "1")
		require.NoError(t, err)
		rows.Close()
	}

	span.End()

	require.Len(t, span, 3)
}

func TestNoSpansAreRecorded(t *testing.T) {
	db, err := sql.Open("ot", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	rows, err := db.QueryContext(context.Background(), "SELECT 1")
	require.NoError(t, err)
	rows.Close()

}
