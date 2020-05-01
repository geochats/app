package main

import (
	"fmt"
	"geochats/pkg/bot"
	"geochats/pkg/client"
	"geochats/pkg/downloader"
	"geochats/pkg/loaders"
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

		tmpDir    = ensureEnv("TMP_DIR")
		publicDir = ensureEnv("PUBLIC_DIR")

		pgDsn = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			ensureEnv("DB_USER"),
			ensureEnv("DB_PASS"),
			ensureEnv("DB_HOST"),
			ensureEnv("DB_PORT"),
			ensureEnv("DB_NAME"),
		)

		tgAppID     = ensureEnv("TG_APP_ID")
		tgAppHash   = ensureEnv("TG_APP_HASH")
		botApiToken = os.Getenv("BOT_API_TOKEN")
		botDisabled = os.Getenv("BOT_DISABLED") != ""

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
		"tmpDir":      tmpDir,
		"publicDir":   publicDir,
		"verbose":     verbose,
		"tgAppID":     tgAppID,
		"tgAppHash":   tgAppHash,
		"botApiToken": botApiToken,
		"listen":      listen,
		"pgDsn":       pgDsn,
	})

	cl, err := client.New(tgAppID, tgAppHash, tmpDir)
	if err != nil {
		log.Panicf("can't create tg client: %v", err)
	}
	if err := client.EnsureBotAuth(cl, botApiToken, 10, 2*time.Second); err != nil {
		log.Panicf("can't auth tg bot: %v", err)
	}

	store, err := storage.New(pgDsn)
	if err != nil {
		log.Panicf("can't create storage: %v", err)
	}

	dl := downloader.NewSyncDownloader(cl, fmt.Sprintf("%s/c", publicDir), "/c")
	loader := loaders.NewChannelInfoLoader(cl, fmt.Sprintf("%s/c", publicDir), "/c")

	if !botDisabled {
		b := bot.New(cl, store, loader, dl, logger)
		go func() {
			if err := b.Run(); err != nil {
				log.Fatalf("error in bot run: %v", err)
			}
		}()
	}

	srv := web_server.New(listen, publicDir, store, logger)
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
