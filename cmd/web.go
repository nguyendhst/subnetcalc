/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nguyendhst/subnetcalc/calc"
	"github.com/spf13/cobra"
)

var PORT string
var SERVER string

// isV4 is alias for calc.VerifyIPv4
var isV4 = calc.VerifyIPv4

// isV6 is alias for calc.VerifyIPv6
var isV6 = calc.VerifyIPv6

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start web server",
	Long:  `Start web server to serve HTML UI, user can specify port and server address`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		router := gin.Default()

		router.LoadHTMLGlob("templates/*")
		router.GET("/", func(c *gin.Context) {
			time.Sleep(10 * time.Second)
			c.HTML(http.StatusOK, "index.html", nil)
		})
		router.POST("/result", func(c *gin.Context) {
			var input calc.IPInput
			addr := c.PostForm("addr")
			if ok, err := isV4(addr); ok && err == nil {
				input = &calc.IPv4{Addr: addr}
			} else if ok, err := isV6(addr); ok && err == nil {
				input = &calc.IPv6{Addr: addr}
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid IP address"})
				return
			}
			res, err := calc.ProcessInput(input)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, res)
		})

		srv := &http.Server{
			Addr:    SERVER + ":" + PORT,
			Handler: router,
		}

		// Initializing the server in a goroutine so that
		// it won't block the graceful shutdown handling below
		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen: %s\n", err)
			}
		}()

		// Listen for the interrupt signal.
		<-ctx.Done()

		// Restore default behavior on the interrupt signal and notify user of shutdown.
		stop()
		log.Println("shutting down gracefully, press Ctrl+C again to force")

		// The context is used to inform the server it has 5 seconds to finish
		// the request it is currently handling
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("Server forced to shutdown: ", err)
		}

		log.Println("Server exiting")
	},
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
