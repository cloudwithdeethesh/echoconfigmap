package controller

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	demov1alpha1 "github.com/cloudwithdeethesh/enchoconfigmap/api/v1alpha1"
)

// RBAC needed by the controller
// +kubebuilder:rbac:groups=demo.deet.dev,resources=echoconfigs,verbs=get;list;watch
// +kubebuilder:rbac:groups=demo.deet.dev,resources=echoconfigs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch

type EchoConfigReconciler struct {
	client.Client
	Scheme *runtime.Scheme // <- add this so main.go/tests can set it
}

func (r *EchoConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// 1) Fetch the EchoConfig instance
	var ec demov1alpha1.EchoConfig
	if err := r.Get(ctx, req.NamespacedName, &ec); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 2) Desired ConfigMap name
	cmName := fmt.Sprintf("echo-%s", ec.Name)

	// 3) Ensure the ConfigMap exists and has the desired message
	var cm corev1.ConfigMap
	err := r.Get(ctx, types.NamespacedName{Name: cmName, Namespace: ec.Namespace}, &cm)
	if apierrors.IsNotFound(err) {
		cm = corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{ // <- use metav1 here
				Name:      cmName,
				Namespace: ec.Namespace,
			},
			Data: map[string]string{"message": ec.Spec.Message},
		}
		if err := r.Create(ctx, &cm); err != nil {
			return ctrl.Result{}, err
		}
		logger.Info("created ConfigMap", "configMap", cmName)
	} else if err == nil {
		if cm.Data == nil {
			cm.Data = map[string]string{}
		}
		if cm.Data["message"] != ec.Spec.Message {
			cm.Data["message"] = ec.Spec.Message
			if err := r.Update(ctx, &cm); err != nil {
				return ctrl.Result{}, err
			}
			logger.Info("updated ConfigMap", "configMap", cmName)
		}
	} else {
		return ctrl.Result{}, err
	}

	// 4) Update status
	needStatus := false
	if ec.Status.ConfigMapName != cmName {
		ec.Status.ConfigMapName = cmName
		needStatus = true
	}
	if ec.Status.ObservedGeneration != ec.GetGeneration() {
		ec.Status.ObservedGeneration = ec.GetGeneration()
		needStatus = true
	}
	if needStatus {
		if err := r.Status().Update(ctx, &ec); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *EchoConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&demov1alpha1.EchoConfig{}).
		Complete(r)
}
