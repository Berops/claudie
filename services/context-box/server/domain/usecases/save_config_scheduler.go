package usecases

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/berops/claudie/proto/pb"
)

// SaveConfigScheduler is a gRPC servie: the function saves config to the DB after receiving it from Scheduler
func (u *Usecases) SaveConfigScheduler(request *pb.SaveConfigRequest) (*pb.SaveConfigResponse, error) {
	config := request.GetConfig()
	log.Info().Msgf("Saving config %s from Scheduler", config.Name)
	// Save new config to the DB
	config.DsChecksum = config.MsChecksum
	config.SchedulerTTL = 0
	err := u.DB.UpdateDs(config)
	if err != nil {
		return nil, fmt.Errorf("error while updating dsChecksum for %s : %w", config.Name, err)
	}

	if config.DesiredState != nil {
		// Update workflow state for K8s clusters. (attached LB clusters included)
		for _, cluster := range config.DesiredState.Clusters {
			if err := u.DB.UpdateWorkflowState(config.Name, cluster.ClusterInfo.Name, config.State[cluster.ClusterInfo.Name]); err != nil {
				return nil, fmt.Errorf("error while updating workflow state for k8s cluster %s in config %s : %w", cluster.ClusterInfo.Name, config.Name, err)
			}
		}
	}
	if err := u.DB.UpdateSchedulerTTL(config.Name, config.SchedulerTTL); err != nil {
		return nil, fmt.Errorf("error while updating schedulerTTL for %s : %w", config.Name, err)
	}

	log.Info().Msgf("Config %s successfully saved from Scheduler", config.Name)
	return &pb.SaveConfigResponse{Config: config}, nil
}
