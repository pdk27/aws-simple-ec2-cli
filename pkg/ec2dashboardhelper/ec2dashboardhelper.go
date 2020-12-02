package ec2dashboardhelper

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/ec2"
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
