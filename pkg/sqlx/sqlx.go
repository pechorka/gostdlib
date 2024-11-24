package sqlx

import (
	"database/sql"
	"fmt"
	"reflect"
	"slices"
	"strings"
)

type DB struct {
	DB *sql.DB
}

type Tx struct {
	TX *sql.Tx
}

type Stmt struct {
	Stmt *sql.Stmt
}

func (db *DB) Get(dest any, query string, args ...any) error {
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Ptr {
		return fmt.Errorf("dest must be a pointer")
	}
	destValue = destValue.Elem()

	rows, err := db.Query(query, args...)
	if err != nil {
		return err
	}

	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	if !rows.Next() {
		return sql.ErrNoRows
	}

	columnIndex, err := columnIndexToDestFieldIndex(cols, destValue.Type())
	if err != nil {
		return err
	}

	destValues, err := mapDestValue(destValue, columnIndex)
	if err != nil {
		return err
	}

	if err := rows.Scan(destValues...); err != nil {
		return err
	}

	if err := rows.Close(); err != nil {
		return err
	}

	return nil
}

func (db *DB) Select(destSlice any, query string, args ...any) error {
	destSliceValue := reflect.ValueOf(destSlice)
	if destSliceValue.Kind() != reflect.Ptr {
		return fmt.Errorf("dest must be a pointer to a slice")
	}
	destSliceValue = destSliceValue.Elem()
	if destSliceValue.Kind() != reflect.Slice {
		return fmt.Errorf("dest must be a pointer to a slice")
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return err
	}

	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	destType := destSliceValue.Type().Elem()
	columnIndex, err := columnIndexToDestFieldIndex(cols, destType)
	if err != nil {
		return err
	}

	for rows.Next() {
		newDestValue := reflect.New(destType).Elem()
		destValues, err := mapDestValue(newDestValue, columnIndex)
		if err != nil {
			return err
		}
		if err := rows.Scan(destValues...); err != nil {
			return err
		}
		destSliceValue.Set(reflect.Append(destSliceValue, newDestValue))
	}

	return nil
}

func mapDestValue(destValue reflect.Value, columnIndex []int) ([]any, error) {
	if destValue.Kind() != reflect.Struct {
		fmt.Println(destValue.Kind())
		return []any{destValue.Addr().Interface()}, nil
	}

	destValues := make([]any, len(columnIndex))
	// iterate over cols and add pointer to corresponding struct field to destValues
	for i := range columnIndex {
		field := destValue.Field(columnIndex[i])
		if field.Kind() == reflect.Ptr && field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		if field.Kind() == reflect.Map {
			field.Set(reflect.MakeMap(field.Type()))
		}
		destValues[i] = field.Addr().Interface()
	}

	return destValues, nil
}

func columnIndexToDestFieldIndex(cols []string, destType reflect.Type) ([]int, error) {
	if destType.Kind() != reflect.Struct {
		// TODO: check that is basic type
		// if dest is basic type, then expect only one column in the result set
		if len(cols) != 1 {
			return nil, fmt.Errorf("expected only one column in the result set for basic type dest")
		}
		return []int{0}, nil
	}

	columnIndexToDestFieldIndex := make([]int, len(cols))

	for i := 0; i < destType.NumField(); i++ {
		destField := destType.Field(i)
		if !destField.IsExported() {
			continue
		}

		destFieldName, ok := destField.Tag.Lookup("db")
		if !ok {
			destFieldName = strings.ToLower(destField.Name)
		}
		if destFieldName == "-" {
			continue
		}

		colIndex := slices.Index(cols, destFieldName)
		if colIndex == -1 {
			return nil, fmt.Errorf("there is no column %s in the result set", destFieldName)
		}

		columnIndexToDestFieldIndex[colIndex] = i
	}

	return columnIndexToDestFieldIndex, nil
}
