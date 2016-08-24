package main

import (
	"net/http"
	"os"
	"time"

	"github.com/boltdb/bolt"
	kitlog "github.com/go-kit/kit/log"
	"golang.org/x/net/context"

	"github.com/gnhuy91/profilesvc"
	boltsvc "github.com/gnhuy91/profilesvc/bolt"
)

func main() {
	var logger kitlog.Logger
	{
		logger = kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))
		logger = kitlog.NewContext(logger).With("ts", kitlog.DefaultTimestampUTC)
		logger = kitlog.NewContext(logger).With("caller", kitlog.DefaultCaller)
	}

	// Local bolt DB
	db, err := bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		logger.Log("bolt", err)
	}
	defer db.Close()

	// Create "profiles" bucket
	err = func(db *bolt.DB, name string) error {
		tx, err := db.Begin(true)
		if err != nil {
			return err
		}
		defer tx.Rollback()

		_, err = tx.CreateBucket([]byte(name))
		if err != nil {
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
		return nil
	}(db, "profiles")

	if err != nil {
		logger.Log("bolt", err)
	}

	var svc profilesvc.Service
	{
		svc = &boltsvc.ProfileService{DB: db}
	}

	ctx := context.Background()
	h := profilesvc.MakeHTTPHandler(ctx, svc, logger)
	logger.Log("exit", http.ListenAndServe(":8080", h))
}