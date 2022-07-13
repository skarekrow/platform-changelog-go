# Platform Changelog API

## Overview

The Platform Changelog is a system for keeping track of changes as they occur
across the platform by leverage different types of notification events, such as
Github and Gitlab webhooks, as well as Deployment Pipeline tasks.

This API will provide JSON responses to the requesting entity, mainly the [Platform
Changelog Frontend](https://www.github.com/redhatinsights/platform-changelog).

Initally, the service supports only Github webhooks authenticated via secret
token, but will eventually also support Gitlab and Deployments hooks from Tekton.

## Architecture

Platform Changelog is a backend API that connects to a backend database for storing
supported incoming events. The current implementation supports a Postgres database
and respondes to incoming requests with JSON responses.

A frontend application has also been developed for displaying this information in
an easy to read, and searchable manner.

## REST API Endpoint

TODO: API Spec

## Adding A Service

To add a service to be supported by platform-changelog, follow these steps:

1. Add the service to `internal/config/services.yaml`
  
  ```yaml
  service-name:
    display_name: "Service Name"
    gh_repo: <https://github.com/org/repo>
    branch: master # branch to be monitored
    namespace: <namespace of the project>
```

2. Submit an MR to this repo. It will be approved by an owner

## Development

A Makefile has been provided for most common operations to get the app up and running.
A compose file is also available for standing up the service in podman.

Docker can be substituted for podman if needed.

### Prequisites

    podman
    podman-compose
    Golang >= 1.16

### Launching

    $> make -B build
    $> make run-db
    $> make run-migration
    $> make run-api DEBUG=1

Note: The `DEBUG` argument allows us to send webhooks without needing the secret token.

The API should now be up and available on `localhost:8000`. You should be able to
see the API in action by visiting `http://localhost:8000/api/v1/services`.

### Testing Webhooks to the API Manually

Launch the API as instructed above, then we can send test webhooks to the API.

Test webhook json is provided in the `tests` directory in this repo.

To send the webhook, you can use curl or `make test-github-webhook`. The curl command is:

`curl -X POST -H "X-Github-Event: push" -H "Content-Type: application/json" --data "@tests/github_webhook.json" http://localhost:8000/api/v1/github-webhook`

From there, you should be able to open a browser and see the results at: http://localhost:8000/api/v1/commits. There should be commits matching the webhook data that was sent.

## Running Tests

TODO: Get some tests in here

# Get Help

This service is owned by the ConsoldeDot Pipeline team. If you have any questions, or
need support with this service, please contact the team on slack @crc-pipeline-team.

You can also raise an Issue in this repo for them to address.
