package proxy

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const (
	//AccountsURL  = "http://localhost:8654"
	AccountsURL = "http://account:8654"
	//SavedURL     = "http://localhost:8666"
	SavedURL = "http://saced:8666"
	//StreamingURL = "http://localhost:7666"
	StreamingURL = "http://streaming:7666"
)

func Run() {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "authorization"}

	r.Use(cors.New(config))

	upload := r.Group("/_upload")
	{
		upload.POST("/avatar", UploadAvatarHandler)
		upload.POST("/files", UploadFileHandler)
		upload.POST("/music/file", UploadMusicFileHandler)
	}

	accounts := r.Group("/_accounts")
	{

		accounts.Any("/*path", AccountsHandler)
	}

	saved := r.Group("/_saved")
	{
		saved.Any("/*path", SavedHandler)
	}
	streaming := r.Group("/_streaming")
	{
		streaming.Any("/*path", StreamingHandler)
	}

	if err := r.Run(":8032"); err != nil {
		return
	}
}

func AccountsHandler(c *gin.Context) {
	r, err := url.Parse(AccountsURL)
	if err != nil {
		panic(err)
	}

	p := httputil.NewSingleHostReverseProxy(r)
	path := "/_accounts" + c.Param("path")

	p.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Host = r.Host
		req.URL.Scheme = r.Scheme
		req.URL.Host = r.Host
		req.URL.Path = path
	}

	p.ServeHTTP(c.Writer, c.Request)
}

func SavedHandler(c *gin.Context) {
	r, err := url.Parse(SavedURL)
	if err != nil {
		panic(err)
	}

	p := httputil.NewSingleHostReverseProxy(r)
	path := "/_saved" + c.Param("path")

	p.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Host = r.Host
		req.URL.Scheme = r.Scheme
		req.URL.Host = r.Host
		req.URL.Path = path
	}

	p.ServeHTTP(c.Writer, c.Request)
}

func StreamingHandler(c *gin.Context) {
	r, err := url.Parse(StreamingURL)
	if err != nil {
		panic(err)
	}

	p := httputil.NewSingleHostReverseProxy(r)
	path := "/_streaming" + c.Param("path")

	p.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Host = r.Host
		req.URL.Scheme = r.Scheme
		req.URL.Host = r.Host
		req.URL.Path = path
	}

	p.ServeHTTP(c.Writer, c.Request)
}
