package models

import (
	"fmt"
	"goChat/utils"
	"strings"
)

type Clos []interface{}

type PostData struct {
	TextCommand string                   `json:"text_command"`
	OpTables    map[string][]interface{} `json:"op_tables"` //tablePath:tableHeaderCols Name
}

//CreatePrompt 根据提交的表格创建提示词
func CreatePrompt(data *PostData) string {
	tablesNames := ""
	index := 0
	for _, tableColsName := range data.OpTables {
		tablesNames += fmt.Sprintf("table%d 列名为 %v ", index, tableColsName)
		index += 1
	}
	data.TextCommand = strings.ReplaceAll(data.TextCommand, "表格", "table")
	var prompt = fmt.Sprintf("现用pandas加载了%d张表 表名分别为 %s"+
		" 需要根据需求生成pandasSql 来分析表数据 要求如下: %s ,直接返回python能执行的,pandasSql语句不需要其多余内容,不需要其他任何解释",
		len(data.OpTables), tablesNames, data.TextCommand)
	return prompt
}

//JustifyPrompt 纠正sql
func JustifyPrompt(errorMsg string) string {
	errorMsg += "\n请修改下报错代码, 能让他运行, 不需要任何解释, 直接输出可以执行的代码"
	return errorMsg
}

func (c Clos) String() string {
	outPut := ""
	for _, v := range c {
		outPut += utils.ToString(v) + " |"
	}
	return outPut
}
