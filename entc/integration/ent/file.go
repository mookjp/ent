// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"bytes"
	"fmt"
	"strconv"

	"fbc/ent/dialect/gremlin"
	"fbc/ent/dialect/sql"
)

// File is the model entity for the File schema.
type File struct {
	config
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// Size holds the value of the "size" field.
	Size int `json:"size,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
}

// FromRows scans the sql response data into File.
func (f *File) FromRows(rows *sql.Rows) error {
	var vf struct {
		ID   int
		Size int
		Name string
	}
	// the order here should be the same as in the `file.Columns`.
	if err := rows.Scan(
		&vf.ID,
		&vf.Size,
		&vf.Name,
	); err != nil {
		return err
	}
	f.ID = strconv.Itoa(vf.ID)
	f.Size = vf.Size
	f.Name = vf.Name
	return nil
}

// FromResponse scans the gremlin response data into File.
func (f *File) FromResponse(res *gremlin.Response) error {
	vmap, err := res.ReadValueMap()
	if err != nil {
		return err
	}
	var vf struct {
		ID   string `json:"id,omitempty"`
		Size int    `json:"size,omitempty"`
		Name string `json:"name,omitempty"`
	}
	if err := vmap.Decode(&vf); err != nil {
		return err
	}
	f.ID = vf.ID
	f.Size = vf.Size
	f.Name = vf.Name
	return nil
}

// Update returns a builder for updating this File.
// Note that, you need to call File.Unwrap() before calling this method, if this File
// was returned from a transaction, and the transaction was committed or rolled back.
func (f *File) Update() *FileUpdateOne {
	return (&FileClient{f.config}).UpdateOne(f)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (f *File) Unwrap() *File {
	tx, ok := f.config.driver.(*txDriver)
	if !ok {
		panic("ent: File is not a transactional entity")
	}
	f.config.driver = tx.drv
	return f
}

// String implements the fmt.Stringer.
func (f *File) String() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("File(")
	buf.WriteString(fmt.Sprintf("id=%v", f.ID))
	buf.WriteString(fmt.Sprintf(", size=%v", f.Size))
	buf.WriteString(fmt.Sprintf(", name=%v", f.Name))
	buf.WriteString(")")
	return buf.String()
}

// id returns the int representation of the ID field.
func (f *File) id() int {
	id, _ := strconv.Atoi(f.ID)
	return id
}

// Files is a parsable slice of File.
type Files []*File

// FromRows scans the sql response data into Files.
func (f *Files) FromRows(rows *sql.Rows) error {
	for rows.Next() {
		vf := &File{}
		if err := vf.FromRows(rows); err != nil {
			return err
		}
		*f = append(*f, vf)
	}
	return nil
}

// FromResponse scans the gremlin response data into Files.
func (f *Files) FromResponse(res *gremlin.Response) error {
	vmap, err := res.ReadValueMap()
	if err != nil {
		return err
	}
	var vf []struct {
		ID   string `json:"id,omitempty"`
		Size int    `json:"size,omitempty"`
		Name string `json:"name,omitempty"`
	}
	if err := vmap.Decode(&vf); err != nil {
		return err
	}
	for _, v := range vf {
		*f = append(*f, &File{
			ID:   v.ID,
			Size: v.Size,
			Name: v.Name,
		})
	}
	return nil
}

func (f Files) config(cfg config) {
	for i := range f {
		f[i].config = cfg
	}
}