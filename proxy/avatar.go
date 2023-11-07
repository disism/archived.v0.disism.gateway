package proxy

import (
	"fmt"
	"github.com/gin-gonic/gin"
	shell "github.com/ipfs/go-ipfs-api"
	"log"
	"os"
)

func UploadAvatarHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(1)
	}
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

	c.JSON(200, &UploadFileObject{
		Name: file.Filename,
		CID:  cid,
		Size: uint64(file.Size),
	})
}
