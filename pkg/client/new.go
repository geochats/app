package client

import (
	"fmt"
	"github.com/Arman92/go-tdlib"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

func New(apiID string, apiHash string, baseDir string) (AbstractClient, error) {
	if apiID == "" {
		return nil, fmt.Errorf("api")
	}

	databaseDir := baseDir + "/var/db"
	if err := os.MkdirAll(databaseDir, 0777); err != nil {
		return nil, fmt.Errorf("can't create database dir: %v", err)
	}

	filesDir := baseDir + "/var/files"
	if err := os.RemoveAll(filesDir); err != nil {
		return nil, fmt.Errorf("can't delete files dir: %v", err)
	}
	if err := os.MkdirAll(filesDir, 0777); err != nil {
		return nil, fmt.Errorf("can't create files dir: %v", err)
	}

	// Create new instance of client
	client := tdlib.NewClient(tdlib.Config{
		APIID:               apiID,
		APIHash:             apiHash,
		SystemLanguageCode:  "en",
		DeviceModel:         "Server",
		SystemVersion:       "1.0.0",
		ApplicationVersion:  "1.0.0",
		UseMessageDatabase:  false,
		UseSecretChats:      false,
		UseFileDatabase:     true,
		UseChatInfoDatabase: true,
		UseTestDataCenter:   false,
		DatabaseDirectory:   databaseDir,
		FileDirectory:       filesDir,
		IgnoreFileNames:     false,
	})
	return client, nil
}

func EnsureBotAuth(client AbstractClient, botToken string, attempts int, sleep time.Duration) error {
	for i := 0; i < attempts; i++ {
		currentState, _ := client.Authorize()
		log.Infof(
			"Trying to auth like a bot. Current state: %#v %#v %v",
			currentState.GetAuthorizationStateEnum(),
			tdlib.AuthorizationStateReadyType,
			currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateReadyType,
		)
		if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitPhoneNumberType {
			_, err := client.CheckAuthenticationBotToken(botToken)
			if err != nil {
				return fmt.Errorf("can't check bot auth token: %v", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateReadyType {
			return nil
		}
		time.Sleep(sleep)
	}
	return fmt.Errorf("bot auth stuck")
}
