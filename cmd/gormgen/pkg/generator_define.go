package pkg

// fieldConfig
type fieldConfig struct {
	FieldName  string
	ColumnName string
	FieldType  string
	HumpName   string
}

// structConfig
type structConfig struct {
	config
	StructName   string
	OnlyFields   []fieldConfig
	OptionFields []fieldConfig
	PkField      fieldConfig
	TemplateName string
}

type ImportPkg struct {
	Pkg string
}

type structHelpers struct {
	Titelize func(string) string
}

type config struct {
	PkgName          string
	Helpers          structHelpers
	QueryBuilderName string
}
