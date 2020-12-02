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
	fmt.Println("Regional Dashboard")
	// Only include running states
	states := []string{
		ec2.InstanceStateNameRunning,
	}

	instances, _ := h.GetInstancesByState(states)
	var data [][]string
	var indexedOptions []string

	data, indexedOptions = table.AppendInstancesForDashboard(data, indexedOptions, instances)

	optionsText := table.BuildTable(data, []string{"Instance", "Recommendations", "Idle Time", "Cost Saving"})
	fmt.Print(optionsText)
	return
}

// Generate dashboard for all regions
func GenerateDashboardWorldWide(h *ec2helper.EC2Helper) error {
	// TODO: Generate dashboard here by calling for `GetDashboardSummaryForRegion` for each region
	fmt.Printf("World-wide Dashboard")

	return nil
}
