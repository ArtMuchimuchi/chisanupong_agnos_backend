package main

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"net/http"
	"unicode"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	host     = "postgres"
	port     = 5432
	database = "mydatabase"
	username = "myuser"
	password = "mypassword"
)

type RequestLog struct {
	Password string `json:"init_password"`
}

type ResponseLog struct {
	NumSteps int `json:"num_of_steps"`
}

var db *sql.DB

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, username, password, database)

	sdb, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatal(err)
	}

	db = sdb

	defer db.Close()

	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	router.POST("/api/strong_password_steps", checkPassword)

	err = router.Run("0.0.0.0:5050")

	if err != nil {
		log.Fatal(err)
	}
}

func checkPassword(c *gin.Context) {
	var newReq RequestLog

	err := c.BindJSON(&newReq)

	if err != nil {
		log.Fatal(err)
	}

	err = createRequestLog(&newReq)

	if err != nil {
		log.Fatal(err)
	}

	//check
	var newRes ResponseLog
	newRes.NumSteps = getStrongSteps(newReq.Password)

	err = createResponseLog(&newRes)

	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, newRes)
}

func getStrongSteps(password string) int {
	const upper int = 0
	const lower int = 1
	const digit int = 2
	const minPassLength int = 6
	const maxPassLength int = 19
	var flagList [3]bool = [3]bool{false, false, false}
	var tempString rune = 0
	var repeatingCount int = 0
	var repeatingFix int = 0
	var slotFix int = 0
	var replaceFix int = 0
	fmt.Println(len(password))
	for _, c := range password {
		if slotFix <= maxPassLength-1 {
			if tempString == 0 {
				tempString = c
			} else if c == tempString {
				repeatingCount++
			} else {
				repeatingCount = 0
				tempString = c
			}
			if repeatingCount == 2 {
				repeatingCount = 0
				tempString = 0
				repeatingFix++
			}
		}
		if unicode.IsUpper(c) {
			flagList[upper] = true
		} else if unicode.IsLower(c) {
			flagList[lower] = true
		} else if unicode.IsDigit(c) {
			flagList[digit] = true
		}

		slotFix++
	}
	if slotFix < minPassLength {
		slotFix = minPassLength - slotFix
	} else if slotFix > maxPassLength {
		slotFix = slotFix - maxPassLength
	} else {
		slotFix = 0
	}
	for _, i := range flagList {
		if i == false {
			replaceFix++
		}
	}
	fmt.Printf("slot %v repeat %v replace %v", slotFix, repeatingFix, replaceFix)
	if len(password) <= maxPassLength {
		return int(math.Max(math.Max(float64(repeatingFix), float64(slotFix)), float64(replaceFix)))
	} else {
		return int(math.Max(float64(repeatingFix), float64(replaceFix))) + slotFix
	}
}

func createRequestLog(req *RequestLog) error {
	_, err := db.Exec("INSERT INTO request(init_password) VALUES ($1);",
		req.Password)
	return err
}

func createResponseLog(res *ResponseLog) error {
	_, err := db.Exec("INSERT INTO response(num_of_steps) VALUES ($1);",
		res.NumSteps)
	return err
}
