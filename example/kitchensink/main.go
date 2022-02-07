package main

import (
	"context"
	"fmt"
	"html/template"
	"path"
	"runtime"
	"strings"
	"time"

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
	golivewire.SetBaseURL("http://localhost:8081")
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

	srv.Use(func(c *gin.Context) {
		c.Request = golivewire.WithRequestContext(c.Request)
		c.Next()
	})

	srv.GET("/home", func(c *gin.Context) {
		done := make(chan struct{})
		ctx, cancelFunc := context.WithTimeout(c.Request.Context(), time.Millisecond*1000)
		defer cancelFunc()

		go func() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("Recover: %s\n", err)
				}
			}()

			c.HTML(200, "home.tmpl", ctx)
			done <- struct{}{}
		}()

		select {
		case <-done:
			return
		case <-ctx.Done():
			c.String(503, "Server timeout exceeded")
		}
	})

	srv.GET("/livewire-dusk/:class", func(c *gin.Context) {
		className := c.Param("class")
		componentName := laravelClassToComponentName(className)

		fmt.Printf("Rendering dusk component: %s\n", componentName)

		c.HTML(200, "livewire-dusk-component.tmpl", gin.H{
			"ctx":  c.Request.Context(),
			"name": componentName,
		})
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
