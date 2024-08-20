package usecases

//import (
//	"fmt"
//
//	"github.com/rs/zerolog/log"
//
//	"github.com/berops/claudie/proto/pb"
//)
//
//// SaveConfigScheduler is a gRPC servie: the function saves config to the DB after receiving it from Scheduler
//func (u *Usecases) SaveConfigScheduler(request *pb.SaveConfigRequest) (*pb.SaveConfigResponse, error) {
//	config := request.GetConfig()
//
//	log.Info().Msgf("Saving config %s from Scheduler", config.Name)
//
//	// Save new config to the DB
//	config.DsChecksum = config.MsChecksum
//	config.SchedulerTTL = 0
//
//	if err := u.DB.UpdateDs(config); err != nil {
//		return nil, fmt.Errorf("error while updating dsChecksum for %s : %w", config.Name, err)
//	}
//
//	// Update workflow state for k8s clusters. (attached LB clusters included)
//	if err := u.DB.UpdateAllStates(config.Name, config.State); err != nil {
//		return nil, fmt.Errorf("error while saving workflow state config %s in MongoDB: %w", config.Name, err)
//	}
//
//	if err := u.DB.UpdateEvents(config.Name, config.Events); err != nil {
//		return nil, fmt.Errorf("error while saving config task events for %s in MongoDB: %w", config.Name, err)
//	}
//
//	if err := u.DB.UpdateSchedulerTTL(config.Name, config.SchedulerTTL); err != nil {
//		return nil, fmt.Errorf("error while updating schedulerTTL for %s : %w", config.Name, err)
//	}
//
//	log.Info().Msgf("Config %s successfully saved from Scheduler", config.Name)
//
//	return &pb.SaveConfigResponse{Config: config}, nil
//}
