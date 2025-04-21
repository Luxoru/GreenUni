// Package mysql provides MySQL adapter implementations for repositories and database interactions.
package mysql

import (
	"backend/internal/db"
	"backend/internal/utils/concurrency"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"strconv"
	"strings"
)

//TODO: Add async capabilities

// Container holds the MySQL connection and thread pool. Implements Database
type Container struct {
	database *sql.DB
	pool     *concurrency.ThreadPool
}

// Configurations holds MySQL configuration including auth and DB name.
type Configurations struct {
	Authentication *db.AuthenticationConfigurations
	DatabaseName   string
}

// GetAuthenticationConfigurations returns a copy of the authentication config.
func (config Configurations) GetAuthenticationConfigurations() db.AuthenticationConfigurations {
	return *config.Authentication
}

// Name returns the name of the SQL driver.
func (sqlDatabase *Container) Name() string {
	return "mysql"
}

var emptyConfig = db.AuthenticationConfigurations{}

// Connect initializes the MySQL connection and thread pool.
func (sqlDatabase *Container) Connect(config Configurations) error {

	if *config.Authentication == emptyConfig {
		return fmt.Errorf("Authentication must be provided")
	}

	host := config.Authentication.Host
	port := config.Authentication.Port
	username := config.Authentication.Username
	password := config.Authentication.Password
	dbName := config.DatabaseName

	if host == "" || port < -1 {
		return fmt.Errorf("invalid Connection string. Host: %s, Port: %d", host, port)
	}

	connectionString := username + ":" + password + "@tcp(" + host + ":" + strconv.Itoa(port) + ")/" + dbName + "?parseTime=true" //using parse time so don't have to manually convert to time.Time

	database, err := sql.Open("mysql", connectionString)
	if err != nil {
		return err
	}

	err = database.Ping()

	if err != nil {
		return err
	}

	sqlDatabase.database = database
	sqlDatabase.pool = concurrency.NewThreadPool(500, 100)
	sqlDatabase.pool.Start()

	return nil
}

// GetThreadPool returns the internal thread pool.
func (sqlDatabase *Container) GetThreadPool() *concurrency.ThreadPool {
	return sqlDatabase.pool
}

// Close closes the database connection
func (sqlDatabase *Container) Close() error {
	return sqlDatabase.database.Close()
}

// Repository represents a MySQL repository backed by a Container.
type Repository struct {
	Database *Container
}

// CreateRepository creates a new repository and runs the given table creation query.
func CreateRepository(Database Container, CreateTableQuery string) (*Repository, error) {
	repo := Repository{
		Database: &Database,
	}
	rows, err := repo.ExecuteQuery(CreateTableQuery, make([]Column, 0), QueryOptions{})
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	return &repo, nil
}

// Column represents a database column with a generic type T
type Column interface {
	GetName() string
	GetValue() interface{}
	GetLength() int
	IsNullable() bool
	GetType() string
}

// BaseColumn implements common functionality for all column types
type BaseColumn struct {
	name     string
	value    interface{}
	length   int
	nullable bool
}

// GetName returns the column name
func (c *BaseColumn) GetName() string {
	return c.name
}

// GetValue returns the column value
func (c *BaseColumn) GetValue() interface{} {
	return c.value
}

// GetLength returns the column length
func (c *BaseColumn) GetLength() int {
	return c.length
}

// IsNullable returns whether the column is nullable
func (c *BaseColumn) IsNullable() bool {
	return c.nullable
}

// CharColumn represents a CHAR column in a MySQL table
type CharColumn struct {
	*BaseColumn
}

// NewCharColumn creates a new CharColumn with a value
func NewCharColumn(name string, value rune) *CharColumn {
	return &CharColumn{
		BaseColumn: &BaseColumn{
			name:     name,
			value:    value,
			length:   1,
			nullable: false,
		},
	}
}

// NewCharColumnWithOptions creates a CharColumn with specified options
func NewCharColumnWithOptions(name string, value rune, length int, nullable bool) *CharColumn {
	return &CharColumn{
		BaseColumn: &BaseColumn{
			name:     name,
			value:    value,
			length:   length,
			nullable: nullable,
		},
	}
}

// GetType returns the MySQL column type
func (c *CharColumn) GetType() string {
	return "CHAR"
}

// IntegerColumn represents an INT column in a MySQL table
type IntegerColumn struct {
	*BaseColumn
	AutoIncrement bool
}

// NewIntegerColumn creates a new IntegerColumn with a value
func NewIntegerColumn(name string, value int64) *IntegerColumn {
	return &IntegerColumn{
		BaseColumn: &BaseColumn{
			name:     name,
			value:    value,
			length:   1,
			nullable: false,
		},
		AutoIncrement: false,
	}
}

// NewIntegerColumnForTable creates an IntegerColumn for table definition
func NewIntegerColumnForTable(name string, nullable bool, autoIncrement bool) *IntegerColumn {
	return &IntegerColumn{
		BaseColumn: &BaseColumn{
			name:     name,
			value:    nil,
			length:   1,
			nullable: nullable,
		},
		AutoIncrement: autoIncrement,
	}
}

// GetType returns the MySQL column type
func (c *IntegerColumn) GetType() string {
	return "INT"
}

// IsAutoIncrement returns whether the column is auto-incrementing
func (c *IntegerColumn) IsAutoIncrement() bool {
	return c.AutoIncrement
}

// UUIDColumn represents a UUID stored as VARCHAR(36)
type UUIDColumn struct {
	*BaseColumn
}

// NewUUIDColumn creates a new UUIDColumn with a value
func NewUUIDColumn(name string, value uuid.UUID) *UUIDColumn {
	return &UUIDColumn{
		BaseColumn: &BaseColumn{
			name:     name,
			value:    value,
			length:   1,
			nullable: false,
		},
	}
}

// NewUUIDColumnForTable creates a UUIDColumn for table definition
func NewUUIDColumnForTable(name string, nullable bool, length int) *UUIDColumn {
	return &UUIDColumn{
		BaseColumn: &BaseColumn{
			name:     name,
			value:    nil,
			length:   length,
			nullable: nullable,
		},
	}
}

// GetType returns the MySQL column type
func (c *UUIDColumn) GetType() string {
	return "VARCHAR(36)"
}

// VarcharColumn represents a VARCHAR column in a MySQL table
type VarcharColumn struct {
	*BaseColumn
}

// NewVarcharColumn creates a new VarcharColumn with a value
func NewVarcharColumn(name string, value string) *VarcharColumn {
	return &VarcharColumn{
		BaseColumn: &BaseColumn{
			name:     name,
			value:    value,
			length:   1,
			nullable: false,
		},
	}
}

// NewVarcharColumnForTable creates a VarcharColumn for table definition
func NewVarcharColumnForTable(name string, nullable bool, length int) *VarcharColumn {
	return &VarcharColumn{
		BaseColumn: &BaseColumn{
			name:     name,
			value:    nil,
			length:   length,
			nullable: nullable,
		},
	}
}

// GetType returns the MySQL column type
func (c *VarcharColumn) GetType() string {
	return "VARCHAR"
}

// TextColumn represents a TEXT column in a MySQL table
type TextColumn struct {
	*BaseColumn
}

// NewTextColumn creates a new TextColumn with a value
func NewTextColumn(name string, value string) *TextColumn {
	return &TextColumn{
		BaseColumn: &BaseColumn{
			name:     name,
			value:    value,
			length:   1,
			nullable: false,
		},
	}
}

// GetType returns the MySQL column type
func (c *TextColumn) GetType() string {
	return "TEXT"
}

// Table defines the schema of a MySQL table
type Table struct {
	Name        string
	Columns     *[]Column
	PrimaryKeys *[]string
}

// GetOrCreateTableQuery generates a SQL query to create a table
func (t *Table) GetOrCreateTableQuery(createIfExists bool) (string, error) {
	if len(*t.Columns) < 1 {
		return "", fmt.Errorf("columns has to be greater or equal to 1")
	}

	createPrefix := "CREATE TABLE "
	if !createIfExists {
		createPrefix += "IF NOT EXISTS "
	}

	sb := strings.Builder{}
	sb.WriteString(createPrefix)
	sb.WriteString("`")
	sb.WriteString(t.Name)
	sb.WriteString("`(")

	columns := *t.Columns

	for _, column := range columns {
		sb.WriteString(column.GetName())
		sb.WriteString(" ")
		sb.WriteString(column.GetType())
		sb.WriteString(" ")

		columnLength := column.GetLength()
		if columnLength > 1 {
			sb.WriteString(fmt.Sprintf("(%d)", columnLength))
		}

		if intColumn, ok := column.(*IntegerColumn); ok {
			if intColumn.IsAutoIncrement() {
				if intColumn.IsNullable() {
					return "", fmt.Errorf("auto incrementing cannot be null")
				}
				sb.WriteString(" AUTO_INCREMENT")
			}
		}

		if !column.IsNullable() {
			sb.WriteString(" NOT NULL")
		}

		sb.WriteString(",")
	}

	primaryKeys := *t.PrimaryKeys

	if len(primaryKeys) > 0 {
		sb.WriteString("PRIMARY KEY (")
		for i, key := range primaryKeys {
			sb.WriteString(key)
			if i < len(primaryKeys)-1 {
				sb.WriteString(", ")
			}
		}
		sb.WriteString(")")
	} else {
		// Remove the trailing comma if there are no primary keys
		query := sb.String()
		query = query[:len(query)-1]
		sb = strings.Builder{}
		sb.WriteString(query)
	}

	sb.WriteString(");")
	return sb.String(), nil
}

// InsertOptions represents options for insert execution.
type InsertOptions struct {
	OnComplete func(result sql.Result)
	OnError    func(error)
	connection *sql.DB
}

// NewInsertOptions creates a new InsertOptions object.
func NewInsertOptions(OnComplete func(res sql.Result), onError func(err error), connection *sql.DB) *InsertOptions {
	return &InsertOptions{
		OnComplete: OnComplete,
		OnError:    onError,
		connection: connection,
	}
}

// QueryOptions represents options for query execution.
type QueryOptions struct {
	OnComplete func(*sql.Rows)
	OnError    func(error)
	Connection *sql.DB
}

// NewQueryOptions creates a new QueryOptions object.
func NewQueryOptions(OnComplete func(rows *sql.Rows), onError func(err error), connection *sql.DB) *QueryOptions {
	return &QueryOptions{
		OnComplete: OnComplete,
		OnError:    onError,
		Connection: connection,
	}
}

func (r *Repository) StartTransaction() (*sql.Tx, error) {
	return r.Database.database.Begin()
}

func (r *Repository) AddExecuteTransaction(tx *sql.Tx, query string, columns []Column) (sql.Result, error) {
	questionMarkCount := strings.Count(query, "?")
	if questionMarkCount != len(columns) {
		return nil, fmt.Errorf("invalid amount of columns for query \"%s\"", query)
	}

	args := make([]interface{}, 0)
	for _, column := range columns {
		args = append(args, column.GetValue())
	}

	return tx.Exec(query, args...)
}

func (r *Repository) CommitTransaction(tx *sql.Tx) error {
	return tx.Commit()
}

// ExecuteInsert executes a SQL insert with the given columns and returns the number of affected rows
func (r *Repository) ExecuteInsert(query string, columns []Column, options InsertOptions) (int64, error) {
	questionMarkCount := strings.Count(query, "?")
	if questionMarkCount != len(columns) {
		return 0, fmt.Errorf("invalid amount of columns for query \"%s\"", query)
	}

	database := options.connection
	if database == nil {
		database = r.Database.database
	}

	stmt, err := database.Prepare(query)
	if err != nil {
		if options.OnError != nil {
			options.OnError(err)
		}
		return 0, err
	}
	defer stmt.Close()

	args := make([]interface{}, 0)
	for _, column := range columns {
		args = append(args, column.GetValue())
	}

	result, err := stmt.Exec(args...)
	if err != nil {
		if options.OnError != nil {
			options.OnError(err)
		}
		return 0, err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		if options.OnError != nil {
			options.OnError(err)
		}
		return 0, err
	}

	if options.OnComplete != nil {
		options.OnComplete(result)
	}

	return affectedRows, nil
}

// ExecuteQuery executes a SQL query and returns the result rows.
func (r *Repository) ExecuteQuery(query string, columns []Column, options QueryOptions) (*sql.Rows, error) {
	if columns != nil {
		questionMarkCount := strings.Count(query, "?")
		if questionMarkCount != len(columns) {
			return nil, fmt.Errorf("invalid amount of columns for query \"%s\"", query)
		}
	}

	database := options.Connection
	if database == nil {
		if r.Database.database == nil {
			return nil, nil
		}
		database = r.Database.database
	}

	stmt, err := database.Prepare(query)
	if err != nil {
		if options.OnError != nil {
			options.OnError(err)
		}
		return nil, err
	}
	defer stmt.Close()

	args := make([]interface{}, 0)
	if columns != nil {
		for _, column := range columns {
			args = append(args, column.GetValue())
		}
	}
	rows, err := stmt.Query(args...)
	if err != nil {
		if options.OnError != nil {
			options.OnError(err)
		}
		return nil, err
	}

	if options.OnComplete != nil {
		options.OnComplete(rows)
	}

	return rows, nil
}
