package scheduler

import (
	"log"
	"time"

	"github.com/rakutentech/shibuya/shibuya/config"
	"github.com/rakutentech/shibuya/shibuya/model"
	smodel "github.com/rakutentech/shibuya/shibuya/scheduler/model"
	apiv1 "k8s.io/api/core/v1"
)

type EngineScheduler interface {
	DeployEngine(projectID, collectionID, planID int64, engineID int, containerConfig *config.ExecutorContainer) error
	CollectionStatus(projectID, collectionID int64, eps []*model.ExecutionPlan) (*smodel.CollectionStatus, error)
	FetchEngineUrlsByPlan(collectionID, planID int64, opts *smodel.EngineOwnerRef) ([]string, error)
	ExposeCollection(ProjectID, collectionID int64) error
	PurgeCollection(collectionID int64) error
	GetDeployedCollections() (map[int64]time.Time, error)
	GetAllNodesInfo() (smodel.AllNodesInfo, error)
	GetPodsMetrics(collectionID, planID int64) (map[string]apiv1.ResourceList, error)
	PodReadyCount(collectionID int64) int
	DownloadPodLog(collectionID, planID int64) (string, error)
}

func NewEngineScheduler(cfg *config.SchedulerConfig) EngineScheduler {
	switch cfg.Kind {
	case "k8s":
		return NewK8sClientManager(cfg)
	case "cloudrun":
		return NewCloudRun(cfg)
	}
	log.Fatalf("Shibuya does not support %s as scheduler", cfg.Kind)
	return nil
}
