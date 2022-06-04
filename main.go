/*
Copyright 2021.

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

package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	corev1 "k8s.io/api/core/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	openshiftquotav1 "github.com/openshift/api/quota/v1"
	quotav1alpha1 "github.com/snapp-cab/quota-operator/api/v1alpha1"
	"github.com/snapp-cab/quota-operator/controllers"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(quotav1alpha1.AddToScheme(scheme))
	utilruntime.Must(openshiftquotav1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

type resourceQuotaValidator struct {
	Client  client.Client
	decoder *admission.Decoder
}

const (
	teamLabel        = "snappcloud.io/team"
	defaultQuotaName = "default"
)

func (v *resourceQuotaValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	ns := &corev1.Namespace{}
	resourcequota := &corev1.ResourceQuota{}

	l, ok := ns.GetLabels()[teamLabel]
	fmt.Printf("%s", l)
	if !ok {
		return admission.Denied("no team found for the project. please join your project to a team")
	}

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

// InjectDecoder injects the decoder.
func (v *resourceQuotaValidator) InjectDecoder(d *admission.Decoder) error {
	v.decoder = d
	return nil
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "bc6545ad.snappcloud.io",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&controllers.QuotaReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Quota")
		os.Exit(1)
	}
	if err = (&quotav1alpha1.Quota{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "Quota")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder
	hookServer := mgr.GetWebhookServer()

	hookServer.Register("/validate-v1-resource-quota", &webhook.Admission{Handler: &resourceQuotaValidator{Client: mgr.GetClient()}})

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
