# Quota Operator

Enforcing per team quota (sum of used resources across all their namespaces) and delegating the per namespace quota to users.

## Instructions

### Development

* `make generate` update the generated code for that resource type.
* `make manifests` Generating CRD manifests.
* `make test` Run tests.

### Build

Export your image name:

```
export IMG=ghcr.io/your-repo-path/image-name:latest
```

* `make build` builds golang app locally.
* `make docker-build` build docker image locally.
* `make docker-push` push container image to registry.

### Run, Deploy
* `make run` run app locally
* `make deploy` deploy to k8s.

### Clean up

* `make undeploy` delete resouces in k8s.


## Metrics

| Metric                                              | Notes
|-----------------------------------------------------|------------------------------------
| controller_runtime_active_workers | Number of currently used workers per controller
| controller_runtime_max_concurrent_reconciles | Maximum number of concurrent reconciles per controller
| controller_runtime_reconcile_errors_total | Total number of reconciliation errors per controller
| controller_runtime_reconcile_time_seconds | Length of time per reconciliation per controller
| controller_runtime_reconcile_total | Total number of reconciliations per controller
| rest_client_request_latency_seconds | Request latency in seconds. Broken down by verb and URL.
| rest_client_requests_total | Number of HTTP requests, partitioned by status code, method, and host.
| workqueue_adds_total | Total number of adds handled by workqueue
| workqueue_depth | Current depth of workqueue
| workqueue_longest_running_processor_seconds | How many seconds has the longest running processor for workqueue been running.
| workqueue_queue_duration_seconds | How long in seconds an item stays in workqueue before being requested
| workqueue_retries_total | Total number of retries handled by workqueue
| workqueue_unfinished_work_seconds | How many seconds of work has been done that is in progress and hasn't been observed by work_duration. Large values indicate stuck threads. One can deduce the number of stuck threads by observing the rate at which this increases.
| workqueue_work_duration_seconds | How long in seconds processing an item from workqueue takes.


## Security

### Reporting security vulnerabilities

If you find a security vulnerability or any security related issues, please DO NOT file a public issue, instead send your report privately to cloud@snapp.cab. Security reports are greatly appreciated and we will publicly thank you for it.

## License

Apache-2.0 License, see [LICENSE](LICENSE).
