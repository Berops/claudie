package testingframework

import (
	"context"
	"fmt"
	"time"

	comm "github.com/berops/claudie/internal/command"
	"github.com/berops/claudie/internal/kubectl"
	"github.com/berops/claudie/internal/utils"
	"github.com/berops/claudie/proto/pb"
	"github.com/berops/claudie/proto/pb/spec"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

const (
	// Deployment which should NOT trigger scale up
	scaleUpDeploymentIgnored = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment-ignore
  labels:
    app: nginx
spec:
  replicas: 4
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - name: nginx
          image: nginx:1.14.2
          ports:
            - containerPort: 80
          resources:
            requests:
              memory: 8000Mi`
	// Deployment which should trigger scale up
	scaleUpDeployment = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment-success
  labels:
    app: nginx
spec:
  replicas: 4
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - name: nginx
          image: nginx:1.14.2
          ports:
            - containerPort: 80
          resources:
            requests:
              memory: 500Mi`
	// Time in which Autoscaler should trigger scale up
	scaleUpTimeout = 3
	// Time in which Autoscaler should trigger scale down
	scaleDownTimeout = 15
)

// testAutoscaler tests the Autoscaler deployment.
func testAutoscaler(ctx context.Context, config *spec.Config) error {
	autoscaledClusters := getAutoscaledClusters(config)
	if len(autoscaledClusters) == 0 {
		// No clusters are currently autoscaled.
		return testLonghornDeployment(ctx, config)
	}

	c, cc := clientConnection()
	defer func() {
		err := cc.Close()
		if err != nil {
			log.Err(err).Msgf("error while closing client connection")
		}
	}()

	var clusterGroup errgroup.Group
	for _, cluster := range autoscaledClusters {
		func(cluster *spec.K8Scluster) {
			clusterGroup.Go(
				func() error {
					log.Info().Msgf("Deploying pods which should be ignored by autoscaler for cluster %s", cluster.ClusterInfo.Name)
					return applyDeployment(cluster, scaleUpDeploymentIgnored)
				})
		}(cluster)
	}

	if err := clusterGroup.Wait(); err != nil {
		return fmt.Errorf("failed to deploy scale up deployment which should be ignored: %w", err)
	}
	// Wait before checking for changes
	log.Info().Msgf("Waiting %d minutes to see if autoscaler starts the scale up", scaleUpTimeout)
	time.Sleep(scaleUpTimeout * time.Minute)

	// Check if build has been started, if yes, error.
	if res, err := c.GetConfigFromDB(context.Background(), &pb.GetConfigFromDBRequest{Id: config.Id, Type: pb.IdType_HASH}); err == nil {
		if !checksumsEqual(res.Config.DsChecksum, res.Config.CsChecksum) {
			return fmt.Errorf("some cluster/s in config %s have been scaled up, when they should not", config.Name)
		} else {
			log.Info().Msgf("Config %s has successfully passed autoscaling test [1/3]", config.Name)
		}
	} else {
		return fmt.Errorf("error while retrieving config %s from DB : %w", config.Name, err)
	}
	// Apply scale up deployment.
	for _, cluster := range autoscaledClusters {
		func(cluster *spec.K8Scluster) {
			clusterGroup.Go(
				func() error {
					log.Info().Msgf("Deploying pods which should trigger scale up by autoscaler for cluster %s", cluster.ClusterInfo.Name)
					return applyDeployment(cluster, scaleUpDeployment)
				})
		}(cluster)
	}

	if err := clusterGroup.Wait(); err != nil {
		return fmt.Errorf("failed to deploy scale up deployment : %w", err)
	}

	// Wait before checking for changes.
	log.Info().Msgf("Waiting %d minutes to see if autoscaler starts the scale up", scaleUpTimeout)
	time.Sleep(scaleUpTimeout * time.Minute)

	// Check if build has been started, if no, error (Scale up).
	if res, err := c.GetConfigFromDB(context.Background(), &pb.GetConfigFromDBRequest{Id: config.Id, Type: pb.IdType_HASH}); err == nil {
		if checksumsEqual(res.Config.DsChecksum, res.Config.CsChecksum) {
			return fmt.Errorf("some cluster/s in config %s have not been scaled up, when they should have", config.Name)
		} else {
			log.Info().Msgf("Config %s has successfully passed autoscaling test [2/3]", config.Name)
		}
	} else {
		return fmt.Errorf("error while retrieving config %s from DB : %w", config.Name, err)
	}

	// Wait until build is finished.
	if err := configChecker(ctx, c, "autoscaling", "scale-up-test", idInfo{id: config.Id, idType: pb.IdType_HASH}); err != nil {
		return err
	}
	// Test longhorn.
	// Get new config from DB with updated counts.
	if res, err := c.GetConfigFromDB(context.Background(), &pb.GetConfigFromDBRequest{Id: config.Id, Type: pb.IdType_HASH}); err != nil {
		return err
	} else {
		if err := testLonghornDeployment(ctx, res.Config); err != nil {
			return err
		}
	}

	for _, cluster := range autoscaledClusters {
		func(cluster *spec.K8Scluster) {
			clusterGroup.Go(
				func() error {
					log.Info().Msgf("Removing pods which should trigger scale down by autoscaler for cluster %s", cluster.ClusterInfo.Name)
					return removeDeployment(cluster, scaleUpDeployment)
				})
		}(cluster)
	}
	if err := clusterGroup.Wait(); err != nil {
		return fmt.Errorf("failed to remove scale up deployment : %w", err)
	}

	// Wait before checking for changes.
	log.Info().Msgf("Waiting %d minutes to let autoscaler start the scale down", scaleDownTimeout)
	time.Sleep(scaleDownTimeout * time.Minute)
	// Check if build has been started, if not, error (Scale down).
	if res, err := c.GetConfigFromDB(context.Background(), &pb.GetConfigFromDBRequest{Id: config.Id, Type: pb.IdType_HASH}); err == nil {
		if checksumsEqual(res.Config.DsChecksum, res.Config.CsChecksum) {
			return fmt.Errorf("some cluster/s in config %s have not been scaled down, when they should have", config.Name)
		} else {
			log.Info().Msgf("Config %s has successfully passed autoscaling test [3/3]", config.Name)
		}
	} else {
		return fmt.Errorf("error while retrieving config %s from DB : %w", config.Name, err)
	}
	// Wait until build is finished.
	if err := configChecker(ctx, c, "autoscaling", "scale-down-test", idInfo{id: config.Id, idType: pb.IdType_HASH}); err != nil {
		return err
	}

	// Test longhorn.
	// Get new config from DB with updated counts.
	if res, err := c.GetConfigFromDB(context.Background(), &pb.GetConfigFromDBRequest{Id: config.Id, Type: pb.IdType_HASH}); err != nil {
		return err
	} else {
		if err := testLonghornDeployment(ctx, res.Config); err != nil {
			return err
		}
	}
	return nil
}

// applyDeployment applies specified deployment into specified cluster.
func applyDeployment(c *spec.K8Scluster, deployment string) error {
	kc := kubectl.Kubectl{Kubeconfig: c.Kubeconfig, MaxKubectlRetries: 5}
	if log.Logger.GetLevel() <= zerolog.InfoLevel {
		prefix := utils.GetClusterID(c.ClusterInfo)
		kc.Stdout = comm.GetStdOut(prefix)
		kc.Stderr = comm.GetStdErr(prefix)
	}
	if err := kc.KubectlApplyString(deployment, "-n", "default"); err != nil {
		return fmt.Errorf("failed to apply deployment on cluster %s : %w", c.ClusterInfo.Name, err)
	}
	return nil
}

// removeDeployment deletes specified deployment from specified cluster.
func removeDeployment(c *spec.K8Scluster, deployment string) error {
	kc := kubectl.Kubectl{Kubeconfig: c.Kubeconfig, MaxKubectlRetries: 5}
	if log.Logger.GetLevel() <= zerolog.InfoLevel {
		prefix := utils.GetClusterID(c.ClusterInfo)
		kc.Stdout = comm.GetStdOut(prefix)
		kc.Stderr = comm.GetStdErr(prefix)
	}
	if err := kc.KubectlDeleteString(deployment, "-n", "default"); err != nil {
		return fmt.Errorf("failed to remove deployment on cluster %s : %w", c.ClusterInfo.Name, err)
	}
	return nil
}
