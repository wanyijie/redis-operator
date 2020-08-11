/*


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

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv2 "github.com/wangyijie/redis-operator/api/v2"
)

// RedisReconciler reconciles a Redis object
type RedisReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=apps.wise2c.com,resources=redis,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.wise2c.com,resources=redis/status,verbs=get;update;patch

func (r *RedisReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("redis", req.NamespacedName)

	var redis appsv2.Redis
	if err := r.Get(ctx, req.NamespacedName, &redis); err != nil {
		log.Error(err, "unable to fetch redis")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	dep := tools.deploymentForRedis(&redis)
	fmt.Println("create deployment:", dep)
	err := r.Client.Create(context.TODO(), dep)
	if err != nil {
		log.Error(err, "create deployment failed")
		return ctrl.Result{}, err
	}

	//
	return ctrl.Result{}, nil
}

func (r *RedisReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv2.Redis{}).
		Complete(r)
}
