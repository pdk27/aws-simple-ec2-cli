package ec2dashboardhelper

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
)


// Generate dashboard for the region
func GenerateDashboardForRegion(sess *session.Session) error {
	// TODO: Generate dashboard here
	fmt.Printf("Regional Dashboard") // Prints `Binary: 100\101`

	return nil
}

// Generate dashboard for all regions
func GenerateDashboardWorldWide(sess *session.Session) error {
	// TODO: Generate dashboard here by calling for `GetDashboardSummaryForRegion` for each region
	fmt.Printf("World-wide Dashboard") // Prints `Binary: 100\101`

	return nil
}