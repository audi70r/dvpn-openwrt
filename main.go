package main

import (
	"context"
	"embed"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/solarlabsteam/dvpn-openwrt/controllers"
	"github.com/solarlabsteam/dvpn-openwrt/services/auth"
	"github.com/solarlabsteam/dvpn-openwrt/services/dvpnconf"
	"github.com/solarlabsteam/dvpn-openwrt/services/keys"
	"github.com/solarlabsteam/dvpn-openwrt/services/node"
	"github.com/solarlabsteam/dvpn-openwrt/services/socket"
	"github.com/solarlabsteam/dvpn-openwrt/utilities/appconf"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
)

const version = "1.0.2"

//go:embed public
var public embed.FS

func main() {

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version":
			fmt.Printf("DVPN Interface version: %s\n", version)
			os.Exit(0)
		}
	}

	// load config
	appconf.LoadConf()

	// load configurations
	if confErr := dvpnconf.LoadConfig(); confErr != nil {
		panic(confErr)
	}

	// load sentinel key storage
	if err := keys.Load(appconf.Paths.SentinelPath()); err != nil {
		panic(err)
	}

	// load node
	node := node.New()

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

	// api group, that requires authorization
	api := r.PathPrefix("/api").Subrouter()
	api.Use(auth.Store.Authenticate)
	api.Path("/node/start/stream").HandlerFunc(node.StartNode).Methods("GET")
	api.Path("/node").HandlerFunc(node.GetNode).Methods("GET")
	api.Path("/node/kill").HandlerFunc(node.KillNode).Methods("POST")
	api.Path("/config").HandlerFunc(controllers.GetConfig).Methods("GET")
	api.Path("/config").HandlerFunc(controllers.PostConfig).Methods("POST")
	api.Path("/keys").HandlerFunc(controllers.ListKeys).Methods("GET")
	api.Path("/keys").HandlerFunc(controllers.AddRecoverKeys).Methods("POST")
	api.Path("/keys").HandlerFunc(controllers.DeleteKeys).Methods("DELETE")
	api.Path("/nat").HandlerFunc(controllers.GetNATInfo).Methods("GET")

	// api group, that does not require authorization
	r.HandleFunc("/api/socket", socket.Handle)
	r.Path("/api/login").HandlerFunc(controllers.Login).Methods("POST")

	// serve embedded static assets
	r.PathPrefix("/").Handler(publicFS)

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
