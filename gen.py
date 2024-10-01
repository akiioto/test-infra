import os

koapps_yaml = """
apps:
  - ko://github.com/kyma-project/test-infra/cmd/tools/pjtester
  - ko://github.com/kyma-project/test-infra/cmd/image-url-helper
  - ko://github.com/kyma-project/test-infra/cmd/markdown-index
  - ko://github.com/kyma-project/test-infra/cmd/tools/usersmapchecker
  - ko://github.com/kyma-project/test-infra/cmd/tools/gcscleaner
  - ko://github.com/kyma-project/test-infra/cmd/tools/diskscollector
  - ko://github.com/kyma-project/test-infra/cmd/tools/clusterscollector
  - ko://github.com/kyma-project/test-infra/cmd/tools/vmscollector
  - ko://github.com/kyma-project/test-infra/cmd/tools/ipcleaner
  - ko://github.com/kyma-project/test-infra/cmd/tools/orphanremover
  - ko://github.com/kyma-project/test-infra/cmd/tools/dnscollector
  - ko://github.com/kyma-project/test-infra/cmd/tools/externalsecretschecker
  - ko://github.com/kyma-project/test-infra/cmd/cloud-run/create-github-issue
  - ko://github.com/kyma-project/test-infra/cmd/cloud-run/move-gcs-bucket
  - ko://github.com/kyma-project/test-infra/cmd/cloud-run/scan-logs-for-secrets
  - ko://github.com/kyma-project/test-infra/cmd/cloud-run/search-github-issue
  - ko://github.com/kyma-project/test-infra/cmd/cloud-run/github-webhook-gateway
  - ko://github.com/kyma-project/test-infra/cmd/cloud-run/cors-proxy
  - ko://github.com/kyma-project/test-infra/cmd/external-plugins/automated-approver
  - ko://github.com/kyma-project/test-infra/cmd/dashboard-token-proxy
"""

apps = [line.split("ko://github.com/kyma-project/test-infra/cmd/")[1].strip() for line in koapps_yaml.splitlines() if line.startswith("  - ko://")]

dockerfile_template = """FROM golang:1.23-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source to the Working Directory inside the container
COPY . .

WORKDIR /app/cmd/{app_full_path}

# Build the Go app with static linking
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:3.20.3

LABEL io.kyma-project.source=github.com/kyma-project/test-infra/cmd/{app_full_path}

# Copy the built Go app from the builder stage
COPY --from=builder /app/cmd/{app_full_path}/main /{app_name}

RUN apk add --no-cache ca-certificates git && \\
    chmod +x /{app_name}
ENTRYPOINT ["/{app_name}"]
"""

workflow_template_pull = """name: build-{app_name}
on:
  pull_request_target:
    types: [ opened, edited, synchronize, reopened, ready_for_review ]
    paths:
      - "cmd/{app_full_path}/*.go"
      - "cmd/{app_full_path}/Dockerfile"
      - "pkg/**"
      - "go.mod"
      - "go.sum"
  push:
    branches:
      - main
    paths:
      - "cmd/{app_full_path}/*.go"
      - "cmd/{app_full_path}/Dockerfile"
      - "pkg/**"
      - "go.mod"
      - "go.sum"
  workflow_dispatch: {{}}

jobs:
  build-image:
    uses: ./.github/workflows/image-builder.yml
    with:
      name: {app_name}
      dockerfile: cmd/{app_full_path}/Dockerfile
      context: .
  print-image:
    runs-on: ubuntu-latest
    needs: build-image
    steps:
      - name: Print image
        run: echo "Image built ${{{{ needs.build-image.outputs.images }}}}"
"""

for app in apps:
    app_full_path = app
    app_name = app_full_path.split('/')[-1]

    dockerfile_dir = os.path.join("cmd", app_full_path)
    dockerfile_path = os.path.join(dockerfile_dir, "Dockerfile")

    os.makedirs(dockerfile_dir, exist_ok=True)

    with open(dockerfile_path, 'w') as dockerfile:
        dockerfile.write(dockerfile_template.format(app_name=app_name, app_full_path=app_full_path))

    workflow_pull_path = f".github/workflows/build-{app_name}.yml"
    with open(workflow_pull_path, 'w') as workflow_pull:
        workflow_pull.write(workflow_template_pull.format(app_name=app_name, app_full_path=app_full_path))

print("git")
