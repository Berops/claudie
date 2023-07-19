package ports

import "github.com/berops/claudie/proto/pb"

type ContextBoxPort interface {
	GetAllConfigs() ([]*pb.Config, error)
	SaveConfig(config *pb.Config) error
	DeleteConfig(configName string) error
	PerformHealthCheck() error
}
