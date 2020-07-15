package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type herokuApp struct {
	db *sql.DB
}

func newApp(db *sql.DB) *herokuApp {
	return &herokuApp{db}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln("Could not open database", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalln("Could not ping to database", err)
	}
	app := newApp(db)
	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.html")
	router.Static("/static", "static")
	router.Static("/gojs", "gojs")

	router.GET("/", app.indexHandler)
	router.GET("/tick", app.tickHandler)

	router.Run(":" + port)

}

func (a *herokuApp) indexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "home.html", nil)
}

func (a *herokuApp) tickHandler(c *gin.Context) {
	if _, err := a.db.Exec(
		"CREATE TABLE IF NOT EXISTS ticks (tick timestamp);",
	); err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error creating database table: %q", err))
		return
	}

	if _, err := a.db.Exec(
		"INSERT INTO ticks VALUES (now());",
	); err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error incrementing tick: %q", err))
		return
	}

	var tick time.Time
	if err := a.db.QueryRow("SELECT tick FROM ticks ORDER BY tick DESC;").
		Scan(&tick); err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error reading ticks: %q", err))
		return
	}
	c.String(http.StatusOK,
		fmt.Sprintf("Read from database: %s\n", tick.String()))
}
