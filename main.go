package main

import (
	"fmt"
	"log"

	"conduit/api"
	"conduit/config"
	db "conduit/db/sqlc"
)

func main() {
	conf, err := config.LoadConfig("dev", "./env")
	if err != nil {
		log.Fatal("config loading error:", err)
	}

	conn, err := db.Connect(conf)
	if err != nil {
		log.Fatal("db connection error:", err)
	}
	defer db.Close(conn)

	if err := db.AutoMigrate(conf); err != nil {
		log.Fatal("db migration error:", err)
	}

	store := db.New(conn)
	server := api.NewServer(conf, store)
	server.MountHandlers()

	addr := fmt.Sprintf(":%s", conf.Port)
	if err := server.Start(addr); err != nil {
		log.Fatal("server error:", err)
	}
}
