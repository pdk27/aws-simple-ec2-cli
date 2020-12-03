package ec2dashboardhelper

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/ec2"
	"simple-ec2/pkg/ec2dashboardhelper/computeOptimizer"
	"simple-ec2/pkg/ec2dashboardhelper/config"
	"simple-ec2/pkg/ec2dashboardhelper/costTracker"
	"simple-ec2/pkg/ec2dashboardhelper/info"
	"simple-ec2/pkg/ec2helper"
	"simple-ec2/pkg/table"
)

// Generate dashboard for the region
func GenerateDashboardForRegion(h *ec2helper.EC2Helper) {
	// TODO: Generate dashboard here
	// Only include running states
	states := []string{
		ec2.InstanceStateNameRunning,
	}

	instances, _ := h.GetInstancesByState(states)
	if instances == nil {
		return
	}
	fmt.Println("\n" + *h.Sess.Config.Region)

	var data [][]string
	var indexedOptions []string

	data, indexedOptions = table.AppendDataForDashboard(data, indexedOptions, instances)

	optionsText := table.BuildTable(data, []string{"Instance", "Recommendations", "Idle Time", "Cost Saving"})
	fmt.Print(optionsText)
	return
}

// Generate dashboard for all regions
func GenerateDashboardWorldWide(h *ec2helper.EC2Helper) error {
	// TODO: Generate dashboard here by calling for `GetDashboardSummaryForRegion` for each region
	regions , _ := h.GetEnabledRegions()
	for _, region := range regions {
		h.ChangeRegion(*region.RegionName)
		GenerateDashboardForRegion(h)
	}
	return nil
}


func GenerateDashboardForRegionWithEverything(cfg config.Config) {
	var instancesInfo []info.InstanceInfo
	instancesInfo = computeOptimizer.GetRecommendation(cfg, instancesInfo)
	//instancesInfo = ec2helper.PopulateInstanceInfo(cfg, instancesInfo) // populate id, type, other information
	instancesInfo = costTracker.PopulateCosts(cfg, instancesInfo)
	//instancesInfo = recommendations.PopulateMetrics(cfg, instancesInfo)
	//instancesInfo = cwMetrics.PopulateMetrics(cfg, instancesInfo)

	info.PrintTable(instancesInfo)
}

// Contains finds a string in the given array
func Contains(slice []*string, val *string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}