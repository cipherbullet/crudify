package main

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sql.DB
}

func (p *PostgresStorage) Connect(connStr string) error {
	// db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	// if err != nil {
	// 	panic("failed to connect database")
	// }

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		return err
	}

	p.db = db
	return nil
}

func (p *PostgresStorage) Disconnect() error {
	if p.db != nil {
		return p.db.Close()
	}
	return nil
}

func (p *PostgresStorage) All(ctx context.Context, modelPtr interface{}) error {
	destType := reflect.TypeOf(modelPtr)
	if destType.Kind() != reflect.Ptr || destType.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("destination must be a pointer to a slice")
	}

	modelElem := destType.Elem().Elem()
	tableName := getTableName(modelElem)

	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return err
	}
	defer rows.Close()

	sliceValue := reflect.ValueOf(modelPtr).Elem()
	for rows.Next() {
		elem := reflect.New(modelElem)

		// Before calling rows.Scan, create separate pointers for each field in the struct
		fields := make([]interface{}, 0)
		for i := 0; i < elem.Elem().NumField(); i++ {
			fieldPtr := elem.Elem().Field(i).Addr().Interface()
			fields = append(fields, fieldPtr)
		}
		if err := rows.Scan(fields...); err != nil {
			return err
		}
		sliceValue.Set(reflect.Append(sliceValue, elem.Elem()))
	}
	return rows.Err()
}

func (p *PostgresStorage) Find(ctx context.Context, modelPtr interface{}, id int) error {
	modelType := reflect.TypeOf(modelPtr).Elem()
    modelElem := reflect.ValueOf(modelPtr).Elem()
	// if destType.Kind() != reflect.Ptr || destType.Elem().Kind() != reflect.Slice {
	//     return fmt.Errorf("destination must be a pointer to a slice")
	// }

	tableName := getTableName(modelType)

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", tableName)

	row := p.db.QueryRowContext(ctx, query, id)

	fields := make([]interface{}, 0)
	for i := 0; i < modelElem.NumField(); i++ {
		fieldPtr := modelElem.Field(i).Addr().Interface()
		fields = append(fields, fieldPtr)
	}

	err := row.Scan(fields...)
	if err != nil {
		return err
	}

	return nil
}

// Create function to insert a new record into the database
func (p *PostgresStorage) Create(ctx context.Context, modelPtr interface{}) error {
	modelType := reflect.TypeOf(modelPtr)
	if modelType.Kind() != reflect.Ptr || modelType.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("model must be a pointer to a struct")
	}

	// Get the model type and table name
	elemType := modelType.Elem()
	tableName := getTableName(elemType)

	var columns []string
	var placeholders []string
	var values []interface{}

	// Iterate over the struct fields
	elemValue := reflect.ValueOf(modelPtr).Elem()
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)

		columnName := field.Tag.Get("db")

		// Skip non-exported fields or primary keys if needed
		if field.PkgPath != "" || columnName == "id" {
			continue
		}

		columns = append(columns, columnName)
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(placeholders)+1))
		values = append(values, elemValue.Field(i).Interface())
	}

	// Construct the INSERT query
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, strings.Join(columns, ", "), strings.Join(placeholders, ", "))

    var insertedID int64
    err := p.db.QueryRowContext(ctx, query, values...).Scan(&insertedID)
	// _, err := p.db.ExecContext(ctx, query, values...)
	if err != nil {
        if err == sql.ErrNoRows {
            // No rows were returned, but that's expected after an insert operation
            // You can manually retrieve the last inserted ID, considering a successful insertion
            insertedID, _ = p.getLastInsertedID(ctx)
        } else {
            return err
        }
    }

    // Update the struct with the inserted ID
    idField := elemValue.FieldByName("Id")
    fmt.Println("idField", idField)
    if idField.IsValid() && idField.CanSet() {
        idField.SetInt(insertedID)
    }

	return nil
}

func (p *PostgresStorage) Delete(ctx context.Context, modelPtr interface{}) error {
    modelValue := reflect.Indirect(reflect.ValueOf(modelPtr))

    // Get the table name and primary key field
    tableName := getTableName(reflect.TypeOf(modelPtr).Elem())
    idField := modelValue.FieldByName("Id")
    if !idField.IsValid() {
        return fmt.Errorf("ID field not found in struct")
    }

    // Construct the DELETE query
    query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", tableName)

    // Execute the DELETE query
    _, err := p.db.ExecContext(ctx, query, idField.Interface())
    if err != nil {
        return err
    }

    return nil
}

func (p *PostgresStorage) getLastInsertedID(ctx context.Context) (int64, error) {
    var lastInsertID int64
    err := p.db.QueryRowContext(ctx, "SELECT lastval()").Scan(&lastInsertID)
    if err != nil {
        return 0, err
    }
    return lastInsertID, nil
}

func getTableName(t reflect.Type) string {
	// Use reflection or struct tags to determine the table name.
	// This example assumes the table name matches the struct name in lowercase.
	return t.Name()
}

func (p *PostgresStorage) generateSelectQuery(model reflect.Value) (string, error) {

	var fields []string
	for i := 0; i < reflect.TypeOf(model).Elem().NumField(); i++ {
		field := reflect.TypeOf(model).Elem().Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag != "" {
			fields = append(fields, dbTag)
		}
	}

	if len(fields) == 0 {
		return "", fmt.Errorf("no db tags found in model")
	}

	// query := fmt.Sprintf("SELECT %s FROM %s", strings.Join(fields, ", "), p.tableName(p.tableName(model)))
	// return query, nil

	return "", nil
}
