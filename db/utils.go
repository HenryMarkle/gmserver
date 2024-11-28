package db

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
)

func MapRow(row *sql.Row, dest interface{}) error {
	ptrValue := reflect.ValueOf(dest)
	if ptrValue.Kind() != reflect.Ptr || ptrValue.Elem().Kind() != reflect.Struct {
		return errors.New("dest must be a pointer to a struct")
	}

	structValue := ptrValue.Elem()
	structType := structValue.Type()

	columns := make([]interface{}, structType.NumField())
	for i := 0; i < structType.NumField(); i++ {
		columns[i] = structValue.Field(i).Addr().Interface()
	}

	err := row.Scan(columns...)
	if err != nil {
		return fmt.Errorf("failed to scan row: %w", err)
	}

	return nil
}

func MapRows(row *sql.Rows, dest interface{}) error {
	ptrValue := reflect.ValueOf(dest)
	if ptrValue.Kind() != reflect.Ptr || ptrValue.Elem().Kind() != reflect.Struct {
		return errors.New("dest must be a pointer to a struct")
	}

	structValue := ptrValue.Elem()
	structType := structValue.Type()

	columns := make([]interface{}, structType.NumField())
	for i := 0; i < structType.NumField(); i++ {
		columns[i] = structValue.Field(i).Addr().Interface()
	}

	err := row.Scan(columns...)
	if err != nil {
		return fmt.Errorf("failed to scan row: %w", err)
	}

	return nil
}
