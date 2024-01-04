/*
Copyright (c) 2019 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package create

import (
	"github.com/openshift-qe/openshift-rosa-cli/cmd/create/proxy"
	"github.com/openshift-qe/openshift-rosa-cli/cmd/create/sg"
	"github.com/openshift-qe/openshift-rosa-cli/cmd/create/subnets"
	"github.com/openshift-qe/openshift-rosa-cli/cmd/create/vpc"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"add"},
	Short:   "Create a resource from stdin",
	Long:    "Create a resource from stdin",
}

func init() {
	Cmd.AddCommand(vpc.Cmd)
	Cmd.AddCommand(sg.Cmd)
	Cmd.AddCommand(subnets.Cmd)
	Cmd.AddCommand(proxy.Cmd)

	// flags := Cmd.PersistentFlags()
	// arguments.AddProfileFlag(flags)
	// arguments.AddRegionFlag(flags)
	// confirm.AddFlag(flags)

	// globallyAvailableCommands := []*cobra.Command{
	// 	accountroles.Cmd, operatorroles.Cmd,
	// 	userrole.Cmd, ocmrole.Cmd,
	// 	oidcprovider.Cmd,
	// }
	// arguments.MarkRegionHidden(Cmd, globallyAvailableCommands)
}
