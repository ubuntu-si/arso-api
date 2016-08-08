package main

import (
	"os"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	m := gin.Default()

	m.Use(static.Serve("/", static.LocalFile("static", true)))
	m.GET(`/potresi.json`, func(c *gin.Context) {
		c.JSON(200, ARSOPotresi())
		return
	})

	m.GET(`/postaje.json`, func(c *gin.Context) {
		c.JSON(200, ARSOVreme())
		return
	})

	m.GET(`/vreme/:postaja`, func(c *gin.Context) {
		name := c.Param("postaja")
		for _, p := range ARSOVreme() {
			if name == p.ID {
				c.JSON(200, p)
				return
			}
		}
		c.JSON(404, gin.H{"Status": "Not found: " + name})
		return
	})

	m.GET(`/potresi.xml`, func(c *gin.Context) {
		c.XML(200, ARSOPotresi())
		return
	})

	m.GET(`/postaje.xml`, func(c *gin.Context) {
		c.XML(200, ARSOVreme())
		return
	})

	m.Run(":" + port)
}
