package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/berops/claudie/internal/manifest"
	"github.com/berops/claudie/proto/pb"
	cbox "github.com/berops/claudie/services/context-box/client"
	"github.com/rs/zerolog/log"

	"gopkg.in/yaml.v3"
)

// processConfigs processes configs concurrently. If an error occurs while
// the file is being processed it's skipped and continues with
// the next one until all are processed. Nothing is done with
// files for which an error occurred, they'll be skipped until
// either corrected or deleted.
func (s *server) processConfigs() error {
	files, err := os.ReadDir(s.manifestDir)
	if err != nil {
		return fmt.Errorf("failed to read dir %q: %w", s.manifestDir, err)
	}

	configs, err := cbox.GetAllConfigs(s.cBox)
	if err != nil {
		return fmt.Errorf("failed to retrieve configs from context-box: %w", err)
	}

	log.Debug().Msgf("%d configs in database | %d files in %v", len(configs.Configs), len(files), s.manifestDir)

	type data struct {
		manifest    *manifest.Manifest
		rawManifest []byte
		path        string
		err         error
	}

	dataChan := make(chan *data, len(files))
	group := sync.WaitGroup{}

	for _, file := range files {
		group.Add(1)

		// Process each of the files concurrently
		// in a separate go-routine skipping over
		// file for which an error occurs.
		go func(entry os.DirEntry) {
			defer group.Done()

			path := filepath.Join(s.manifestDir, entry.Name())
			var rawManifest []byte
			var err error
			var m manifest.Manifest

			defer func() {
				dataChan <- &data{
					manifest:    &m,
					rawManifest: rawManifest,
					path:        path,
					err:         err,
				}
			}()

			if rawManifest, err = os.ReadFile(path); err != nil {
				return
			}

			if err = yaml.Unmarshal(rawManifest, &m); err != nil {
				return
			}

			err = m.Validate()
		}(file)
	}

	go func() {
		group.Wait()
		close(dataChan)
	}()

	// Collect data from files with no error.
	for data := range dataChan {
		var removed bool
		configs.Configs, removed = remove(configs.Configs, data.manifest.Name)

		if data.err != nil {
			log.Error().Msgf("Skipping over file %v due to error : %v", data.path, data.err)
			continue
		}

		cfg := &pb.Config{
			Name:     data.manifest.Name,
			Manifest: string(data.rawManifest),
		}

		if _, err := cbox.SaveConfigFrontEnd(s.cBox, &pb.SaveConfigRequest{Config: cfg}); err != nil {
			log.Error().Msgf("Skip saving config %v due to error : %v", data.manifest.Name, err)
			continue
		}

		log.Info().Msgf("File %s has been saved to the database", data.path)

		// if the config is not in the DB we start to track it.
		if !removed {
			for _, cluster := range data.manifest.Kubernetes.Clusters {
				if _, ok := s.inProgress.Load(cluster.Name); !ok {
					s.inProgress.Store(cluster.Name, cfg)
				}
			}
		}
	}

	for _, config := range configs.Configs {
		if _, ok := s.deletingConfigs.Load(config.Id); ok {
			continue
		}

		s.deletingConfigs.Store(config.Id, nil)

		for _, cluster := range config.GetCurrentState().GetClusters() {
			if _, ok := s.inProgress.Load(cluster.ClusterInfo.Name); !ok {
				s.inProgress.Store(cluster.ClusterInfo.Name, config)
			}
		}

		go func(config *pb.Config) {
			log.Info().Msgf("Deleting config: %v", config.Name)

			if err := cbox.DeleteConfig(s.cBox, &pb.DeleteConfigRequest{Id: config.Id, Type: pb.IdType_HASH}); err != nil {
				log.Error().Msgf("Failed to the delete %s with id %s : %v", config.Name, config.Id, err)
			}
			s.deletingConfigs.Delete(config.Id)
		}(config)
	}

	return nil
}

// remove deletes the config with the specified name from the slice.
// If not present the original slice is returned.
func remove(configs []*pb.Config, configName string) ([]*pb.Config, bool) {
	for index, config := range configs {
		if config.Name == configName {
			configs = append(configs[0:index], configs[index+1:]...)
			return configs, true
		}
	}
	return configs, false
}
