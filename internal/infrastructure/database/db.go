package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"github.com/Gustrb/ccanalytics/internal/infrastructure/common"
)

var (
	db *sql.DB
)

var (
	ErrDBAlreadyConnected = fmt.Errorf("database connection already established")
)

func Connect(ctx context.Context, dataSourceName string) (common.CleanupFunction, error) {
	if db != nil {
		return nil, ErrDBAlreadyConnected
	}

	var err error
	db, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	cleanup := func() error {
		return db.Close()
	}

	if err := db.PingContext(ctx); err != nil {
		if cerr := cleanup(); cerr != nil {
			return nil, fmt.Errorf("failed to ping database: %w; also failed to clean up database connection: %v", err, cerr)
		}

		return nil, err
	}

	return cleanup, nil
}

// getScanDest returns pointers to struct fields in the order of the given columns,
// using sql struct tags (or lowercase field name) to match column names.
func getScanDest(obj any, columns []string) ([]any, error) {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("getScanDest: expected struct, got %s", v.Kind())
	}
	t := v.Type()

	columnToField := make(map[string]int)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		colName := field.Tag.Get("sql")
		if colName == "" {
			colName = strings.ToLower(field.Name)
		}
		columnToField[colName] = i
	}

	var dest []any
	for _, col := range columns {
		colLower := strings.ToLower(col)
		i, ok := columnToField[colLower]
		if !ok {
			return nil, fmt.Errorf("getScanDest: no struct field for column %q", col)
		}
		dest = append(dest, v.Field(i).Addr().Interface())
	}
	return dest, nil
}

func SelectContext[T any](ctx context.Context, query string, args ...any) ([]*T, error) {
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	var result []*T

	for rows.Next() {
		var o T
		dest, err := getScanDest(&o, columns)
		if err != nil {
			return nil, err
		}
		if err := rows.Scan(dest...); err != nil {
			return nil, err
		}
		result = append(result, &o)
	}

	return result, nil
}

type HasID interface {
	GetID() int
	SetID(int)
}

// extractInsertArgs uses reflection to extract struct field values for INSERT,
// excluding the ID field (by name "ID" or sql tag "id"). Values are returned
// in struct field order. The INSERT query must have placeholders in the same order.
func extractInsertArgs(obj any) ([]any, error) {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("extractInsertArgs: expected struct, got %s", v.Kind())
	}
	t := v.Type()
	var args []any
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		if field.Name == "ID" {
			continue
		}
		if tag := field.Tag.Get("sql"); tag == "id" {
			continue
		}
		args = append(args, v.Field(i).Interface())
	}
	return args, nil
}

func InsertContext[T HasID](ctx context.Context, query string, obj T) error {
	args, err := extractInsertArgs(obj)
	if err != nil {
		return err
	}

	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	obj.SetID(int(id))

	return nil
}

func ExecContext(ctx context.Context, query string) (sql.Result, error) {
	return db.ExecContext(ctx, query)
}

func WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(ctx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			return fmt.Errorf("transaction function error: %w; also failed to rollback transaction: %v", err, rerr)
		}

		return err
	}

	return tx.Commit()
}
