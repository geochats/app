package main

import (
	"fmt"
	"geochats/pkg/bot"
	"geochats/pkg/client"
	"geochats/pkg/collector"
	"geochats/pkg/collector/loaders"
	"geochats/pkg/downloader"
	"geochats/pkg/storage"
	"geochats/pkg/web_server"
	"github.com/Arman92/go-tdlib"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

func main() {
	var (
		verbose = os.Getenv("VERBOSE") != ""

		workDir   = ensureEnv("VAR_DIR")
		publicDir = ensureEnv("PUBLIC_DIR")
		dbFile    = ensureEnv("DB_FILE")

		tgAppID     = ensureEnv("TG_APP_ID")
		tgAppHash   = ensureEnv("TG_APP_HASH")
		botApiToken = os.Getenv("BOT_API_TOKEN")

		listen = ensureEnv("LISTEN")

		logger = log.StandardLogger()
	)

	log.SetLevel(log.InfoLevel)
	tdlib.SetLogVerbosityLevel(1)
	if verbose {
		log.SetLevel(log.DebugLevel)
		tdlib.SetLogVerbosityLevel(5)
	}
	logger.Infof("Options: %#v", log.Fields{
		"workDir":     workDir,
		"publicDir":   publicDir,
		"verbose":     verbose,
		"tgAppID":     tgAppID,
		"tgAppHash":   tgAppHash,
		"botApiToken": botApiToken,
		"listen":      listen,
		"dbFile":      dbFile,
	})

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

	srv := web_server.New(listen, publicDir, cl, store, loader, logger)
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
