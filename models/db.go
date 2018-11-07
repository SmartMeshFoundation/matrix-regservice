package models

import (
	"fmt"

	"log"
	"os"

	"github.com/SmartMeshFoundation/matrix-regservice/params"
	"github.com/jinzhu/gorm"
	//_ "github.com/jinzhu/gorm/dialects/mysql"    //for gorm
	//_ "github.com/jinzhu/gorm/dialects/postgres" //for gorm
	_ "github.com/jinzhu/gorm/dialects/sqlite" //for gorm
)

var db *gorm.DB

//SetUpDB init db
func SetUpDB(path string) {
	var err error
	db, err = gorm.Open("sqlite3", path)
	if err != nil {
		panic("failed to connect database")
	}
	if params.DebugMode {
		db = db.Debug()
		db.LogMode(true)
	}
	//db.SetLogger(gorm.Logger{revel.TRACE})
	db.SetLogger(log.New(os.Stdout, "\r\n", 0))
	db.AutoMigrate(&User{})
	return
}

//CloseDB release connection
func CloseDB() {
	err := db.Close()
	if err != nil {
		log.Printf(fmt.Sprintf("closedb err %s", err))
	}
}

//SetupTestDB for test only
func SetupTestDB() {
	dbPath := "/tmp/test.db"
	err := os.Remove(dbPath)
	if err != nil {
		//ignore
	}
	SetUpDB(dbPath)
}
