package custom_webhook

import (
    "context"

    openshiftquotav1 "github.com/openshift/api/quota/v1"
    corev1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/types"
    "sigs.k8s.io/controller-runtime/pkg/client"
    "sigs.k8s.io/controller-runtime/pkg/log"
    "sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

//+kubebuilder:webhook:path=/validate-v1-resource-quota,mutating=false,failurePolicy=fail,sideEffects=None,groups="",resources=resourcequotas,verbs=create;update;delete,versions=v1,name=vresourcequota.kb.io,admissionReviewVersions={v1,v1beta1}

type ResourceQuotaValidator struct {
    Client client.Client
}

const (
    teamLabel = "snappcloud.io/team"
    enforceLabel = "quota.snappcloud.io/enforce"
)

func (v *ResourceQuotaValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
    log := log.FromContext(ctx)
    if req.Operation == "UPDATE" {
        ns := &corev1.Namespace{}
        err := v.Client.Get(context.TODO(), types.NamespacedName{Name: req.Namespace}, ns)
        if err != nil {
            log.Error(err, "error getting namespace", "name", req.Namespace)
            return admission.Denied("error on getting namespace")
        }
        l, ok := ns.GetLabels()[teamLabel]
        if !ok {
            return admission.Denied("no team found for the project. please join your project to a team")
        }
        if l,ok := ns.GetLabels()[enforceLabel]; ok{
            if l == "false" || l == ""{
                return admission.Allowed("ignoring resourcequota")
            }
        }else{
            return admission.Denied("no enforce label found for the project")
        }
        crq := &openshiftquotav1.ClusterResourceQuota{}
        err = v.Client.Get(context.TODO(), types.NamespacedName{Name: l}, crq)
        if err != nil {
            log.Error(err, "error getting clusterResourceQuota", "name", l)
            return admission.Denied("no team quota found. please request a quota for your team in cloud-support")
        }
        return admission.Allowed("updating resourcequota")
    } else if req.Operation == "DELETE" {
        if req.Name == "default" {
            return admission.Denied("default resourcequota cannot be deleted")
        }
        return admission.Allowed("DELETE")
    } else {
        return admission.Allowed("Allowed")
    }
}