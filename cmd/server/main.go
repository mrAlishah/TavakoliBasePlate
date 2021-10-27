package main

import (
	"GolangTraining/internal/logger"
	"context"
	"flag"
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	GitCommit     string    = "Development"
	BuildTime     string    = time.Now().Format(time.RFC1123Z)
	ContainerName string    = "local"
	StartTime     time.Time = time.Now()
)

func main() {
	const op = "server.main"

	StartTime = time.Now()

	// Default Config file based on the environment variable
	defaultConfigFile := "configs/config-local.yaml"
	if env := os.Getenv("APP_MODE"); env != "" {
		defaultConfigFile = fmt.Sprintf("configs/config-%s.yaml", env)
	}
	logrus.Infof("loading config %s", defaultConfigFile)

	// Load Master Config File
	var configFile string
	flag.StringVar(&configFile, "c", defaultConfigFile, "The environment configuration file of application")
	flag.StringVar(&configFile, "config", defaultConfigFile, "The environment configuration file of application")
	flag.Usage = usage
	flag.Parse()

	// Print Start Ascii Art
	printAsciiArt()
	// Set Commit Metrics
	//gitMetrics()

	// Setting up the main context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// Loading the config file
	cfg, err := LoadConfig(configFile)
	if err != nil {
		logrus.Fatal(errors.Wrapf(err, "failed to load config: %s", op))
	}

	// Setup Logger
	loggy := logger.CreateLogger(cfg.Logger)

	// Get OS Container Name
	hostname, err := os.Hostname()
	if err != nil {
		loggy.Fatal(errors.WithMessage(err, op))
	}
	loggy.Infof("[OK] Hostname acquired :%s", hostname)
	loggy.Infof("Service Name: %s", cfg.Server)

	// Commit, BuildTime
	loggy.Infof("[OK] Commit Number:%s, Build Time: %s", GitCommit, BuildTime)

	logrus.Infof("loading config %s", defaultConfigFile)

	server := NewServer(cfg, loggy)

	// Initialize the Server Dependencies
	err = server.Initialize(ctx)
	if err != nil {
		loggy.Fatal(errors.Wrapf(err, "failed to initialize server: %s", op))
	}

	done := make(chan bool, 1)
	quiteSignal := make(chan os.Signal, 1)
	signal.Notify(quiteSignal, syscall.SIGINT, syscall.SIGTERM)

	// Graceful shutdown goroutine
	go server.GracefulShutdown(quiteSignal, done)

	// Start server in blocking mode
	if cfg.Server.Enabled {
		server.Start(ctx)
	}

	// Wait for graceful shutdown signal
	<-done

	// Kill other background jobs
	cancel()
	loggy.Info("Waiting for background jobs to finish their works...")

	// Wait for all other background jobs to finish their works
	server.Wait()

	loggy.Info("Master App Shutdown successfully, see you next time ;-)")
}

func usage() {
	usageStr := `
Usage: GolangTraining [options]
Options:
	-c,  --config   <config file name>   Path of yaml configuration file
`
	fmt.Printf("%s\n", usageStr)
	os.Exit(0)
}

func printAsciiArt() {
	// https://patorjk.com/software/taag/#p=display&f=ANSI%20Shadow&t=Project
	fmt.Println(aurora.Green(` GolangTraining PROJECT`))
}
