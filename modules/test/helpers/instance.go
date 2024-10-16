/*
Copyright 2023 Red Hat
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

package helpers

import (
	"github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"

	k8s_errors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

// DeleteInstance deletes a specified resource and waits until it's fully removed from the cluster.
//
// Example usage:
//
//	DeferCleanup(th.DeleteInstance, metadata_instance)
func (tc *TestHelper) DeleteInstance(instance client.Object, opts ...client.DeleteOption) {
	// We have to wait for the controller to fully delete the instance
	tc.Logger.Info(
		"Deleting", "Name", instance.GetName(),
		"Namespace", instance.GetNamespace(),
		"Kind", instance.GetObjectKind(),
	)

	gomega.Eventually(func(g gomega.Gomega) {
		name := types.NamespacedName{Name: instance.GetName(), Namespace: instance.GetNamespace()}
		err := tc.K8sClient.Get(tc.Ctx, name, instance)
		// if it is already gone that is OK
		if k8s_errors.IsNotFound(err) {
			return
		}
		g.Expect(err).ShouldNot(gomega.HaveOccurred())

		g.Expect(tc.K8sClient.Delete(tc.Ctx, instance, opts...)).Should(gomega.Succeed())

		err = tc.K8sClient.Get(tc.Ctx, name, instance)
		g.Expect(k8s_errors.IsNotFound(err)).To(gomega.BeTrue())
	}, tc.Timeout, tc.Interval).Should(gomega.Succeed())

	tc.Logger.Info(
		"Deleted", "Name", instance.GetName(),
		"Namespace", instance.GetNamespace(),
		"Kind", instance.GetObjectKind().GroupVersionKind().Kind,
	)
}
