package main

import (
	"context"
	"embed"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/solarlabsteam/dvpn-openwrt/controllers"
	"github.com/solarlabsteam/dvpn-openwrt/services/keys"
	"github.com/solarlabsteam/dvpn-openwrt/services/socket"
	"github.com/solarlabsteam/dvpn-openwrt/utilities/appconf"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
)

//go:embed public

var public embed.FS

func main() {
	// load config
	appconf.LoadConf()

	// load sentinel key storage
	if err := keys.Load(appconf.Paths.SentinelDir); err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	if _, homeSet := os.LookupEnv("HOME"); !homeSet {
		os.Setenv("PATH", appconf.Paths.BinDir)
		os.Setenv("HOME", appconf.Paths.HomeDir)
	}
	// for development: serve static assets from public folder
	//publicFS := http.FileServer(http.Dir("./public"))

	// for production: embed static assets into binary
	publicDir, _ := fs.Sub(public, "public")
	publicFS := http.FileServer(http.FS(publicDir))

	r.PathPrefix("/").Handler(publicFS) // serve embedded static assets
	r.HandleFunc("/api/node/start/stream", controllers.StartNodeStreamStd)
	r.HandleFunc("/api/node", controllers.GetNode)
	r.HandleFunc("/api/node/kill", controllers.KillNode)
	r.HandleFunc("/api/config", controllers.Config)
	r.HandleFunc("/api/socket", socket.Handle)
	r.HandleFunc("/api/keys", controllers.ListKeys)
	r.HandleFunc("/api/keys/add", controllers.AddRecoverKeys)

	srv := &http.Server{
		Addr:         fmt.Sprintf("%v:%v", appconf.Server.Addr, appconf.Server.Port),
		WriteTimeout: appconf.Server.WriteTimeout,
		ReadTimeout:  appconf.Server.ReadTimeout,
		IdleTimeout:  appconf.Server.IdleTimeout,
		Handler:      r,
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), appconf.Server.HttpServerGracefulShutdownTimeout)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
