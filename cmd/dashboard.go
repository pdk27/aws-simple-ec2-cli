// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package cmd

import (
	"fmt"
	"simple-ec2/pkg/cli"
	"simple-ec2/pkg/ec2dashboardhelper"
	"simple-ec2/pkg/ec2dashboardhelper/config"
	"simple-ec2/pkg/ec2helper"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/spf13/cobra"
)

// connectCmd represents the connect command
var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "dashboard to an AWS Resources",
	Long:  `Dashboard allows you to view summary of all your AWS resources and recommendations for cost-saving`,
	Run:   dashboard,
}

var	cfg config.Config

// Add flags
func init() {
	rootCmd.AddCommand(dashboardCmd)

	dashboardCmd.Flags().StringVarP(&regionFlag, "region", "r", "", "The region for which you would like to see AWS resources and recommendations")
	dashboardCmd.Flags().StringVarP(&granularityFlag, "granularity", "g", "DAILY", "The AWS cost granularity. Choose from [MONTHLY, DAILY, HOURLY] (Default: DAILY)")
	dashboardCmd.Flags().StringVarP(&costTypeFlag, "costType", "c", "BlendedCost,UnblendedCost",
		"The type of costs. Choose from [AmortizedCost, BlendedCost, NetAmortizedCost, NetUnblendedCost, UnblendedCost] (Default: 'BlendedCost','UnblendedCost'")
	dashboardCmd.Flags().IntVarP(&evalPeriodInDaysFlag, "evaluationPeriodInDays", "p", 7, "The evaluation period for costs and metrics in days. (Default: 7)")
}

// The main function
func dashboard(cmd *cobra.Command, args []string) {
	if !ValidateDashboardFlags() {
		return
	}

	// Start a new session, with the default credentials and config loading
	sess := session.Must(session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable}))
	config := config.Config{
		AWSSession: sess,
		Region: regionFlag,
		Granularity: granularityFlag,
		CostType: costTypeFlag,
		EvaluationPeriodInDays: evalPeriodInDaysFlag,
	}
	ec2helper.GetDefaultRegion(sess)

	printFlags(cmd.Flags())
	notes := "NOTE: " +
		"\n\t* All costs are in USD" +
		"\n\t* Averages and max values showed below are calculated over the configured evaluation period."
	fmt.Println(notes)

	dashboardSummary(config)
}

// Fetch dashboard summary for AWS resources
func dashboardSummary(config config.Config) {
	h := ec2helper.New(config.AWSSession)

	var err error
	// Override region if specified
	if regionFlag != "" {
		h.ChangeRegion(regionFlag)
		//err = GetDashboardSummaryForRegion(h)
		ec2dashboardhelper.GenerateDashboardForRegionWithEverything(config)
	} else {
		err = GetDashboardSummaryWorldWide(h)
	}

	if cli.ShowError(err, "Generating dashboard failed") {
		return
	}
}

// Validate flags using simple rules. Return true if the flags are validated, false otherwise
func ValidateDashboardFlags() bool {
	// TODO: add any flag validations here
	return true
}

// Get the information of the instance and connect to it
func GetDashboardSummaryForRegion(h *ec2helper.EC2Helper) error {
	ec2dashboardhelper.GenerateDashboardForRegion(h)
	return nil
}

// Get the information of the instance and connect to it
func GetDashboardSummaryWorldWide(h *ec2helper.EC2Helper) error {
	err := ec2dashboardhelper.GenerateDashboardWorldWide(h)
	if err != nil {
		return err
	}

	return nil
}