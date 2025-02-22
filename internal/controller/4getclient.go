package controller

import (
	"context"
	gateway "gateway/api/v1beta1"
	"gateway/internal/controller/ontap"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
)

func (r *StorageVirtualMachineReconciler) reconcileGetClient(ctx context.Context,
	svmCR *gateway.StorageVirtualMachine,
	adminSecret *corev1.Secret, host string, trustSSL bool,
	log logr.Logger) (*ontap.Client, error) {

	log.Info("STEP 4: Create ONTAP client")

	oc, err := ontap.NewClient(
		string(adminSecret.Data["username"]),
		string(adminSecret.Data["password"]),
		host, svmCR.Spec.SvmDebug, trustSSL)

	if err != nil {
		log.Error(err, "Error creating ONTAP client - requeueing")
		_ = r.setConditionONTAPCreation(ctx, svmCR, CONDITION_STATUS_FALSE)
		return oc, err
	}

	log.Info("ONTAP client created")
	_ = r.setConditionONTAPCreation(ctx, svmCR, CONDITION_STATUS_TRUE)

	return oc, nil

}
