package schema

import (
	"context"
	"fmt"

	"fbc/ent/dialect"
	"fbc/ent/dialect/sql"
)

// MySQL is a mysql migration driver.
type MySQL struct {
	dialect.Driver
	version string
}

// init loads the MySQL version from the database for later use in the migration process.
func (d *MySQL) init(ctx context.Context, tx dialect.Tx) error {
	rows := &sql.Rows{}
	if err := tx.Query(ctx, "SHOW VARIABLES LIKE 'version'", []interface{}{}, rows); err != nil {
		return fmt.Errorf("mysql: querying mysql version %v", err)
	}
	defer rows.Close()
	if !rows.Next() {
		return fmt.Errorf("mysql: version variable was not found")
	}
	version := make([]string, 2)
	if err := rows.Scan(&version[0], &version[1]); err != nil {
		return fmt.Errorf("mysql: scanning mysql version: %v", err)
	}
	d.version = version[1]
	return nil
}

func (d *MySQL) tableExist(ctx context.Context, tx dialect.Tx, name string) (bool, error) {
	query, args := sql.Select(sql.Count("*")).From(sql.Table("INFORMATION_SCHEMA.TABLES").Unquote()).
		Where(sql.EQ("TABLE_SCHEMA", sql.Raw("(SELECT DATABASE())")).And().EQ("TABLE_NAME", name)).Query()
	return exist(ctx, tx, query, args...)
}

func (d *MySQL) fkExist(ctx context.Context, tx dialect.Tx, name string) (bool, error) {
	query, args := sql.Select(sql.Count("*")).From(sql.Table("INFORMATION_SCHEMA.TABLE_CONSTRAINTS").Unquote()).
		Where(sql.EQ("TABLE_SCHEMA", sql.Raw("(SELECT DATABASE())")).And().EQ("CONSTRAINT_TYPE", "FOREIGN KEY").And().EQ("CONSTRAINT_NAME", name)).Query()
	return exist(ctx, tx, query, args...)
}

// table loads the current table description from the database.
func (d *MySQL) table(ctx context.Context, tx dialect.Tx, name string) (*Table, error) {
	rows := &sql.Rows{}
	query, args := sql.Select("column_name", "column_type", "is_nullable", "column_key", "column_default", "extra", "character_set_name", "collation_name").
		From(sql.Table("INFORMATION_SCHEMA.COLUMNS").Unquote()).
		Where(sql.EQ("TABLE_SCHEMA", sql.Raw("(SELECT DATABASE())")).And().EQ("TABLE_NAME", name)).Query()
	if err := tx.Query(ctx, query, args, rows); err != nil {
		return nil, fmt.Errorf("mysql: reading table description %v", err)
	}
	defer rows.Close()
	t := &Table{Name: name}
	for rows.Next() {
		c := &Column{}
		if err := c.ScanMySQL(rows); err != nil {
			return nil, fmt.Errorf("mysql: %v", err)
		}
		if c.PrimaryKey() {
			t.PrimaryKey = append(t.PrimaryKey, c)
		}
		t.Columns = append(t.Columns, c)
	}
	return t, nil
}

func (d *MySQL) setRange(ctx context.Context, tx dialect.Tx, name string, value int) error {
	return tx.Exec(ctx, fmt.Sprintf("ALTER TABLE `%s` AUTO_INCREMENT = %d", name, value), []interface{}{}, new(sql.Result))
}

func (d *MySQL) cType(c *Column) string                { return c.MySQLType(d.version) }
func (d *MySQL) tBuilder(t *Table) *sql.TableBuilder   { return t.MySQL(d.version) }
func (d *MySQL) cBuilder(c *Column) *sql.ColumnBuilder { return c.MySQL(d.version) }