package ports

import "github.com/berops/claudie/proto/pb"

type MongoDBPort interface {
	GetConfig(id string, idType pb.IdType) (*pb.Config, error)
	DeleteConfig(id string, idType pb.IdType) error
	GetAllConfigs() ([]*pb.Config, error)
	SaveConfig(config *pb.Config) error
	UpdateSchedulerTTL(name string, newTTL int32) error
	UpdateBuilderTTL(name string, newTTL int32) error
	UpdateMsToNull(id string) error
	UpdateDs(config *pb.Config) error
	UpdateCs(config *pb.Config) error
	UpdateWorkflowState(configName, clusterName string, workflow *pb.Workflow) error
}
