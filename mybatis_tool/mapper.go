package mybatis_tool

import (
	"encoding/xml"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"code_tools/util"

	"github.com/xwb1989/sqlparser"
)

type (
	// 解析 Mapper 的 SQL 语法树得到涉及的表
	Mapper struct {
		// 固定名字 XMLName
		XMLName   xml.Name
		Namespace string `xml:"namespace,attr"`
		Sql       []sql  `xml:"sql"`
		Select    []sql  `xml:"select"`
		Insert    []sql  `xml:"insert"`
		Update    []sql  `xml:"update"`
		Delete    []sql  `xml:"delete"`
	}

	sql struct {
		Comment string `xml:",comment"`
		Id      string `xml:"id,attr"`
		SQL     string `xml:",innerxml"`
	}
)

func (m Mapper) FeaturesName() string {
	return "Mapper 清单"
}

func (m Mapper) Before(i ...interface{}) (results []interface{}) {
	return
}

func (m Mapper) FileNames() []string {
	return []string{"mybatis_list", "mybatis_comma_list"}
}

func (m Mapper) Titles() [][]string {
	return [][]string{
		{"id", "comment", "type", "tables"},
		{"id", "comment", "type", "tables"},
	}
}

func (m Mapper) Lines(rootPath string, path string) [][][]string {
	if filepath.Ext(path) != ".xml" {
		return nil
	}

	model := Mapper{}
	err := util.ParseXml(path, &model)
	if err != nil {
		return nil
	}
	if model.XMLName.Local != "mapper" {
		return nil
	}
	if model.Namespace != "" {
		model.Namespace += "."
	}

	sqlMap := make(map[string]string)
	for _, s := range model.Sql {
		sqlMap[s.Id] = xmlSQL(s.SQL)
	}

	result := [][][]string{{}, {}}
	appendData(result, sqlMap, path, "select", model.Namespace, model.Select)
	appendData(result, sqlMap, path, "insert", model.Namespace, model.Insert)
	appendData(result, sqlMap, path, "update", model.Namespace, model.Update)
	appendData(result, sqlMap, path, "delete", model.Namespace, model.Delete)
	return result
}

var paramRegExp = regexp.MustCompile(`[#$]{(?:\w+\.)?(\w+)(?:,[^}]+)?}`)
var xmlRegExp = regexp.MustCompile(`<!--[\s\S]*-->|</?if[^>]*>\s*|</where>\s*|</set>\s*|,\s*</trim>|,\s*</set>|<!\[CDATA\[|]]>`)
var whereRegExp = regexp.MustCompile(`<where>\s*(?i:AND)?`)
var setRegExp = regexp.MustCompile(`<set>\s*`)
var trimRegExp = regexp.MustCompile(`<trim prefix="([^"]*)"[^>]*>\s*`)
var trim2RegExp = regexp.MustCompile(`,\s*</trim>`)
var spaceLineRegExp = regexp.MustCompile(`\n\s*\n`)

func xmlSQL(sql string) string {
	sql = paramRegExp.ReplaceAllString(sql, `$1`)
	sql = xmlRegExp.ReplaceAllString(sql, ``)
	sql = whereRegExp.ReplaceAllString(sql, `WHERE `)
	sql = setRegExp.ReplaceAllString(sql, `SET `)
	sql = trimRegExp.ReplaceAllString(sql, `$1`)
	sql = trim2RegExp.ReplaceAllString(sql, `)`)
	sql = spaceLineRegExp.ReplaceAllString(sql, "\n")
	return sql
}

var includeRegExp = regexp.MustCompile(`<include refid="(\w+)" */?>(?:</include>)?`)
var positionRegExp = regexp.MustCompile(`position (\d+)`)

func appendData(result [][][]string, sqlMap map[string]string, path, sqlType, namespace string, selectSQL []sql) {
	for _, s := range selectSQL {
		sql := xmlSQL(s.SQL)
		includeId := includeRegExp.FindStringSubmatch(sql)
		if includeId != nil {
			sql = includeRegExp.ReplaceAllString(sql, sqlMap[includeId[1]])
		}
		stmt, err := sqlparser.Parse(sql)
		if err != nil {
			msg := positionRegExp.FindStringSubmatch(err.Error())
			if msg != nil {
				index, err := strconv.Atoi(msg[1])
				if err == nil {
					sql = util.Position(sql, index)
				}
			}
			util.Logger.Printf("sqlparser.Parse(...): %v\n"+
				"file:///%v\n"+
				"%v%v\n"+
				"%v\n",
				err,
				filepath.ToSlash(path),
				namespace, s.Id,
				sql)
			continue
		}
		// 转为切片
		var tables []string
		switch stmt := stmt.(type) {
		case *sqlparser.Select:
			tables = appendTableExprs(tables, stmt.From)
		case *sqlparser.Update:
			tables = appendTableExprs(tables, stmt.TableExprs)
		case *sqlparser.Delete:
			tables = appendTableExprs(tables, stmt.TableExprs)
		case *sqlparser.Insert:
			tables = append(tables, stmt.Table.Name.String())
		}
		id := namespace + s.Id
		for _, table := range tables {
			result[0] = append(result[0], []string{
				id,
				s.Comment,
				sqlType,
				table,
			})
		}
		result[1] = append(result[1], []string{
			id,
			s.Comment,
			sqlType,
			strings.Join(tables, ", "),
		})
	}
}

func appendTableExprs(tables []string, tableExprs sqlparser.TableExprs) []string {
	for _, tableExpr := range tableExprs {
		tables = appendTableExpr(tables, tableExpr)
	}
	return tables
}

func appendTableExpr(tables []string, tableExpr sqlparser.TableExpr) []string {
	switch tableExpr := tableExpr.(type) {
	case *sqlparser.AliasedTableExpr:
		switch t := tableExpr.Expr.(type) {
		case sqlparser.TableName:
			tables = append(tables, t.Name.String())
		case *sqlparser.Subquery:
			// recursion through appendTableExprs()
			tables = appendTableExprs(tables, t.Select.(*sqlparser.Select).From)
		}
	case *sqlparser.JoinTableExpr:
		tables = appendTableExpr(tables, tableExpr.LeftExpr)
		tables = appendTableExpr(tables, tableExpr.RightExpr)
	case *sqlparser.ParenTableExpr:
		// recursion through appendTableExprs()
		tables = appendTableExprs(tables, tableExpr.Exprs)
	}
	return tables
}
