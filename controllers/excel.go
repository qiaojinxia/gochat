package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
	"goChat/conf"
	"goChat/models"
	"goChat/utils"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"strings"
)

func Tutorial(c *gin.Context) {
	c.HTML(http.StatusOK, "tutorial.html", gin.H{
		"result": c.Param("content"),
	})
}

func SaveExcel(c *gin.Context) {
	fileName := c.Query("file_name")
	// 获取文件
	file, err := c.FormFile("uploadFile")
	if err != nil {
		// 处理错误
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if fileName == "" {
		// 生成一个新的KSUID
		id, err := ksuid.NewRandom()
		if err != nil {
			// 处理错误
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// 生成随机文件名
		fileName = id.String() + path.Ext(file.Filename)
	}
	filePath := filepath.Join("data", fileName)

	// 保存文件
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		// 处理错误
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回文件名
	c.JSON(http.StatusOK, gin.H{"filename": fileName})
}

func GetExcel(c *gin.Context) {
	excelName := c.Param("name")
	// 读取本地xlsx文件
	filePath := filepath.Join("data", excelName)
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 将xlsx文件转为二进制数据
	excelData, err := f.WriteToBuffer()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 设置响应头信息
	c.Header("Content-Disposition", "attachment; filename=test.xlsx")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelData.Bytes())
}

func Excel(c *gin.Context) {
	c.HTML(http.StatusOK, "excel.html", gin.H{})
}

func ChatExcelAbout(c *gin.Context) {
	c.HTML(http.StatusOK, "about.html", gin.H{})
}

func ChatExcelGenerate(c *gin.Context) {

	postData := &models.PostData{}
	if err := c.BindJSON(&postData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	correctionCt := 0 //纠正计数
	prompt := models.CreatePrompt(postData)
R:
	resultSql, err := utils.RequestRemoteGptApi(prompt)

	log.Println("接收提示词:", postData.TextCommand)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resultSql = utils.MatchChatGptPythonCode(resultSql)
	if resultSql == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error text"})
		return
	}
	resultSql = utils.MatchHandlerSqlTable(utils.MatchHandlerSqlTableS(resultSql))
	fileName, err := utils.GenerateID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	optTables := make(map[string]string)
	index := 0
	for urlPath, _ := range postData.OpTables {
		urlPath = utils.FilterTrim(urlPath)
		optTables[fmt.Sprintf("table%d", index)] += fmt.Sprintf("data/%s", urlPath)
		index++
	}
	args1, err := json.Marshal(optTables)
	if err != nil {
		log.Fatal(err)
	}
	if resultSql == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "错误的命令!"})
		return
	}
	args2 := resultSql
	args3 := fileName
	// 传入参数 格式 python script\sql_exc_template.py "table1 data/filename1.xlsx table2 data/filename2.xlsx" "SELECT *,数量 * 单价 AS 总金额 FROM table1|table2" saveFilePath
	// 运行python脚本来处理表格
	cmdOut, err := utils.RunPythonScript(filepath.Join("script", "sql_exc_template.py"),
		string(args1), args2, args3)
	if err != nil {
		log.Printf("failed to run python script: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if strings.Contains(cmdOut, "Error") {
		prompt = models.JustifyPrompt(cmdOut)
		correctionCt += 1
		if strings.Contains(cmdOut, "syntax error") && correctionCt < conf.GetConfigInt("setting", "sql_correction_times") {
			log.Printf("failed to run python scriptOut: %v\n", cmdOut)
			goto R
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": cmdOut})
		return
	}
	exists := utils.FileExists(path.Join("data", utils.FilterTrim(cmdOut)))
	if !exists {
		c.JSON(http.StatusInternalServerError, "excel not create")
		return
	}

	// 返回文件名
	c.JSON(http.StatusOK, gin.H{"filename": cmdOut})
}
