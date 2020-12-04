package ec2dashboardhelper

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/ec2"
	"simple-ec2/pkg/ec2dashboardhelper/config"
	"simple-ec2/pkg/ec2dashboardhelper/costTracker"
	"simple-ec2/pkg/ec2dashboardhelper/info"
	"simple-ec2/pkg/ec2helper"
	"simple-ec2/pkg/table"
)

// Generate dashboard for the region
func GenerateDashboardForRegion(h *ec2helper.EC2Helper, config config.Config) {
	// TODO: Generate dashboard here
	// Only include running states
	GenerateDashboardWithEverything(config)

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
func GenerateDashboardWorldWide(h *ec2helper.EC2Helper, conf config.Config) error {
	// TODO: Generate dashboard here by calling for `GetDashboardSummaryForRegion` for each region
	regions , _ := h.GetEnabledRegions()
	for _, region := range regions {
		// This throws error as it is internal region
		if *region.RegionName == "ap-northeast-3" {
			panic("Internal region can't be used. Please select a different region.")
		}

		//h.ChangeRegion(*region.RegionName)
		//c := config.Config{
		//	AWSSession: h.Sess,
		//	Region:     *region.RegionName,
		//}
		GenerateDashboardWithEverything(conf)
	}
	return nil
}

func GenerateDashboardWithEverything(cfg config.Config) {
	// create dashboard tracked by instance id
	result := make(map[string]info.InstanceInfo)
	//result = computeOptimizer.PopulateRecommendations(cfg)
	//info.PrintTable(result)

	instancesInfoCosts := costTracker.PopulateCostsAndType(cfg)
	result = info.Merge(instancesInfoCosts, result)
	info.PrintTable(result)

	//instancesInfoMetrics := cwMetrics.PopulateMetrics(cfg)
	//result = info.Merge(instancesInfoCosts, result)

	// print cost tracker tables by categories - these will be separate tables
	if cfg.ShowCostsByCategories {
		result = info.Merge(costTracker.PopulateRegion(cfg), result)
		info.PrintTable(result)
		result = info.Merge(costTracker.PopulateCapacityType(cfg), result)
		info.PrintTable(result)
	}

	//info.PrintTable(result)
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