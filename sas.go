package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var router *gin.Engine

func main() {
	//-------------------------------------------------------

	//-------------------------------------------------------
	// Set the router as the default one shipped with Gin
	router := gin.Default()
	router.GET("/heartBeat", heartBeatHandler)

	// Start and run the server
	router.Run(":3000")
}

func heartBeatHandler(c *gin.Context) {
	db, err := sql.Open("mysql", "root@unix(/tmp/mysql.sock)/proftp")
	checkErr(err)

	rows, err := db.Query("SELECT userid, passwd, uid, gid FROM users")
	checkErr(err)

	for rows.Next() {
		var userid string
		var passwd string
		var uid int
		var gid int
		err = rows.Scan(&userid, &passwd, &uid, &gid)
		checkErr(err)
		fmt.Println(userid)
		fmt.Println(passwd)
		fmt.Println(uid)
		fmt.Println(gid)
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{
		"message": "Jokes handler not implemented yet",
	})
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
