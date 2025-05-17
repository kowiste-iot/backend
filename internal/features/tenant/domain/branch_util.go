package domain

import (
	"bytes"
	"embed"
	"html/template"
	"io"
)

//go:embed branch.sql
var schema embed.FS

func GetBranchSchemaSQL(schemaName string, version int) ([]byte, error) {
	t, err := template.ParseFS(schema, "branch.sql")
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)

	err = t.Execute(buf, struct {
		SchemaName string
	}{
		SchemaName: schemaName,
	})
	if err != nil {
		return nil, err
	}

	sql, err := io.ReadAll(buf)
	if err != nil {
		return nil, err
	}

	return sql, nil
}