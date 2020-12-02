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
	"simple-ec2/pkg/ec2dashboardhelper"
	"simple-ec2/pkg/cli"
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

// Add flags
func init() {
	rootCmd.AddCommand(dashboardCmd)

	dashboardCmd.Flags().StringVarP(&regionFlag, "region", "r", "",
		"The region for which you would like to see AWS resources and recommendations")
}

// The main function
func dashboard(cmd *cobra.Command, args []string) {
	if !ValidateDashboardFlags() {
		return
	}

	// Start a new session, with the default credentials and config loading
	sess := session.Must(session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable}))
	ec2helper.GetDefaultRegion(sess)
	h := ec2helper.New(sess)

	dashboardSummary(h)
}

// Fetch dashboard summary for AWS resources
func dashboardSummary(h *ec2helper.EC2Helper) {
	var err error
	// Override region if specified
	if regionFlag != "" {
		h.ChangeRegion(regionFlag)
		err = GetDashboardSummaryForRegion(h)
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
	err := ec2dashboardhelper.GenerateDashboardForRegion(h.Sess)
	if err != nil {
		return err
	}

	return nil
}

// Get the information of the instance and connect to it
func GetDashboardSummaryWorldWide(h *ec2helper.EC2Helper) error {
	err := ec2dashboardhelper.GenerateDashboardWorldWide(h.Sess)
	if err != nil {
		return err
	}

	return nil
}