package main

import (
	"fmt"
	"geochats/pkg/bot"
	"geochats/pkg/client"
	"geochats/pkg/client/downloader"
	"geochats/pkg/collector"
	"geochats/pkg/collector/loaders"
	"geochats/pkg/storage"
	"geochats/pkg/web_server"
	"github.com/Arman92/go-tdlib"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"time"
)

func main() {
	log.SetLevel(log.InfoLevel)
	workDir := ensureEnv("WORK_DIR")
	publicDir := ensureEnv("PUBLIC_DIR")
	verbose := os.Getenv("VERBOSE") != ""
	tgAppID := ensureEnv("TG_APP_ID")
	tgAppHash := ensureEnv("TG_APP_HASH")
	botApiToken := os.Getenv("BOT_API_TOKEN")
	listen := ensureEnv("LISTEN")
	dbFile := ensureEnv("DB")
	logger := log.StandardLogger()

	tdlib.SetLogVerbosityLevel(1)
	if verbose {
		log.SetLevel(log.DebugLevel)
		tdlib.SetLogVerbosityLevel(5)
	}
	if workDir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		workDir, err = filepath.EvalSymlinks(cwd)
		if err != nil {
			log.Fatal(err)
		}
	}

	cl, err := client.New(tgAppID, tgAppHash, workDir)
	if err != nil {
		log.Panicf("can't create tg client: %v", err)
	}
	if err := client.EnsureBotAuth(cl, botApiToken, 10, 2*time.Second); err != nil {
		log.Panicf("can't auth tg bot: %v", err)
	}

	store, err := storage.New(dbFile)
	if err != nil {
		log.Panicf("can't create storage: %v", err)
	}

	dl := downloader.NewSyncDownloader(cl, fmt.Sprintf("%s/c", publicDir), "/c")
	loader := loaders.NewChannelInfoLoader(cl, fmt.Sprintf("%s/c", publicDir), "/c")

	b := bot.New(cl, store, loader, dl, logger)
	go func() {
		if err := b.Run(); err != nil {
			log.Fatalf("error in bot run: %v", err)
		}
	}()

	col := collector.New(cl, loader, store, logger)
	go func() {
		for {
			if err := col.UpdateGroups(); err != nil {
				log.Errorf("can't update groups: %v", err)
			}
			time.Sleep(5 * time.Second)
		}
	}()

	srv := web_server.New(listen, cl, store, loader, logger)
	if err := srv.Listen(); err != nil {
		log.Panicf("can't create storage: %v", err)
	}
}

func ensureEnv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		log.Panicf("env variable `%s` is empty", name)
	}
	return v
}
