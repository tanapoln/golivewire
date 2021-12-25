package main

import (
	"fmt"
	"html/template"
	"path"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
	"github.com/tanapoln/golivewire"

	_ "github.com/tanapoln/golivewire/example/kitchensink/component"
)

func currentDir() string {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		panic("error get caller")
	}
	return path.Dir(filename)
}

func main() {
	golivewire.EnableMethodCamelCaseSupport = true
	golivewire.CORSOptions = &cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
	}

	srv := gin.Default()

	srv.Static("/static", path.Join(currentDir(), "public", "static"))
	srv.SetFuncMap(template.FuncMap{
		"livewire": golivewire.LivewireTemplateFunc,
	})
	srv.LoadHTMLGlob(path.Join(currentDir(), "templates", "**"))

	srv.GET("/home", func(c *gin.Context) {
		c.HTML(200, "home.tmpl", nil)
	})

	srv.GET("/livewire-dusk/:class", func(c *gin.Context) {
		className := c.Param("class")
		componentName := laravelClassToComponentName(className)

		fmt.Printf("Rendering dusk component: %s\n", componentName)

		c.HTML(200, "livewire-dusk-component.tmpl", componentName)
	})

	srv.Use(gin.WrapH(golivewire.NewAjaxHandler()))

	err := srv.Run("127.0.0.1:8081")
	if err != nil {
		panic(err)
	}
}

func laravelClassToComponentName(className string) string {
	dot := strings.ReplaceAll(className, "\\", ".")
	return strings.ToLower(dot)
}
