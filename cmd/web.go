/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

var PORT string
var SERVER string
var templates = make(map[string]*template.Template, 2)
var templateName = []string{
	"layout.html",
	"index.html",
}
var templateDir = "./templates/"

// webCmd represents the web command
var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start web server",
	Long:  `Start web server to serve HTML UI, user can specify port and server address`,
	Run: func(cmd *cobra.Command, args []string) {
		rMux := mux.NewRouter()
		loadTemplates()
		s := &http.Server{
			Addr:         fmt.Sprintf("%s:%s", SERVER, PORT),
			Handler:      rMux,
			ErrorLog:     nil,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout:  10 * time.Second,
		}

		rMux.HandleFunc("/", DefaultHandler)

		go func() {
			fmt.Printf("Listening on %s:%s\n", SERVER, PORT)
			err := s.ListenAndServe()
			if err != nil {
				fmt.Println(err)
				Halt(err)
			} else {
				fmt.Println("Visit " + SERVER + ":" + PORT)
			}
		}()

		waitSignal := make(chan os.Signal, 1)
		signal.Notify(waitSignal, os.Interrupt)
		sig := <-waitSignal
		fmt.Println("Got signal:", sig)
		time.Sleep(2 * time.Second)
		s.Shutdown(context.TODO())
	},
}

func loadTemplates() {
	for _, name := range templateName {
		t, err := template.ParseFiles(templateDir + name)
		if err != nil {
			fmt.Println(err)
			Halt(err)
		}
		templates[name] = t
	}
}

// DefaultHandler handles default request
func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	err := templates["layout.html"].Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		Halt(err)
	}
}

func init() {
	rootCmd.AddCommand(webCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// webCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	webCmd.Flags().StringVarP(&PORT, "port", "p", "63000", "Set port number")
	webCmd.Flags().StringVarP(&SERVER, "server", "s", "127.0.0.1", "Set server address")
}
