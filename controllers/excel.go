package controllers

import (
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
	"net/http"
	"path"
	"path/filepath"
)

func SaveExcel(c *gin.Context) {
	// 获取文件
	file, err := c.FormFile("uploadFile")
	if err != nil {
		// 处理错误
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 生成一个新的KSUID
	id, err := ksuid.NewRandom()
	if err != nil {
		// 处理错误
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 生成随机文件名
	fileName := id.String() + path.Ext(file.Filename)
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
