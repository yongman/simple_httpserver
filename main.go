package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

var http_port int
var root_path string

func init() {
	flag.IntVar(&http_port, "p,port", 8008, "http server port")
	flag.StringVar(&root_path, "r,root", "./data", "root path")
}

var html_before_tmpl string = `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8" />
    <title>Simple HTTP</title>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/semantic-ui/1.11.8/semantic.min.css"/>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/jquery/2.1.3/jquery.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/semantic-ui/1.11.8/semantic.min.js"></script>
</head>
<body>

<form class="ui form" action="/upload/" enctype="multipart/form-data" method="post">
  <div class="field">
    <label>File</label>
    <input type="file" name="filename">
  </div>
  <button class="ui button" type="submit">upload</button>
</form>
<div id="filelist">
  <table class="ui celled striped table">
  <thead>
    <tr><th colspan="3">
      Files
    </th>
  </tr></thead>
  <tbody>
`
var html_after_tmpl string = `
  </tbody>
</table>
</div>

</body>
</html>
`

func main() {
	flag.Parse()
	rt := gin.Default()
	rt.GET("/ui", func(c *gin.Context) {
		//gererate file list
		inner := func(filename, filesize, lasttime string) string {
			tr := "<tr><td class='collapsing'><a href='/file/" + filename + "'><i class='file icon'></i>" + filename + "</a></td>"
			tr = tr + "<td>" + filesize + "MB</td>"
			tr = tr + "<td class='right aligned collapsing'>" + lasttime + "</td>"
			tr = tr + "</td>"
			return tr
		}
		fileinfos, err := ioutil.ReadDir(root_path)
		if err != nil {
			fmt.Println(err)
			return
		}
		trs := ""
		for _, fi := range fileinfos {
			trs = trs + inner(fi.Name(), fmt.Sprintf("%v", fi.Size()/1024/1024), fi.ModTime().String())
		}
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, html_before_tmpl+trs+html_after_tmpl)
	})
	rt.StaticFS("/file/", http.Dir(root_path))
	rt.POST("/upload/", func(c *gin.Context) {
		file, header, err := c.Request.FormFile("filename")
		if err != nil {
			fmt.Println(err)
			return
		}
		filename := header.Filename
		out, err := os.Create(root_path + "/" + filename)
		if err != nil {
			fmt.Println(err)
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			fmt.Println(err)
		}
		c.Redirect(http.StatusMovedPermanently, "/ui")
	})
	rt.Run(":" + strconv.Itoa(http_port))
}
