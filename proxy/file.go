package proxy

import (
	"fmt"
	"github.com/gin-gonic/gin"
	shell "github.com/ipfs/go-ipfs-api"
	"log"
	"net/http"
	"os"
)

type UploadFileObject struct {
	Name string `json:"name,omitempty"`
	CID  string `json:"cid,omitempty"`
	Size uint64 `json:"size,omitempty"`
}

const (
	//ipfs       = "localhost:5001"
	ipfs       = "ipfs:5001"
	cidVersion = 1
)

func UploadFileHandler(c *gin.Context) {
	var fs []UploadFileObject

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	files := form.File["files"]
	for _, file := range files {
		sh := shell.NewShell(ipfs)

		b, err := file.Open()
		if err != nil {
			log.Fatal(err)
		}

		cid, err := sh.Add(b, shell.CidVersion(cidVersion))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s", err)
			os.Exit(1)
		}
		fs = append(fs, UploadFileObject{
			Name: file.Filename,
			CID:  cid,
			Size: uint64(file.Size),
		})
	}

	c.JSON(200, gin.H{
		"message": "OK",
		"files":   fs,
	})
}
