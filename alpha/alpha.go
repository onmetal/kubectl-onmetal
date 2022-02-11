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

	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	computev1alpha1 "github.com/onmetal/onmetal-api/apis/compute/v1alpha1"
	storagev1alpha1 "github.com/onmetal/onmetal-api/apis/storage/v1alpha1"
)

func Run(ctx context.Context, c client.Client, namespace string, args []string) error {
	log := controllerruntime.LoggerFrom(ctx)
	log.Info("")

	var tree = map[string][]string{}

	switch args[0] {
	case "machinepool":
		if err := machineTree(ctx, c, tree, namespace); err != nil {
			return err
		}
	case "storagepool":
		if err := storageTree(ctx, c, tree, namespace); err != nil {
			return err
		}
	}

	for pool, items := range tree {
		println("-", pool)
		for _, item := range items {
			println(" ", "-", item)
		}
	}

	return nil
}

func machineTree(ctx context.Context, c client.Client, tree map[string][]string, namespace string) error {
	var ml computev1alpha1.MachineList
	var mpl computev1alpha1.MachinePoolList

	if err := c.List(ctx, &ml, client.InNamespace(namespace)); err != nil {
		return err
	}
	if err := c.List(ctx, &mpl, client.InNamespace(namespace)); err != nil {
		return err
	}

	for _, mp := range mpl.Items {
		tree[mp.Name] = []string{}
	}
	for _, m := range ml.Items {
		if _, ok := tree[m.Spec.MachinePool.Name]; ok {
			tree[m.Spec.MachinePool.Name] = append(tree[m.Spec.MachinePool.Name], m.Name)
		} else {
			tree["unPooled"] = append(tree["unPooled"], m.Name)
		}
	}

	return nil
}

func storageTree(ctx context.Context, c client.Client, tree map[string][]string, namespace string) error {
	var vl storagev1alpha1.VolumeList
	var spl storagev1alpha1.StoragePoolList

	if err := c.List(ctx, &vl, client.InNamespace(namespace)); err != nil {
		return err
	}
	if err := c.List(ctx, &spl, client.InNamespace(namespace)); err != nil {
		return err
	}

	for _, sp := range spl.Items {
		tree[sp.Name] = []string{}
	}
	for _, v := range vl.Items {
		if _, ok := tree[v.Spec.StoragePool.Name]; ok {
			tree[v.Spec.StoragePool.Name] = append(tree[v.Spec.StoragePool.Name], v.Name)
		} else {
			tree["unPooled"] = append(tree["unPooled"], v.Name)
		}
	}

	return nil
}
