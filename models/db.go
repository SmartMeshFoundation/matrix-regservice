package models

import (
	"encoding/gob"
	"fmt"
	"os"
	"time"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/asdine/storm"
	gobcodec "github.com/asdine/storm/codec/gob"
	"github.com/coreos/bbolt"
	"github.com/ethereum/go-ethereum/common"
)

var db *storm.DB

//SetupDB init db
func SetupDB(dbPath string) error {
	var err error
	log.Trace(fmt.Sprintf("dbpath=%s", dbPath))
	needCreateDb := !common.FileExist(dbPath)
	db, err = storm.Open(dbPath, storm.BoltOptions(os.ModePerm, &bolt.Options{Timeout: 1 * time.Second}), storm.Codec(gobcodec.Codec))
	if err != nil {
		return err
	}
	initDB()
	if needCreateDb {

	}
	return nil
}
func init() {
	gob.Register(&User{})
}
func initDB() {
	err := db.Init(&User{})
	if err != nil {
		log.Error(fmt.Sprintf("db err %s", err))
	}
}
