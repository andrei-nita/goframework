package main

import (
	"context"
	"crypto/tls"
	_ "embed"
	"errors"
	"fmt"
	"github.com/NYTimes/gziphandler"
	fk "github.com/andrei-nita/goframework/framework"
	"github.com/andrei-nita/goframework/framework/auth"
	"github.com/andrei-nita/goframework/middleware"
	"github.com/andrei-nita/goframework/routes"
	"github.com/arl/statsviz"
	"github.com/gorilla/csrf"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	certFile = "localhost+2.pem"
	keyFile  = "localhost+2-key.pem"
)

var (
	httpServer  = &http.Server{}
	httpsServer = &http.Server{}
	csrfProtect func(http.Handler) http.Handler
)

func main() {
	// log
	err := fk.CreateLog(false)
	if err != nil {
		log.Fatalln(err)
	}
	defer fk.CloseLog()

	// recover
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in main", r)
		}
	}()

	// csrf key
	csrfProtect = csrf.Protect([]byte(fk.Server.CSRF))

	// auth
	err = auth.OpenAuth()
	if err != nil {
		log.Fatalln(err)
	}
	defer auth.CloseAuth()

	// create server
	mux := http.NewServeMux()
	httpServerSetup(mux)
	httpsServerSetup(mux)

	// stats
	statsviz.Register(mux)

	// configure ssl server
	if fk.Server.UseSSL {
		if fk.Server.Mode == fk.ModeProd {
			httpsServerSetupProduction()
		} else if fk.Server.Mode == fk.ModeDev {
			mkcertPath, err := fk.CreateMkcertIfNotExists()
			fk.LogFatalErr(err)
			err = fk.GenerateSSL(mkcertPath)
			fk.LogFatalErr(err)
			httpsServerSetupDevelopment()
		}
	}

	if fk.Server.CacheStaticFiles {
		mux.Handle("/static/public/", middleware.Cache(http.StripPrefix("/static/public/", http.FileServer(http.Dir("static/public")))))
	} else {
		mux.Handle("/static/public/", http.StripPrefix("/static/public/", http.FileServer(http.Dir("static/public"))))
	}

	routes.Setup(mux)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Println("http server closed", err)
		}
	}()

	if fk.Server.UseSSL {
		go func() {
			if err := httpsServer.ListenAndServeTLS("", ""); err != nil && errors.Is(err, http.ErrServerClosed) {
				log.Println("https server closed", err)
			}
		}()
	}

	// reload browser
	if fk.Server.Mode == fk.ModeDev {
		fk.ReloadBrowser()
		defer fk.CloseReloadBrowser()
	}

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default sends syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Println("http server shutdown", err)
	}

	if err := httpsServer.Shutdown(ctx); err != nil {
		log.Println("https server shutdown", err)
	}

	log.Println("main() done")
}

func handlerRedirect(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Connection", "close")
	url := fmt.Sprintf("https://%s%s%s", fk.Server.Domain, fk.Server.PortSSL, r.RequestURI)
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}

func httpServerSetup(mux *http.ServeMux) {
	httpServer.Addr = fk.Server.Port
	httpServer.ReadTimeout = 5 * time.Second
	httpServer.WriteTimeout = 10 * time.Second
	httpServer.Handler = gziphandler.GzipHandler(middleware.Secure(csrfProtect(middleware.IsAuthHandler(mux))))
}

func httpsServerSetup(mux *http.ServeMux) {
	httpsServer.Addr = fk.Server.PortSSL
	httpsServer.ReadTimeout = 5 * time.Second
	httpsServer.WriteTimeout = 10 * time.Second
	httpsServer.Handler = gziphandler.GzipHandler(middleware.Secure(csrfProtect(middleware.IsAuthHandler(mux))))
}

func httpsServerSetupDevelopment() {
	httpServer.Handler = http.HandlerFunc(handlerRedirect)
	cer, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalln(err)
	}
	httpsServer.TLSConfig = &tls.Config{Certificates: []tls.Certificate{cer}, MinVersion: tls.VersionTLS12}
}

func httpsServerSetupProduction() {
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(fmt.Sprint(fk.Server.Domain)), // Your domain here
		Cache:      autocert.DirCache("certs"),                           // Folder for storing certificates
	}
	httpServer.Handler = certManager.HTTPHandler(nil)
	httpsServer.TLSConfig = &tls.Config{GetCertificate: certManager.GetCertificate, MinVersion: tls.VersionTLS12}
}
