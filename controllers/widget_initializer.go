/*
Copyright 2022.

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

package controllers

import (
	"context"
	"fmt"

	kcpclient "github.com/kcp-dev/apimachinery/pkg/client"
	apisv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/apis/v1alpha1"
	tenancyv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/tenancy/v1alpha1"
	"github.com/kcp-dev/logicalcluster"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// WidgetInitializer reconciles a Widget object
type WidgetInitializer struct {
	client.Client
	Scheme        *runtime.Scheme
	APIExportPath string
}

// +kubebuilder:rbac:groups=data.my.domain,resources=widgets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=data.my.domain,resources=widgets/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=data.my.domain,resources=widgets/finalizers,verbs=update

// Reconcile TODO
func (r *WidgetInitializer) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger = logger.WithValues("clusterName", req.ClusterName)
	logger.Info("reconciling ClusterWorkspace of the widget type")

	clusterName := fmt.Sprintf("%s:%s", req.ClusterName, req.Name)
	logger.Info("the expected clusterName", "clusterName", clusterName)

	ctx = kcpclient.WithCluster(ctx, logicalcluster.New(clusterName))

	apiBinding := &apisv1alpha1.APIBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "workload",
			ClusterName: clusterName,
		},
		Spec: apisv1alpha1.APIBindingSpec{
			Reference: apisv1alpha1.ExportReference{
				Workspace: &apisv1alpha1.WorkspaceExportReference{
					ExportName: "kubernetes",
					Path:       r.APIExportPath,
				},
			},
		},
	}

	return ctrl.Result{}, r.Client.Create(ctx, apiBinding)
}

// SetupWithManager sets up the controller with the Manager.
func (r *WidgetInitializer) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&tenancyv1alpha1.ClusterWorkspace{}).
		Complete(r)
}
