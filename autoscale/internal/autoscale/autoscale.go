package autoscale

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ferretcode/switchyard/autoscale/internal/railway"
	"github.com/ferretcode/switchyard/autoscale/internal/railway/gql"
	"github.com/ferretcode/switchyard/autoscale/internal/repositories"
	"github.com/ferretcode/switchyard/autoscale/pkg/types"
)

const (
	spikeThreshold = 0.7
	spikeWindow    = 3
)

var (
	lastUpscaleTime     time.Time
	lastDownscaleTime   time.Time
	consecutiveHighLoad int
	consecutiveLowLoad  int
)

type AutoscaleService struct {
	Logger              *slog.Logger
	Config              *types.Config
	GqlQueries          *railway.QueryService
	Queries             *repositories.Queries
	Context             context.Context
	ServiceHistoryCache *types.ServiceHistoryCache
}

func NewAutoscaleService(logger *slog.Logger, config *types.Config, gqlQueries *railway.QueryService, queries *repositories.Queries, context context.Context, serviceHistoryCache *types.ServiceHistoryCache) AutoscaleService {
	return AutoscaleService{
		Logger:              logger,
		Config:              config,
		GqlQueries:          gqlQueries,
		Queries:             queries,
		Context:             context,
		ServiceHistoryCache: serviceHistoryCache,
	}
}

func (a *AutoscaleService) StartAutoscaling() {
	monitoringTimer := time.NewTicker(a.Config.MonitoringInterval)
	defer monitoringTimer.Stop()

	for {
		<-monitoringTimer.C

		registeredServices, err := a.Queries.ListServices(a.Context)
		if err != nil {
			a.Logger.Error("error fetching registered services", "err", err)
			continue
		}

		project, err := a.GqlQueries.QueryProjectData()
		if err != nil {
			a.Logger.Error("error fetching project state", "err", err)
		}

		validServices := getValidServiceIds(registeredServices, project.Project.Services.Edges)

		for _, serviceId := range validServices {
			a.processServiceId(serviceId, project)
		}

	}
}

func (a *AutoscaleService) processServiceId(validService ValidService, project *gql.ProjectData) {
	startDate := time.Now().Format(time.RFC3339)

	history, ok := a.ServiceHistoryCache.ServiceHistories[validService.ServiceId]
	if !ok {
		history = types.MetricHistory{
			CPU:    make([]float64, 0, a.Config.MetricHistorySize),
			Memory: make([]float64, 0, a.Config.MetricHistorySize),
			Times:  make([]time.Time, 0, a.Config.MetricHistorySize),
		}
		a.ServiceHistoryCache.ServiceHistories[validService.ServiceId] = history
	}

	metrics, err := a.GqlQueries.QueryServiceMetrics(validService.ServiceId, startDate)
	if err != nil {
		a.Logger.Error("error fetching service metrics", "err", err)
		return
	}

	currentReplicas := a.getCurrentReplicas(project, validService.ServiceId)

	cpuPercent, memPercent := extractMetrics(metrics)

	now := time.Now()
	history.CPU = append(history.CPU, cpuPercent)
	history.Memory = append(history.Memory, memPercent)
	history.Times = append(history.Times, now)

	if len(history.CPU) > a.Config.MetricHistorySize {
		history.CPU = history.CPU[1:]
		history.Memory = history.Memory[1:]
		history.Times = history.Times[1:]
	}

	a.ServiceHistoryCache.ServiceHistories[validService.ServiceId] = history

	avgCpu := calculateWeightedAverage(history.CPU)
	avgMem := calculateWeightedAverage(history.Memory)

	hasCpuSpike := detectSpike(history.CPU, spikeThreshold)
	hasMemSpike := detectSpike(history.Memory, spikeThreshold)

	cpuTrend := calculateTrend(history.CPU, history.Times)
	memTrend := calculateTrend(history.Memory, history.Times)

	scalingDecision, reason, err := a.makeScalingDecision(
		types.ScalingContext{
			CpuPercent:      cpuPercent,
			MemPercent:      memPercent,
			AvgCpu:          avgCpu,
			AvgMem:          avgMem,
			HasCpuSpike:     hasCpuSpike,
			HasMemSpike:     hasMemSpike,
			CpuTrend:        cpuTrend,
			MemTrend:        memTrend,
			CurrentReplicas: currentReplicas,
			Now:             now,
			Service:         validService.Service,
		},
	)
	if err != nil {
		a.Logger.Error("error making scaling decision", "err", err)
		return
	}

	if scalingDecision != 0 {
		newReplicas := currentReplicas + scalingDecision
		if newReplicas >= int(validService.Service.MinReplicaCount) && newReplicas <= int(validService.Service.MaxReplicaCount) {
			a.Logger.Info("scaling decision reached",
				"current_replicas", currentReplicas,
				"new_replicas", newReplicas,
				"reason", reason,
			)

			err = a.GqlQueries.MutationUpdateReplicas(a.Config.RailwayEnvironmentId, validService.ServiceId, a.Config.RailwaySelectedRegion, newReplicas)
			if err != nil {
				a.Logger.Error("error scaling service", "err", err)
			}

			err = a.GqlQueries.MutationServiceInstanceRedeploy(a.Config.RailwayEnvironmentId, validService.ServiceId)
			if err != nil {
				a.Logger.Error("error redeploying scaled service", "err", err)
			}
		}
	}

	a.Logger.Info("monitoring metrics",
		"current-cpu", fmt.Sprintf("%.2f%%", cpuPercent*100),
		"current-mem", fmt.Sprintf("%.2f%%", memPercent*100),
		"avg-cpu", fmt.Sprintf("%.2f%%", avgCpu*100),
		"avg-mem", fmt.Sprintf("%.2f%%", avgMem*100),
		"cpu-trend", fmt.Sprintf("%.4f", cpuTrend),
		"mem-trend", fmt.Sprintf("%.4f", memTrend),
		"cpu-spike", hasCpuSpike,
		"mem-spike", hasMemSpike,
		"replicas", currentReplicas,
		"consecutive-high", consecutiveHighLoad,
		"consecutive-low", consecutiveLowLoad,
	)
}

func getValidServiceIds(registeredServices []repositories.Service, projectServiceNodes []gql.ServiceEdge) []ValidService {
	var validServiceIds []ValidService

	for _, registeredService := range registeredServices {
		for _, service := range projectServiceNodes {
			if registeredService.ServiceID == service.Node.Id && registeredService.Enabled {
				validServiceIds = append(validServiceIds, ValidService{
					ServiceId: service.Node.Id,
					Service:   registeredService,
				})
			}
		}
	}

	return validServiceIds
}
