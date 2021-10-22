package testingframework

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/Berops/platform/proto/pb"
	cbox "github.com/Berops/platform/services/context-box/client"
	"github.com/Berops/platform/urls"
	"github.com/Berops/platform/utils"
	"github.com/rs/zerolog/log"

	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testDir    = "tests"
	maxTimeout = 3600 // max allowed time for operation to finish in [seconds]
)

// ClientConnection will return new client connection to Context-box
func ClientConnection() pb.ContextBoxServiceClient {
	cc, err := utils.GrpcDialWithInsecure("context-box", urls.ContextBoxURL)
	if err != nil {
		log.Fatal().Err(err)
	}

	// Creating the client
	c := pb.NewContextBoxServiceClient(cc)
	return c
}

// TestPlatform will start all the test cases specified in tests directory
func TestPlatform(t *testing.T) {
	var err error
	c := ClientConnection()
	log.Info().Msg("----Starting the tests----")

	// loop through the directory and list files inside
	files, err := ioutil.ReadDir(testDir)
	if err != nil {
		log.Fatal().Msgf("Error while trying to read test sets: %v", err)
	}

	// save all the test set paths
	var pathsToSets []string
	for _, f := range files {
		if f.IsDir() {
			log.Info().Msgf("Found test set: %s", f.Name())
			setTestDir := filepath.Join(testDir, f.Name())
			pathsToSets = append(pathsToSets, setTestDir)
		}
	}

	// apply test sets sequentially - while framework is still in dev
	for _, path := range pathsToSets {
		err = applyTestSet(path, c)
		if err != nil {
			log.Fatal().Msgf("Error while processing %s : %v", path, err)
			break
		}
	}

	require.NoError(t, err)
}

// applyTestSet function will apply test set sequantially to a platform
func applyTestSet(pathToSet string, c pb.ContextBoxServiceClient) error {
	done := make(chan string)
	var id string

	log.Info().Msgf("Working on the test set: %s", pathToSet)

	files, err := ioutil.ReadDir(pathToSet)
	if err != nil {
		log.Fatal().Msgf("Error while trying to read test configs: %v", err)
	}

	for _, file := range files {
		setFile := filepath.Join(pathToSet, file.Name())
		manifest, errR := ioutil.ReadFile(setFile)
		if errR != nil {
			log.Fatal().Err(errR)
		}

		id, err = cbox.SaveConfigFrontEnd(c, &pb.SaveConfigRequest{
			Config: &pb.Config{
				Name:     file.Name(),
				Id:       id,
				Manifest: string(manifest),
			},
		})

		if err != nil {
			log.Fatal().Msgf("Error while saving a config: %v", err)
			return err
		}
		go configChecker(done, c, id, file.Name())
		// wait until test config has been processed
		if res := <-done; res != "ok" {
			log.Error().Msg(res)
			return fmt.Errorf(res)
		}
	}
	// delete the nodes
	log.Info().Msgf("Deleting the clusters from test set: %s", pathToSet)
	err = cbox.DeleteConfig(c, id)
	if err != nil {
		return err
	}

	return nil
}

// configChecker function will check if the config has been applied every 30s
func configChecker(done chan string, c pb.ContextBoxServiceClient, configID string, configName string) {
	var counter int
	sleepSec := 30
	for {
		elapsedSec := counter * sleepSec
		// if CSchecksum == DSchecksum, the config has been processed
		config, err := c.GetConfigById(context.Background(), &pb.GetConfigByIdRequest{
			Id: configID,
		})
		if err != nil {
			emsg := fmt.Sprintf("Got error while waiting for config to finish: %v", err)
			log.Fatal().Msg(emsg)
		}
		if config != nil {
			cfg := config.Config
			if len(cfg.ErrorMessage) > 0 {
				emsg := cfg.ErrorMessage
				log.Error().Msg(emsg)
				done <- emsg
				return
			}
			currentState := cfg.CurrentState.String()
			desiredState := cfg.DesiredState.String()
			// if currentState == desiredState, the config has been processed
			if len(currentState) > 0 && currentState == desiredState {
				break
			}
		}
		if elapsedSec == maxTimeout {
			emsg := fmt.Sprintf("Test took too long... Aborting on timeout %d seconds", maxTimeout)
			log.Fatal().Msg(emsg)
			done <- emsg
			return
		}
		time.Sleep(time.Duration(sleepSec) * time.Second)
		counter++
		log.Info().Msgf("Waiting for %s to finish... [ %ds elapsed ]", configName, elapsedSec)
	}
	// send signal that config has been processed, unblock the applyTestSet
	done <- "ok"
}
