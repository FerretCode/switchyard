package autoscale

import (
	"time"

	"github.com/ferretcode/switchyard/autoscale/internal/railway/gql"
	"github.com/ferretcode/switchyard/autoscale/pkg/types"
)

func (a *AutoscaleService) getCurrentReplicas(project *gql.ProjectData, serviceId string) int {
	if project == nil {
		return 0
	}

	for _, service := range project.Project.Services.Edges {
		if service.Node.Id != serviceId {
			continue
		}

		for _, serviceInstance := range service.Node.ServiceInstances.Edges {
			if serviceInstance.Node.EnvironmentId != a.Config.RailwayEnvironmentId {
				continue
			}

			return serviceInstance.Node.LatestDeployment.Meta.ServiceManifest.Deploy.MultiRegionConfig[a.Config.RailwaySelectedRegion].NumReplicas
		}
	}
	return 0
}

func (a *AutoscaleService) makeScalingDecision(ctx types.ScalingContext) (int, string, error) {
	upscaleCooldown, err := time.ParseDuration(ctx.Service.UpscaleCooldown)
	if err != nil {
		return 0, "invalid-upscale-cooldown", err
	}
	downscaleCooldown, err := time.ParseDuration(ctx.Service.DownscaleCooldown)
	if err != nil {
		return 0, "invalid-downscale-cooldown", err
	}

	// emergency high usage
	if (ctx.CpuPercent > 0.9 || ctx.MemPercent > 0.9) && ctx.Now.Sub(lastUpscaleTime) > upscaleCooldown/3 {
		lastUpscaleTime = ctx.Now
		consecutiveLowLoad = 0
		return 2, "emergency-high-usage", nil
	}

	// spike detection
	if (ctx.HasCpuSpike || ctx.HasMemSpike) && ctx.Now.Sub(lastUpscaleTime) > upscaleCooldown/2 {
		lastUpscaleTime = ctx.Now
		consecutiveLowLoad = 0
		return 1, "spike-detected", nil
	}

	// predictive scaling
	isSustainedHighLoad := ctx.AvgCpu > ctx.Service.RailwayCpuUpscaleThreshold || ctx.AvgMem > ctx.Service.RailwayMemoryUpscaleThreshold
	isIncreasingTrend := ctx.CpuTrend > 0.01 || ctx.MemTrend > 0.01
	isNotInLowLoadZone := ctx.AvgCpu > ctx.Service.RailwayCpuDownscaleThreshold || ctx.AvgMem > ctx.Service.RailwayMemoryDownscaleThreshold

	if (isSustainedHighLoad || (isIncreasingTrend && isNotInLowLoadZone)) && ctx.Now.Sub(lastUpscaleTime) > upscaleCooldown {
		lastUpscaleTime = ctx.Now
		consecutiveLowLoad = 0
		return 1, "proactive-upscale", nil
	}

	// downscaling
	isLowLoad := ctx.AvgCpu < ctx.Service.RailwayCpuDownscaleThreshold && ctx.AvgMem < ctx.Service.RailwayMemoryDownscaleThreshold

	if isLowLoad {
		consecutiveLowLoad++
	} else {
		consecutiveLowLoad = 0
	}

	if consecutiveLowLoad >= 4 &&
		ctx.CurrentReplicas > int(ctx.Service.MinReplicaCount) &&
		ctx.Now.Sub(lastDownscaleTime) > downscaleCooldown {

		lastDownscaleTime = ctx.Now
		consecutiveLowLoad = 0
		return -1, "sustained-low-usage", nil
	}

	// no scaling required
	return 0, "no-scaling", nil
}

func extractMetrics(metrics *gql.MetricsData) (float64, float64) {
	if metrics == nil {
		return 0, 0
	}

	var memUsage, cpuUsage, memLimit, cpuLimit float64

	for _, metric := range metrics.Metrics {
		if len(metric.Values) == 0 {
			continue
		}

		switch metric.Measurement {
		case "MEMORY_USAGE_GB":
			memUsage = metric.Values[0].Value
		case "CPU_USAGE":
			cpuUsage = metric.Values[0].Value
		case "MEMORY_LIMIT_GB":
			memLimit = metric.Values[0].Value
		case "CPU_LIMIT":
			cpuLimit = metric.Values[0].Value
		}
	}

	cpuPercent := 0.0
	memPercent := 0.0

	if cpuLimit > 0 {
		cpuPercent = cpuUsage / cpuLimit
	}
	if memLimit > 0 {
		memPercent = memUsage / memLimit
	}

	return cpuPercent, memPercent
}

func calculateWeightedAverage(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	if len(values) == 1 {
		return values[0]
	}

	var totalValue float64
	var totalWeight float64

	for i, val := range values {
		weight := float64(i + 1)
		totalValue += val * weight
		totalWeight += weight
	}

	if totalWeight == 0 {
		return 0
	}

	return totalValue / totalWeight
}

func detectSpike(values []float64, threshold float64) bool {
	if len(values) < spikeWindow {
		return false
	}

	recentValues := values[len(values)-spikeWindow:]
	spikeCount := 0

	for _, val := range recentValues {
		if val >= threshold {
			spikeCount++
		}
	}

	return spikeCount >= (spikeWindow+1)/2
}

func calculateTrend(values []float64, times []time.Time) float64 {
	if len(values) < 3 || len(values) != len(times) {
		return 0
	}

	n := len(values)
	start := 0
	if n > 6 {
		start = max(n-15, 0)
	}

	sumX, sumY, sumXY, sumX2 := 0.0, 0.0, 0.0, 0.0
	startTime := times[start]

	for i := start; i < n; i++ {
		x := times[i].Sub(startTime).Seconds()
		y := values[i]
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	count := float64(n - start)
	if count < 2 {
		return 0
	}

	slopePerSecond := (count*sumXY - sumX*sumY) / (count*sumX2 - sumX*sumX)

	return slopePerSecond * 60
}
