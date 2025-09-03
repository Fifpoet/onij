package main

import (
	"bytes"
	"context"
	"fmt"
	"onij/util/boost/collection"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/iancoleman/strcase"
	"github.com/xwb1989/sqlparser"
)

var (
	logDebug  = color.New(color.FgBlue).Printf
	logInfo   = color.New(color.FgGreen).Printf
	logError  = color.New(color.FgRed).Printf
	logNotice = color.New(color.FgYellow).Printf
)

/*
func (f *FilePages) Value() (driver.Value, error) {
	if f == nil {
		return "", nil
	}
	buf, err := json.Marshal(f)
	if err != nil {
		return nil, err
	}
	return string(buf), nil
}

func (f *FilePages) Scan(value interface{}) error {
	if f == nil {
		return nil
	}

	var b []byte
	switch v := value.(type) {
	case string:
		b = []byte(v)
	case []byte:
		b = v
	default:
		return nil
	}

	if string(b) == "" {
		return nil
	}

	return json.Unmarshal(b, &f)
}
*/

type SpecType int32

type generator struct {
	DestPath string
	DdlPath  string
}

func (g *generator) Run() {
	_, _ = logNotice("==================== Generate Database Model ====================\n")
	_, _ = logNotice("DDLPath: %s\n", g.DdlPath)
	_, _ = logNotice("ModelDir: %s\n\n", g.DestPath)

	_, _ = logInfo("[INFO] Prepare to get DDL content...\n")
	ddlContent, err := g.getDDLContent()
	if err != nil {
		_, _ = logError("[ERROR] Get DDL content failed: %v\n", err)
		return
	}
	_, _ = logDebug("%s\n\n", string(ddlContent))

	_, _ = logInfo("[INFO] Parsing DDL...\n\n")
	ddl, err := g.parseDDL(ddlContent)
	if err != nil {
		_, _ = logError("[ERROR] Parse DDL failed: %v\n", err)
		return
	}

	_, _ = logInfo("[INFO] Generating model file...\n\n")
	err = g.writeFile(ddl)
	if err != nil {
		_, _ = logError("[ERROR] Generate model file failed: %v\n", err)
		return
	}

	_, _ = logNotice("==================== Generate Database Model Done ====================\n")
}

func (g *generator) writeFile(ddl *sqlparser.DDL) error {
	fName, content := g.genFileContent(ddl)
	if err := os.WriteFile(fName, content, 0644); err != nil {
		return err
	}
	cmd := exec.Command("gofmt", "-w", fName)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (g *generator) genFileContent(ddl *sqlparser.DDL) (fileName string, fileContent []byte) {
	tpl := `package model

/*********** 表字段常量定义 **********/

%s

/*********** 表结构定义 **********/
%s

%s

/*********** 查询器与更新器 **********/

%s

%s

`
	colConstStr := g.genColumnNameConsts(ddl)
	tableStruct := g.genTableModel(ddl)
	tableMethods := g.genTableModelMethods(ddl)
	tableQuerier := g.genTableQuerier(ddl)
	tableUpdater := g.genTableUpdater(ddl)

	fileName = filepath.Join(g.DestPath, fmt.Sprintf("%s_model_gen.go", ddl.NewName.Name.String()))
	fileContent = []byte(fmt.Sprintf(tpl, colConstStr, tableStruct, tableMethods, tableQuerier, tableUpdater))
	return
}

func (g *generator) genTableUpdater(ddl *sqlparser.DDL) string {
	var buf bytes.Buffer
	modelName := strcase.ToCamel(ddl.NewName.Name.String())
	updaterName := fmt.Sprintf("%sUpdater", modelName)

	buf.WriteString(fmt.Sprintf("// %s 表 %s 更新器\n", updaterName, ddl.NewName.Name.String()))
	buf.WriteString(fmt.Sprintf("type %s cdb.BasicUpdater\n\n", updaterName))
	// buf.WriteString(fmt.Sprintf("func New%s() %s {\nreturn map[string]any{}\n}\n\n", updaterName, updaterName))
	buf.WriteString(fmt.Sprintf("func (u %s) super() cdb.BasicUpdater {\nreturn (cdb.BasicUpdater)(u)\n}\n\n", updaterName))
	for _, col := range ddl.TableSpec.Columns {
		colName := col.Name.String()             // DB列名
		fieldName := strcase.ToCamel(colName)    // 字段名
		typeName := MysqlTypeToGo[col.Type.Type] // 类型名

		buf.WriteString(fmt.Sprintf("func (u %s) %s(v %s) %s {\nu.super().Add(%s_%s, v)\nreturn u\n}\n\n",
			updaterName, fieldName, typeName, updaterName, modelName, fieldName))
	}

	buf.WriteString(fmt.Sprintf("func (u %s) ToMap() map[string]any {\nreturn u.super().ToMap()\n}\n\n", updaterName))
	buf.WriteString(fmt.Sprintf("func (u %s) IsEmpty() bool {\n return u.super().IsEmpty() \n}\n\n", updaterName))
	return buf.String()
}

func (g *generator) genTableQuerier(ddl *sqlparser.DDL) string {
	var buf bytes.Buffer
	modelName := strcase.ToCamel(ddl.NewName.Name.String())
	querierName := fmt.Sprintf("%sQuerier", modelName)

	buf.WriteString(fmt.Sprintf("// %s 表 %s 查询器\n", querierName, ddl.NewName.Name.String()))
	buf.WriteString(fmt.Sprintf("type %s cdb.BasicQuerier\n\n", querierName))
	// buf.WriteString(fmt.Sprintf("func New%s() *%s {\nreturn &%s{map[string]any{}}\n}\n\n", querierName, querierName, querierName))
	buf.WriteString(fmt.Sprintf("func (q %s) super() cdb.BasicQuerier {\nreturn (cdb.BasicQuerier)(q)\n}\n\n", querierName))
	buf.WriteString(fmt.Sprintf("func (q %s) ToOptions() []cdb.Option {\nreturn q.super().ToOptions()\n}\n\n", querierName))
	buf.WriteString(fmt.Sprintf("func (q %s) WithAsc(cols ...string) %s {\nq.super().WithAsc(cols...)\nreturn q\n}\n\n", querierName, querierName))
	buf.WriteString(fmt.Sprintf("func (q %s) WithDesc(cols ...string) %s {\nq.super().WithDesc(cols...)\nreturn q\n}\n\n", querierName, querierName))
	buf.WriteString(fmt.Sprintf("func (q %s) WithLimit(limit int) %s {\nq.super().WithLimit(limit)\nreturn q\n}\n\n", querierName, querierName))
	buf.WriteString(fmt.Sprintf("func (q %s) WithUnscoped() %s {\nq.super().WithUnscoped()\nreturn q\n}\n\n", querierName, querierName))

	colMap := map[string]*sqlparser.ColumnDefinition{}
	for _, col := range ddl.TableSpec.Columns {
		colMap[col.Name.String()] = col
	}

	// 生成索引函数
	for _, idx := range ddl.TableSpec.Indexes {
		if idx.Info.Primary {
			continue
		}
		funcName := strcase.ToCamel(idx.Info.Name.String())

		if !idx.Info.Unique && !strings.HasPrefix(funcName, "Idx") {
			funcName = "Idx" + funcName
		}
		if idx.Info.Unique && !strings.HasPrefix(funcName, "Uniq") && !strings.HasPrefix(funcName, "Uk") {
			funcName = "Uniq" + funcName
		}

		var argsStr, bodyStr bytes.Buffer
		for i, col := range idx.Columns {
			colName := col.Column.String() // DB列名
			colDef := colMap[colName]
			fieldName := strcase.ToCamel(colName) // 字段名
			argName := strings.ToLower(fieldName[:1]) + fieldName[1:]
			typeName := MysqlTypeToGo[colDef.Type.Type]

			if golangKeywords.Contains(argName) {
				argName += "Field"
			}

			argsStr.WriteString(fmt.Sprintf("%s %s", argName, typeName))
			bodyStr.WriteString(fmt.Sprintf("%s(%s)", fieldName, argName))
			if i != len(idx.Columns)-1 {
				argsStr.WriteString(",")
				bodyStr.WriteString(".")
			}
		}

		buf.WriteString(fmt.Sprintf("func (q %s) %s(%s) %s {\n return q.%s\n}\n\n", querierName, funcName, argsStr.String(), querierName, bodyStr.String()))
	}

	for _, col := range ddl.TableSpec.Columns {
		colName := col.Name.String()          // DB列名
		fieldName := strcase.ToCamel(colName) // 字段名

		// 生成查询方法
		buf.WriteString(fmt.Sprintf("func (q %s) %s(v any) %s {\n", querierName, fieldName, querierName))
		buf.WriteString(fmt.Sprintf("q.super().Add(%s_%s, v)\nreturn q\n}\n\n", modelName, fieldName))
	}

	return buf.String()
}

func (g *generator) genTableModelMethods(ddl *sqlparser.DDL) string {
	var buf bytes.Buffer
	modelName := strcase.ToCamel(ddl.NewName.Name.String())

	// 生成构造函数
	buf.WriteString(fmt.Sprintf("func New%s() *%s {\nreturn &%s{}\n}\n\n", modelName, modelName, modelName))

	// 生成 table name
	buf.WriteString(fmt.Sprintf("func (%s) TableName() string {\n return TableName_%s\n}\n\n", modelName, modelName))

	ptCols := collection.NewSet[string]()
	for _, idx := range ddl.TableSpec.Indexes {
		if !idx.Info.Primary && !idx.Info.Unique {
			continue
		}
		for _, col := range idx.Columns {
			colName := col.Column.String() // DB列名
			ptCols.Add(colName)
		}
	}

	// 生成 UpdatableColumns
	buf.WriteString(fmt.Sprintf("func (%s) UpdatableColumns() []string {\n return []string{\n", modelName))
	for _, col := range ddl.TableSpec.Columns {
		colName := col.Name.String() // DB列名
		if ptCols.Contains(colName) {
			continue
		}
		buf.WriteString(fmt.Sprintf("%s_%s,\n", modelName, strcase.ToCamel(colName)))
	}
	buf.WriteString("}\n}\n\n")

	// 生成 primary name
	// for _, idx := range ddl.TableSpec.Indexes {
	// 	if !idx.Info.Primary {
	// 		continue
	// 	}
	// 	col := idx.Columns[0]
	// 	colName := col.Column.String()        // DB列名
	// 	fieldName := strcase.ToCamel(colName) // 字段名
	//
	// 	buf.WriteString(fmt.Sprintf("func (%s) PrimaryKeyName() string {\n return %s_%s\n}\n\n", modelName, modelName, fieldName))
	// 	buf.WriteString(fmt.Sprintf("func (m %s) GetPrimaryKey() int64 {\n return m.%s}\n\n", modelName, fieldName))
	// 	break
	// }

	return buf.String()
}

func (g *generator) genTableModel(ddl *sqlparser.DDL) string {
	var buf bytes.Buffer
	modelName := strcase.ToCamel(ddl.NewName.Name.String())

	buf.WriteString(fmt.Sprintf("\n// %s 表 %s 结构定义\n", modelName, ddl.NewName.Name.String()))
	buf.WriteString(fmt.Sprintf("type %s struct {\n", modelName))
	for _, col := range ddl.TableSpec.Columns {
		colName := col.Name.String()             // DB列名
		fieldName := strcase.ToCamel(colName)    // 字段名
		typeName := MysqlTypeToGo[col.Type.Type] // 类型名
		if typeName == "time.Time" && colName == "deleted_at" {
			typeName = "gorm.DeletedAt"
		}

		// 生成字段声名
		buf.WriteString(fmt.Sprintf("%s %s", fieldName, typeName))

		// 生成 gorm tag
		buf.WriteString(fmt.Sprintf("`gorm:\"column:%s\" json:\"%s\"`", colName, colName))

		// 生成 字段注释
		if col.Type.Comment != nil {
			buf.Write(append([]byte("// "), col.Type.Comment.Val...))
		}

		// 换行
		buf.WriteString("\n")
	}
	buf.WriteString("}\n\n")

	return buf.String()
}

func (g *generator) genColumnNameConsts(ddl *sqlparser.DDL) string {
	var buf bytes.Buffer
	modelName := strcase.ToCamel(ddl.NewName.Name.String())
	hasTime := false
	hasDeleted := false

	buf.WriteString("const (\n")
	buf.WriteString(fmt.Sprintf("TableName_%s = \"%s\"\n\n", modelName, ddl.NewName.Name.String()))
	for _, col := range ddl.TableSpec.Columns {
		colName := col.Name.String()             // DB列名
		fieldName := strcase.ToCamel(colName)    // 字段名
		typeName := MysqlTypeToGo[col.Type.Type] // 类型名

		hasTime = hasTime || typeName == "time.Time"
		hasDeleted = hasDeleted || colName == "deleted_at"

		buf.WriteString(fmt.Sprintf(`%s_%s = "%s"`, modelName, fieldName, colName))
		buf.WriteString("\n")
	}
	buf.WriteString(")\n")

	if !hasDeleted && !hasTime {
		return buf.String()
	}

	prefix := "\nimport (\n"
	prefix += "cdb \"code.chenji.com/pkg/common/component/db\"\n"
	if hasTime {
		prefix += "\"time\"\n"
	}
	if hasDeleted {
		prefix += `"gorm.io/gorm"`
	}
	prefix += "\n)\n\n"

	return prefix + buf.String()
}

func (g *generator) parseDDL(ddlContent []byte) (*sqlparser.DDL, error) {
	stmt, err := sqlparser.ParseStrictDDL(string(ddlContent))
	if err != nil {
		return nil, err
	}
	ddl, ok := stmt.(*sqlparser.DDL)
	if !ok {
		return nil, fmt.Errorf("invalid ddl file content, parse to DDL failed")
	}
	return ddl, nil
}

func (g *generator) getDDLContent() ([]byte, error) {
		ddlContent, err := os.ReadFile(g.DdlPath)
		if err != nil {
			return nil, err
		}
		return ddlContent, nil
	
}

func NewContext(PSM string) context.Context {
	return context.Background()
}

// MysqlTypeToGo copy from https://github.com/gohouse/converter/blob/master/table2struct.go
var MysqlTypeToGo = map[string]string{
	"int":                "int32",
	"integer":            "int32",
	"tinyint":            "int16",
	"smallint":           "int16",
	"mediumint":          "int32",
	"bigint":             "int64",
	"int unsigned":       "int64",
	"integer unsigned":   "int64",
	"tinyint unsigned":   "int64",
	"smallint unsigned":  "int64",
	"mediumint unsigned": "int64",
	"bigint unsigned":    "int64",
	"bit":                "int64",
	"bool":               "bool",
	"enum":               "string",
	"set":                "string",
	"varchar":            "string",
	"char":               "string",
	"tinytext":           "string",
	"mediumtext":         "string",
	"text":               "string",
	"longtext":           "string",
	"blob":               "string",
	"tinyblob":           "string",
	"mediumblob":         "string",
	"longblob":           "string",
	"date":               "time.Time", // time.Time or string
	"datetime":           "time.Time", // time.Time or string
	"timestamp":          "time.Time", // time.Time or string
	"time":               "time.Time", // time.Time or string
	"float":              "float32",
	"double":             "float64",
	"decimal":            "float64",
	"binary":             "string",
	"varbinary":          "string",
	"json":               "string",
}

var golangKeywords = collection.NewSet[string](
	"break", "default", "func", "interface", "select",
	"case", "defer", "go", "map", "struct",
	"chan", "else", "goto", "package", "switch",
	"const", "fallthrough", "if", "range", "type",
	"continue", "for", "import", "return", "var",
)
