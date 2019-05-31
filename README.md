# Mock AWS SDK Example

A simple example showing how to mock the Go AWS SDK using the provided interfaces.
By using an interface, you can easily mock responses from the SDK in your unit tests.

This example uses the CloudWatch API to get the estimated size of an S3 bucket.

## Running

```
go run main.go s3 --bucket my-bucket
```

## Tests

```
go test
```
