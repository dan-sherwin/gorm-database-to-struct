package pgtypes

import (
	"gorm.io/gorm"
)

// DataTypeMap returns a mapping from PostgreSQL types to Go types
// for use with GORM Gen's WithDataTypeMap.
func DataTypeMap() map[string]func(columnType gorm.ColumnType) string {
	return map[string]func(columnType gorm.ColumnType) string{
		"text[]":                        use("pgtypes.StringArray"),
		"varchar[]":                     use("pgtypes.StringArray"),
		"integer[]":                     use("pgtypes.Int32Array"),
		"int4[]":                        use("pgtypes.Int32Array"),
		"int8[]":                        use("pgtypes.Int64Array"),
		"bigint[]":                      use("pgtypes.Int64Array"),
		"bool[]":                        use("pgtypes.BoolArray"),
		"boolean[]":                     use("pgtypes.BoolArray"),
		"uuid[]":                        use("pgtypes.UUIDArray"),
		"float8[]":                      use("pgtypes.Float64Array"),
		"double precision[]":            use("pgtypes.Float64Array"),
		"timestamptz[]":                 use("pgtypes.TimeArray"),
		"timestamp[]":                   use("pgtypes.TimeArray"),
		"timestamp with time zone[]":    use("pgtypes.TimeArray"),
		"timestamp without time zone[]": use("pgtypes.TimeArray"),
		"interval":                      use("pgtypes.Duration"),
		"interval[]":                    use("pgtypes.DurationArray"),
	}
}

func use(goType string) func(columnType gorm.ColumnType) string {
	return func(columnType gorm.ColumnType) string {
		return goType
	}
}
