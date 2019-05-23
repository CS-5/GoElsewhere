package main

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/caarlos0/env/v5"
	"github.com/labstack/echo"
	_ "github.com/mattn/go-sqlite3"
)

type (
	config struct {
		Name    string `env:"HTTP_NAME" envDefault:"http://0.0.0.0"`
		IP      string `env:"HTTP_IP" envDefault:""`
		Port    string `env:"HTTP_PORT" envDefault:"80"`
		UI      string `env:"UI_PATH" envDefault:"admin"`
		DbDir   string `env:"DB_DIRECTORY" envDefault:"data"`
		Default string `env:"DEFAULT_URL" envDefault:"https://google.com"`
	}

	RedirectEntry struct {
		ID      int64  `json:"id" form:"id" query:"id"`
		Code    string `json:"code" form:"code" query:"code"`
		URL     string `json:"url" form:"url" query:"url"`
		Link    string `json:"link" form:"link" query:"link"`
		Created string `json:"created" form:"created" query:"created"`
		Hits    int    `json:"hits" form:"hits" query:"hits"`
	}

	apiResponse struct {
		Good  bool          `json:"good"`
		Entry RedirectEntry `json:"entry"`
		Error string        `json:"error"`
	}
)

var (
	cfg       config
	db        *sql.DB
	redirects = make(map[string]RedirectEntry)
	urlRegex  = regexp.MustCompile(`(?:(?:https?)://)(?:([-A-Z0-9+&@#/%=~_|$?!:,.]*)|[-A-Z0-9+&@#/%=~_|$?!:,.])*(?:([-A-Z0-9+&@#/%=~_|$?!:,.]*)|[A-Z0-9+&@#/%=~_|$])`) //Inb4 someone complains this is a dumb and inefficent regex and/or someone also complains this isn't the right way to deal with global vars in Go
)

const (
	badEnvVar string = "The %s enviorment variable is not set correctly. You entered: '%s'. Make sure it is formatted correctly per the documentation here: (TODO: Wiki)"
)

func main() {
	cfg = config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Println("Error reading configuration file")
		fmt.Println(err)
		badStart("Error reading configuration file. See the console for details.")
	}

	if !urlRegex.MatchString(cfg.Name) {
		badStart(fmt.Sprintf(badEnvVar, "HTTP_NAME", cfg.Name))
	}

	if !urlRegex.MatchString(cfg.Default) {
		badStart(fmt.Sprintf(badEnvVar, "DEFAULT_URL", cfg.Default))
	}

	ok := false
	db, ok = dbConnect()
	if !ok {
		badStart("There was a problem connecting to the database. See server console for details.")
	}

	//Store the redirects found in the database in memory
	redirects = buildRedirects()

	//Start the web server
	e := echo.New()
	e.Debug = false
	e.HidePort = true
	e.HideBanner = true

	e.GET("/*", redirect)
	e.Static("/"+cfg.UI, "public")
	//e.POST("/api/login", apiLogin)
	e.GET("/api/list", apiList)
	e.POST("/api/create", apiCreate)
	e.PATCH("/api/update", apiUpdate)
	e.DELETE("/api/delete", apiRemove)

	bind := cfg.IP + ":" + cfg.Port
	fmt.Println("Starting server on: " + bind)
	e.Start(bind)
}

//Start up a tiny go http server to inform the user of the problem
func badStart(message string) {
	fmt.Println("Bad Start! -> '" + message + "'")

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(message))
	})

	if cfg.Port != "" {
		http.ListenAndServe(cfg.IP+":"+cfg.Port, nil)
	} else {
		http.ListenAndServe(":80", nil)
	}
}

func redirect(c echo.Context) error {
	code := strings.TrimPrefix(c.Request().URL.Path, "/")

	if code == "" {
		return c.Redirect(http.StatusSeeOther, cfg.Default)
	}

	redirect, ok := redirects[code]
	if !ok {
		return c.String(http.StatusBadRequest, "No redirect matching the code: '"+code+"' found!")
	}

	redirect.Hits++
	_, err := db.Exec("UPDATE redirects SET hits = '" + strconv.Itoa(redirect.Hits) + "' where code = '" + code + "';")
	if err != nil {
		fmt.Println("Unable to increase hits value in database. See the following error:")
		fmt.Println(err.Error())
	}

	return c.Redirect(http.StatusSeeOther, redirect.URL)
}

func apiCreate(c echo.Context) error {
	r := RedirectEntry{}

	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusInternalServerError, apiResponse{
			Good:  false,
			Entry: r,
			Error: err.Error(),
		})
	}

	if r.URL == "" {
		return c.JSON(http.StatusUnprocessableEntity, apiResponse{
			Good:  false,
			Entry: r,
			Error: "You did not submit a URL. A URL is required to perform a redirect.",
		})
	}

	if !urlRegex.MatchString(r.URL) {
		return c.JSON(http.StatusUnprocessableEntity, apiResponse{
			Good:  false,
			Entry: r,
			Error: "You submitted: '" + r.URL + "' which does not appear to be a valid URL in the form of: 'http/https://example.com/sub'",
		})
	}

	code := r.Code
	if code == "" {
		code = generateCode(6)
		finish := false
		for !finish {
			if _, exists := redirects[code]; exists {
				code = generateCode(6)
			} else {
				finish = true
			}
		}
	} else if _, exists := redirects[code]; exists {
		return c.JSON(http.StatusUnprocessableEntity, apiResponse{
			Good:  false,
			Entry: r,
			Error: "The code: '" + code + "' already exists. Please try again.",
		})
	}

	result, err := db.Exec("INSERT INTO redirects (code, url, created, hits) VALUES ('" + code + "', '" + r.URL + "', date('now'), '0')")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, apiResponse{
			Good:  false,
			Entry: r,
			Error: err.Error(),
		})
	}

	id, _ := result.LastInsertId()

	redirects[code] = RedirectEntry{
		ID:      id,
		Code:    code,
		URL:     r.URL,
		Link:    cfg.Name + "/" + code,
		Created: time.Now().Format("01/02/2006"),
		Hits:    0,
	}

	return c.JSON(http.StatusOK, apiResponse{
		Good:  true,
		Entry: redirects[code],
		Error: "",
	})
}

func apiList(c echo.Context) error {
	return c.JSON(http.StatusOK, redirects) //Needs formatted (Maybe?)
}

func apiRemove(c echo.Context) error {
	r := RedirectEntry{}

	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusInternalServerError, apiResponse{
			Good:  false,
			Entry: r,
			Error: err.Error(),
		})
	}

	if r.Code == "" {
		return c.JSON(http.StatusUnprocessableEntity, apiResponse{
			Good:  false,
			Entry: r,
			Error: "You did not submit a code. A code is required to perform a redirect.",
		})
	}

	_, err := db.Exec("DELETE FROM redirects WHERE code = '" + r.Code + "';")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, apiResponse{
			Good:  false,
			Entry: r,
			Error: err.Error(),
		})
	}

	delete(redirects, r.Code)

	return c.JSON(http.StatusOK, apiResponse{
		Good:  true,
		Entry: r,
		Error: "",
	})
}

func apiUpdate(c echo.Context) error {
	r := RedirectEntry{}

	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusInternalServerError, apiResponse{
			Good:  false,
			Entry: r,
			Error: err.Error(),
		})
	}

	if r.Code == "" {
		return c.JSON(http.StatusUnprocessableEntity, "You must specify a code")
	}

	if !urlRegex.MatchString(r.URL) {
		return c.JSON(http.StatusUnprocessableEntity, apiResponse{
			Good:  false,
			Entry: r,
			Error: "You submitted: '" + r.URL + "' which does not appear to be a valid URL in the form of: 'http/https://example.com/sub'",
		})
	}

	_, err := db.Exec("UPDATE redirects SET url = '" + r.URL + "' where code = '" + r.Code + "';")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, apiResponse{
			Good:  false,
			Entry: r,
			Error: err.Error(),
		})
	}

	redirects[r.Code] = RedirectEntry{
		URL: r.URL,
	}

	return c.JSON(http.StatusOK, apiResponse{
		Good:  true,
		Entry: redirects[r.Code],
		Error: "",
	})
}

func buildRedirects() map[string]RedirectEntry {
	//Populate list of redirects with that which is stored in the DB
	results, err := db.Query("SELECT id, code, url, created, hits FROM redirects")
	if err != nil {
		fmt.Println(err)
	}

	for results.Next() {
		var (
			id   int64
			code string
			url  string
			date time.Time
			hits string
		)
		if err := results.Scan(&id, &code, &url, &date, &hits); err != nil {
			fmt.Println(err)
		}

		h, _ := strconv.Atoi(hits)

		redirects[code] = RedirectEntry{
			ID:      id,
			URL:     url,
			Link:    cfg.Name + "/" + code,
			Created: date.Format("01/02/2006"),
			Hits:    h,
		}
	}

	return redirects
}

//TODO: This function could probably use some work.
func dbConnect() (*sql.DB, bool) {
	//If the database file does not exist, create it and initialize the table
	dbPath := "/" + cfg.DbDir
	if _, err := os.Stat(dbPath + "/config.db"); os.IsNotExist(err) {
		fmt.Println("Creating configuration db file.")
		os.MkdirAll(dbPath, 0775)
		os.Create(dbPath + "/config.db")

		db, err := sql.Open("sqlite3", dbPath+"/config.db")
		if err != nil {
			fmt.Println("Error opening the database file. Might be a permissions issue, or might be something else.")
			fmt.Println(err)
			return nil, false
		}

		_, err = db.Exec("CREATE TABLE `redirects` ( `id` INTEGER PRIMARY KEY AUTOINCREMENT, `code` TEXT NOT NULL UNIQUE, `url` TEXT NOT NULL, `created` datetime NOT NULL, `hits` INTEGER )")
		if err != nil {
			fmt.Println("Error creating the 'Redirects' table.")
			fmt.Println(err)
			return nil, false
		}
		db.Close()
	} else if err != nil {
		fmt.Println("Other Error")
		fmt.Println(err)
		return nil, false
	}

	fmt.Println("Opening DB Connection")
	//Open DB Connection
	var err error
	db, err := sql.Open("sqlite3", dbPath+"/config.db")
	if err != nil {
		fmt.Println("There was a problem opening the database connection")
		fmt.Println(err)

		return db, false
	}

	return db, true
}

func generateCode(length int) string {
	var chars = "0123456789abcdefghijklmnopqrstuvwxyz"

	var bytes = make([]byte, length)
	var op = byte(len(chars))

	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = chars[b%op]
	}
	return string(bytes)
}
