package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func StoreDirectory(c *gin.Context) {

}

func ListDirectory(c *gin.Context) {
	server := GetContextServer(c)
	if server == nil {
		c.Abort()
	}

	fmt.Println(server)
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, "test")
}

func CopyFile(c *gin.Context) {

}

func MoveFile(c *gin.Context) {

}

func CompressFile(c *gin.Context) {

}

func DecompressFile(c *gin.Context) {

}

func StatFile(c *gin.Context) {

}

func ReadFileContents(c *gin.Context) {

}

func WriteFileContents(c *gin.Context) {

}

func DeleteFile(c *gin.Context) {

}

func DownloadFile(c *gin.Context) {

}
