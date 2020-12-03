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
	"github.com/spf13/pflag"
	"simple-ec2/pkg/config"
)

// Used for flags
var (
	instanceIdConnectFlag string
	isInteractive         bool
	isSaveConfig          bool
	regionFlag            string
	instanceIdFlag        []string
)


var (
	granularityFlag 	 string
	costTypeFlag         string
	evalPeriodInDaysFlag int
)

var flagConfig = config.SimpleInfo{}

// PrintFlags prints all flags of a command, if set
func printFlags(flags *pflag.FlagSet) {
	f := make(map[string]interface{})
	flags.Visit(func(flag *pflag.Flag) {
		f[flag.Name] = flag.Value
	})

	if len(f) > 0 {
		fmt.Println("\nFlags:")
		for key, value := range f {
			fmt.Printf("%s: %s\n", key, value)
		}
		fmt.Println()
	}
}
