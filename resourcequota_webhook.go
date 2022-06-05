package main

//+kubebuilder:webhook:path=/validate-v1-resource-quota,mutating=false,failurePolicy=fail,sideEffects=None,groups="",resources=resourcequotas,verbs=create;update;delete,versions=v1,name=vresourcequota.kb.io,admissionReviewVersions={v1,v1beta1}
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=resourcequotas,verbs=get;watch;list
