package models

import "gorm.io/cli/gorm/genconfig"

var _ = genconfig.Config{
	IncludeInterfaces: []any{"Query*"},
	IncludeStructs:    []any{User{}, Article{}, ArticleCategory{}, ArticleComment{}, Notice{}},
}

type Query[T any] interface {
	// SELECT * FROM @@table WHERE id=@id
	GetByID(id int64) (T, error)

	// SELECT * FROM @@table WHERE @@column=@value
	FilterWithColumn(column string, value string) (T, error)

	// SELECT * FROM @@table WHERE @@column LIKE @value
	FilterWithColumnLike(column string, value string) ([]T, error)
}
