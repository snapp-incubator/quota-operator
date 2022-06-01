package v1alpha1

import (
	"context"
	"fmt"
	"net/http"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:webhook:path=/validate-v1-resource-quota,mutating=false,failurePolicy=fail,groups="",resources=resourcequotas,verbs=create;update,versions=v1,name=vresourcequota.kb.io

// podValidator validates Pods
type resourceQuotaValidator struct {
	Client  client.Client
	decoder *admission.Decoder
}

// podValidator admits a pod if a specific annotation exists.
func (v *resourceQuotaValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	resourcequota := &corev1.ResourceQuota{}

	err := v.decoder.Decode(req, resourcequota)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	key := "example-mutating-admission-webhook"
	anno, found := resourcequota.Annotations[key]
	if !found {
		return admission.Denied(fmt.Sprintf("missing annotation %s", key))
	}
	if anno != "foo" {
		return admission.Denied(fmt.Sprintf("annotation %s did not have value %q", key, "foo"))
	}

	return admission.Allowed("")
}

// podValidator implements admission.DecoderInjector.
// A decoder will be automatically injected.

// InjectDecoder injects the decoder.
func (v *resourceQuotaValidator) InjectDecoder(d *admission.Decoder) error {
	v.decoder = d
	return nil
}
