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

package delete

import (
	"github.com/spf13/cobra"
	"gitlab.cee.redhat.com/openshift-group-I/ocm_aws/cmd/delete/tag"
	"gitlab.cee.redhat.com/openshift-group-I/ocm_aws/cmd/delete/vpc"
)

var Cmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"del"},
	Short:   "Delete a resource from stdin",
	Long:    "Delete a resource from stdin",
}

func init() {
	Cmd.AddCommand(vpc.Cmd)
	Cmd.AddCommand(tag.Cmd)
}
