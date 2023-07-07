package usecases

import (
	"fmt"
	"sync"

	"github.com/berops/claudie/internal/utils"
	"github.com/berops/claudie/proto/pb"
)

// ConfigProcessor will fetch new config from the context-box microservice.
// Each received config will be processed in a separate go-routine.
// If a sync.WaitGroup is supplied it will call the Add(1) and then the Done() method on it after
// the go-routine finishes the work, if nil it will be ignored.
func (u *Usecases) ConfigProcessor(contextBoxGrpcClient pb.ContextBoxServiceClient, waitGroup *sync.WaitGroup) error {
	// Pull an item from the context-box scheduler queue
	response, err := u.ContextBox.GetConfigScheduler(contextBoxGrpcClient)
	if err != nil {
		return fmt.Errorf("error while getting Scheduler config: %w", err)
	}

	config := response.GetConfig()
	if config == nil {
		return nil
	}

	if waitGroup != nil {
		// We received a non-nil config thus we add a new worker to the wait group.
		waitGroup.Add(1)
	}

	go func() {
		if waitGroup != nil {
			defer waitGroup.Done()
		}

		logger := utils.CreateLoggerWithProjectName(config.Name)

		logger.Info().Msgf("Processing config")

		// Process (build desired state) the config
		if configProcessingErr := u.processConfig(config, contextBoxGrpcClient); configProcessingErr != nil {
			logger.Err(configProcessingErr).Msgf("Error while processing config")

			// Save processing error message to config
			if err := u.saveErrorMessageToConfig(config, contextBoxGrpcClient, configProcessingErr); err != nil {
				logger.Err(err).Msgf("Failed to save error to the config")
			}
		}
		logger.Info().Msgf("Config have been successfully processed")
	}()

	return nil
}

// processConfig contains the core logic of processing a config
// returns error if not successful, nil otherwise
func (u *Usecases) processConfig(config *pb.Config, contextBoxGrpcClient pb.ContextBoxServiceClient) error {
	// Create desired state
	config, err := u.CreateDesiredState(config)
	if err != nil {
		return fmt.Errorf("error while creating a desired state: %w", err)
	}

	// After constructing the desired state for the config
	// save it to the context-box DB
	err = u.ContextBox.SaveConfigScheduler(config, contextBoxGrpcClient)
	if err != nil {
		return fmt.Errorf("error while saving the config: %w", err)
	}

	return nil
}

// saveErrorMessageToConfig saves error message to the config
// Returns error if not successful, nil otherwise
func (u *Usecases) saveErrorMessageToConfig(config *pb.Config, contextBoxGrpcClient pb.ContextBoxServiceClient, err error) error {
	if config.State == nil {
		// If no state map created yet, set the error under the key "scheduler".
		config.State = make(map[string]*pb.Workflow)
		config.State["scheduler"] = &pb.Workflow{
			Stage:       pb.Workflow_SCHEDULER,
			Status:      pb.Workflow_ERROR,
			Description: err.Error(),
		}
	} else {
		// As scheduler does not process config on per cluster basis, the error will be saved under every cluster.
		for _, state := range config.State {
			state.Stage = pb.Workflow_SCHEDULER
			state.Status = pb.Workflow_ERROR
			state.Description = fmt.Sprintf("Error encountered in Scheduler, halting processing of all clusters : %v", err)
		}
	}

	return u.ContextBox.SaveConfigScheduler(config, contextBoxGrpcClient)
}
