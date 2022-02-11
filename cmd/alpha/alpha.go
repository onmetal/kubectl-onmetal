// Copyright 2021 OnMetal authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package alpha

import (
	"context"
	"errors"

	"github.com/onmetal/kubectl-onmetal/alpha"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Command(restClientGetter genericclioptions.RESTClientGetter) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "alpha tree",
		Args: cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			return Run(ctx, restClientGetter, args)
		},
	}

	return cmd
}

func Run(ctx context.Context, restClientGetter genericclioptions.RESTClientGetter, args []string) error {
	if args != nil && len(args) == 0 {
		return errors.New("missing argument")
	}

	cfg, err := restClientGetter.ToRESTConfig()
	if err != nil {
		return errors.New("error getting rest config: " + err.Error())
	}

	namespace, _, err := restClientGetter.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return errors.New("error determining target namespace: " + err.Error())
	}

	c, err := client.New(cfg, client.Options{})
	if err != nil {
		return errors.New("error creating client: " + err.Error())
	}

	return alpha.Run(ctx, c, namespace, args)
}
