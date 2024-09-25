package routes

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"

	api "apidemo/handlers/api"
	auth "apidemo/handlers/auth"
	middleware "apidemo/middleware"
)

// 自定义模板解析器
func parseTemplate(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	// 正则表达式匹配 <include file="xxx.html"></include> 或 <include file="xxx.html" />
	re := regexp.MustCompile(`<include\s+file="([a-zA-Z0-9./_-]+)"\s*/?>`)
	matches := re.FindAllStringSubmatch(string(content), -1)

	for _, match := range matches {
		includeFile := match[1]
		includeContent, err := ioutil.ReadFile(filepath.Join("html", includeFile))
		if err != nil {
			return "", err
		}
		// 替换 <include file="xxx.html"></include> 或 <include file="xxx.html" /> 为实际内容
		content = []byte(strings.ReplaceAll(string(content), match[0], string(includeContent)))
	}

	return string(content), nil
}

func SetupRouter(r *gin.Engine) {
	r_auth := r.Group("/auth")
	{
		r_auth.POST("/login", auth.Login)
	}
	r_api := r.Group("/api")
	{
		r_api.Use(middleware.JwtAuthMiddleware())
		r_api.GET("/users", api.GetUsers)
		r_api.GET("/status", api.GetStatus)
	}

	// Add 404 Not Found handler
	r.NoRoute(func(c *gin.Context) {
		// Get the request path and convert it to a file path
		path := c.Request.URL.Path
		if path == "/" {
			path = "/index"
		}
		filePath := filepath.Join("html", path+".html")

		// Check if the file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
			return
		}

		// Parse the template file
		parsedContent, err := parseTemplate(filePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		// Create and render the template
		tmpl, err := template.New("webpage").Funcs(template.FuncMap{
			"include": func(name string) (template.HTML, error) {
				content, err := ioutil.ReadFile(filepath.Join("html", name))
				if err != nil {
					return "", err
				}
				return template.HTML(content), nil
			},
		}).Parse(parsedContent)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		// Execute the template and write the result to the response
		c.Writer.Header().Set("Content-Type", "text/html")
		c.Writer.WriteHeader(http.StatusOK) // 确保状态码为 200
		if err := tmpl.Execute(c.Writer, nil); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
	})
}
