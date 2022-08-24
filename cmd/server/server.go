package main

import (
	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
	"introvert_test_mission/internal/server"
	"time"
)

func main() {

	var opts struct {
		Host string `long:"host" env:"HOST" default:"0.0.0.0"`
		Port string `long:"port" env:"PORT" default:"8080"`

		MongoURL        string `long:"mongo_url" env:"MONGO_URL" default:"mongodb://root:12345@0.0.0.0:27017"`
		MongoCollection string `long:"mongo_collection" env:"MONGO_COLLECTION" default:"entities"`
		MongoDBName     string `long:"mongo_db_name" env:"MONGO_DB_NAME" default:"test"`

		CacheUpdatePeriod int `long:"cache_update_period" env:"CACHE_UPDATE_PERIOD" default:"1" help:"Update cache period, in seconds"`
	}

	if _, err := flags.Parse(&opts); err != nil {
		log.Fatalf("Failed to parse args: %s", err)
	}

	mdb, err := server.NewMongoDB(opts.MongoURL, opts.MongoDBName, opts.MongoCollection)
	if err != nil {
		log.Fatal(err)
	}
	defer mdb.Close()

	cache := server.NewCache(time.Duration(opts.CacheUpdatePeriod)*time.Second, mdb)

	if err := server.Run(opts.Host, opts.Port, mdb, cache); err != nil {
		log.Fatal(err)
	}

}
