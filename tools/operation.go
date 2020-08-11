package tools

import (
	"fmt"

	appsv2 "github.com/wangyijie/redis-operator/api/v2"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func deploymentForRedis(app *appsv2.Redis) *appsv1.Deployment {
	ls := labelsForRedis("redis")
	replicas := app.Spec.Size
	fmt.Println("size:", replicas)

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.ObjectMeta.Name,
			Namespace: app.ObjectMeta.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: "redis:alpine",
						Name:  "redis",
						Ports: []corev1.ContainerPort{{
							ContainerPort: 6379,
							Name:          "redis",
						}},
					}},
				},
			},
		},
	}
	// Set Memcached instance as the owner and controller
	return dep
}

func labelsForRedis(name string) map[string]string {
	return map[string]string{"app": "memcached", "memcached_cr": name}
}

// func configMapForRedis(app *appsv2.Redis) *corev1.ConfigMap {
// 	ls := labelsForRedis("redis")
// 	cm := &corev1.ConfigMap{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      app.ObjectMeta.Name,
// 			Namespace: app.ObjectMeta.Namespace,
// 			Labels:    ls,
// 		},
// 		Data: ,
// 	}

// }
