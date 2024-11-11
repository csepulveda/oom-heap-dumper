# Memory Profiler Sidecar for Kubernetes

This application is designed to be used as a sidecar container in Kubernetes pods that experience Out of Memory (OOM) issues. The app monitors processes within the pod, identifying any process (PID) that exceeds a defined memory usage percentage. When a process surpasses the designated threshold, the application triggers a memory heap dump using the default Go `pprof` endpoint and uploads the heap dump to an S3 bucket for analysis.

## Features

- **Threshold-Based Monitoring**: Continuously monitors memory usage by process and identifies any process exceeding the critical memory percentage.
- **Automatic Heap Dump Collection**: Calls the Go `pprof` endpoint to collect a memory heap dump when a process crosses the memory threshold.
- **S3 Storage Integration**: Uploads the collected memory heap dump to an S3 bucket, making it easy to retrieve and analyze the file.

## Requirements

- AWS credentials with permissions to upload to the specified S3 bucket.
- **A web service running with `pprof` enabled on the application to be monitored**. This sidecar relies on the monitored application exposing a `pprof` endpoint, accessible at `/debug/pprof/heap`.
- Environment variables:
  - `AWS_REGION`: AWS region of the S3 bucket.
  - `BUCKET`: Name of the S3 bucket where memory heap dumps will be stored.

## Security Requirements

To monitor processes within the Kubernetes pod, this application requires elevated privileges. Ensure that you set `securityContext.privileged: true` in the Kubernetes configuration for the sidecar container.

Additionally, enable the `shareProcessNamespace` setting in the pod specification. This allows the sidecar container to access process information across all containers within the pod. [For more details, refer to the Kubernetes documentation on sharing the process namespace.](https://kubernetes.io/docs/tasks/configure-pod-container/share-process-namespace/)

## Environment Variables

The following environment variables are required for the application to function:

- **`AWS_REGION`**: The AWS region where the S3 bucket is located (e.g., `us-west-2`).
- **`BUCKET`**: The name of the S3 bucket for storing memory heap dumps.

## Usage

To deploy this application as a sidecar container in your Kubernetes pod, add it to your pod's container specification. Here’s an example configuration:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: memory-profiler
spec:
  shareProcessNamespace: true
  containers:
    - name: main-app
      image: your-main-app-image
      # Main application container
    - name: memory-profiler-sidecar
      image: your-sidecar-image
      securityContext:
        privileged: true
      env:
        - name: AWS_REGION
          value: "us-west-2"
        - name: BUCKET
          value: "your-s3-bucket-name"
      # Add resource limits if needed
```

## Configuration

The application can be configured to set the memory usage threshold that triggers a heap dump. Set this threshold as an environment variable (`MEMORY_CRITICAL_PERCENTAGE`). For example:

```yaml
- name: MEMORY_CRITICAL_PERCENTAGE
  value: "80"
```

The above setting means the application will trigger a heap dump if any process exceeds 80% of memory usage.

## How It Works

1. The application continuously monitors processes within the Kubernetes pod.
2. When a process exceeds the specified critical memory percentage, it calls the Go `pprof` endpoint (`/debug/pprof/heap`) to collect a memory heap dump.
3. The heap dump is saved with a filename indicating the PID and memory usage percentage and is uploaded to the specified S3 bucket for further analysis.

## Dependencies

- [AWS SDK for Go](https://aws.amazon.com/sdk-for-go/)
- [Go Fiber](https://gofiber.io/) (if applicable, based on your setup)

## Example

To test the application locally or in a Kubernetes environment, set up your AWS credentials and specify the required environment variables:

```bash
export AWS_REGION=us-west-2
export BUCKET=my-s3-bucket
export MEMORY_CRITICAL_PERCENTAGE=80
export WATCH_TIME=500
```

Then run the application to start monitoring memory usage.

## Notes

- This sidecar container should be used alongside applications prone to high memory usage or memory leaks.
- Ensure the sidecar has permissions to write to the specified S3 bucket.
- Access to the Go `pprof` endpoint should be secure and limited to prevent unauthorized access to heap dumps.

## Credits

This application it's based on this awesome proyect: https://github.com/ricardomaraschini/oomhero

## License

This project is open-source and available under the [Apache License: Version 2](LICENSE).
