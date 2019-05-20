package main

import (
	"os"
	"fmt"
	"regexp"
	"strings"
	"net/http"
	"crypto/rand"
	"database/sql"

	"github.com/caarlos0/env/v5"
	"github.com/labstack/echo"
	_ "github.com/mattn/go-sqlite3"
)

type (
	Config struct {
		Name string `env:"HTTP_NAME" envDefault:"http://0.0.0.0"`
		IP string `env:"HTTP_IP" envDefault:""`
		Port string `env:"HTTP_PORT" envDefault:"80"`
		UI string `env:"UI_PATH" envDefault:"admin"`
		DbDir string `env:"DB_DIRECTORY" envDefault:"data"`
		Default string `env:"DEFAULT_URL" envDefault:"https://google.com"`
	}

	Redirect struct {
		Code string `json:"code"`
		URL string `json:"url"`
	}
)

var (
	cfg Config
	db *sql.DB
	redirects = make(map[string]string)
	urlRegex = regexp.MustCompile(`(?:(?:https?)://)(?:([-A-Z0-9+&@#/%=~_|$?!:,.]*)|[-A-Z0-9+&@#/%=~_|$?!:,.])*(?:([-A-Z0-9+&@#/%=~_|$?!:,.]*)|[A-Z0-9+&@#/%=~_|$])`) //Inb4 someone complains this is a dumb and inefficent regexp and/or someone also complains this isn't the right way to deal with global vars in Go
)

func main() {
	cfg = Config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Println("Error reading configuration file")
		fmt.Println(err)
		badStart("Error reading configuration file. See the console for details.")
	}

	if !urlRegex.MatchString(cfg.Name) {
		fmt.Println("The HTTP_NAME enviorment variable is not set correctly. You entered: \"" + cfg.Name + "\". Make sure it is in the form: http/https://your.domain.com") //Probably make a string formatter thingy for this
		badStart("The HTTP_NAME enviorment variable is not set correctly. You entered: \"" + cfg.Name + "\". Make sure it is in the form: http/https://your.domain.com")
	}

	if !urlRegex.MatchString(cfg.Default) {
		fmt.Println("The DEFAULT_URL enviorment variable is not set correctly. You entered: \"" + cfg.Name + "\". Make sure it is in the form: http/https://your.domain.com")
		badStart("The DEFAULT_URL enviorment variable is not set correctly. You entered: \"" + cfg.Name + "\". Make sure it is in the form: http/https://your.domain.com")
	}
	
	ok := false
	db, ok = dbConnect()
	if !ok {
		badStart("There was a problem connecting to the database. See server console for details.")
	}
	redirects = buildRedirects()

	//Start the web services
	e := echo.New()
	e.Debug = false
	e.HidePort = true
	e.HideBanner = true
	
	e.GET("/*", redirect)
	e.Static("/" + cfg.UI, "web")
	//e.POST("/api/login", apiLogin)
	e.GET("/api/list", apiList)
	//e.GET("/api/stats", apiStats)
	e.POST("/api/create", apiCreate)
	e.DELETE("/api/remove", apiRemove)
	e.PATCH("/api/update", apiUpdate)

	bind := cfg.IP + ":" + cfg.Port
	fmt.Println("Starting server on: " + bind)
	e.Start(bind)
}

//Start up a tiny go http server to inform the user of the problem
func badStart(message string) {
	fmt.Println("Bad Start!")

	http.HandleFunc("/", func (w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(message))
	})

	if cfg.Port != "" {
		http.ListenAndServe(cfg.IP + ":" + cfg.Port, nil)
	} else {
		http.ListenAndServe(":80", nil)
	}
}

func redirect(c echo.Context) error {
	code := strings.TrimPrefix(c.Request().URL.Path, "/")

	if code == "" {
		return c.Redirect(http.StatusSeeOther, cfg.Default)
	}

	url, ok := redirects[code]
	if !ok {
		return c.String(http.StatusBadRequest, "No redirect matching the code: \"" + code + "\" found!")
	}
	return c.Redirect(http.StatusSeeOther, url)
}

func apiCreate(c echo.Context) error {
	r := Redirect{}

	if err := c.Bind(&r); err != nil {
		return c.String(http.StatusInternalServerError, err.Error()) //Needs formatted
	}

	if r.URL == "" {
		return c.JSON(http.StatusUnprocessableEntity, "You must enter a URL") //Needs formatted
	}

	if !urlRegex.MatchString(r.URL) {
		return c.JSON(http.StatusUnprocessableEntity, r.URL + " is not a valid URL. Make sure it is in the form: 'http/https://your.domain.com' before entry.'") //Needs formatted
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
		return c.JSON(http.StatusUnprocessableEntity, "A redirect with this code already exists") //Needs formatted
	}

	fmt.Println(redirects)

	_, err := db.Exec("INSERT INTO redirects (code, url, created) VALUES ('"+code+"', '"+r.URL+"', date('now'))")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error()) //Needs formatted
	}

	redirects[code] = r.URL
	fmt.Println(redirects)

	return c.JSON(http.StatusOK, cfg.Name + "/" + code) //Needs formatted
}

func apiList(c echo.Context) error {
	return c.JSON(http.StatusOK, redirects) //Needs formatted (Maybe?)
}

func apiRemove(c echo.Context) error {
	r := Redirect{}

	if err := c.Bind(&r); err != nil {
		return c.String(http.StatusInternalServerError, err.Error()) //Needs formatted
	}

	if r.Code == "" {
		return c.JSON(http.StatusUnprocessableEntity, "You must enter a code") //Needs formatted
	}

	
	_, err := db.Exec("DELETE FROM redirects WHERE code = '" + r.Code + "';")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error()) //Needs formatted
	}

	delete(redirects, r.Code)

	return c.JSON(http.StatusOK, "Deleted.") //Needs formatted
}

func apiUpdate(c echo.Context) error {
	r := Redirect{}

	if err := c.Bind(&r); err != nil {
		return c.String(http.StatusInternalServerError, err.Error()) //Needs formatted
	}

	if r.Code == "" {
		return c.JSON(http.StatusUnprocessableEntity, "You must specify a code")
	}

	if !urlRegex.MatchString(r.URL) {
		return c.JSON(http.StatusUnprocessableEntity, r.URL + " is not a valid URL. Make sure it is in the form: 'http/https://your.domain.com' before entry.'") //Needs formatted
	}

	_, err := db.Exec("UPDATE redirects SET url = '" + r.URL + "' where code = '" + r.Code + "';")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error()) //Needs formatted
	}

	redirects[r.Code] = r.URL
	
	return c.JSON(http.StatusOK, "Updated.")//Needs formatted
}

func buildRedirects() map[string]string {
	//Populate list of redirects with that which is stored in the DB
	results, err := db.Query("SELECT code, url FROM redirects")
	if err != nil { 
		fmt.Println(err)
	}
 
	for results.Next() {
		var (
			code string
			url string
		)
		if err := results.Scan(&code, &url); err != nil {
			fmt.Println(err)
		}

		redirects[code] = url
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

		db, err := sql.Open("sqlite3", dbPath + "/config.db")
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
	db, err := sql.Open("sqlite3", dbPath + "/config.db")
	if err != nil {
		fmt.Println(err)
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
