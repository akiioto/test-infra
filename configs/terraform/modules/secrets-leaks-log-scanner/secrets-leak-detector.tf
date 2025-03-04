data "google_iam_policy" "run_invoker" {
  binding {
    role    = "roles/run.invoker"
    members = ["serviceAccount:${google_service_account.secrets_leak_detector.email}"]
  }
}

# Create a service account for Eventarc trigger and Workflows
resource "google_service_account" "secrets_leak_detector" {
  account_id  = "secrets-leak-detector"
  description = "Identity of secrets leak detector application."
}

# Grant the logWriter role to the service account
resource "google_project_iam_member" "project_log_writer" {
  member  = "serviceAccount:${google_service_account.secrets_leak_detector.email}"
  project = data.google_project.project.id
  role    = "roles/logging.logWriter"
}

# Grant the workflows.invoker role to the service account
resource "google_project_iam_member" "project_workflows_invoker" {
  project = data.google_project.project.id
  role    = "roles/workflows.invoker"
  member  = "serviceAccount:${google_service_account.secrets_leak_detector.email}"
}

locals {
  scan_logs_for_secrets_yaml = templatefile(("${path.module}/../../../../pkg/gcp/workflows/secrets-leak-detector.yaml"), {
    scan-logs-for-secrets-url = google_cloud_run_service.secrets_leak_log_scanner.status[0].url
    move-gcs-bucket-url       = google_cloud_run_service.gcs_bucket_mover.status[0].url
    search-github-issue-url   = google_cloud_run_service.github_issue_finder.status[0].url
    create-github-issue-url   = google_cloud_run_service.github_issue_creator.status[0].url
    send-slack-message-url    = var.slack_message_sender_url
    }
  )
}


resource "google_workflows_workflow" "secrets_leak_detector" {
  name            = "secrets-leak-detector"
  region          = "europe-west3"
  description     = "Workflow is triggered by message published to prowjobs PubSub topic and scans prowjobs logs for secrets."
  service_account = google_service_account.secrets_leak_detector.id
  source_contents = local.scan_logs_for_secrets_yaml
}

resource "google_eventarc_trigger" "secrets_leak_detector_workflow" {
  name     = "secrets-leak-detector"
  location = "europe-west3"
  matching_criteria {
    attribute = "type"
    value     = "google.cloud.pubsub.topic.v1.messagePublished"
  }
  destination {
    workflow = google_workflows_workflow.secrets_leak_detector.id
  }
  service_account = google_service_account.secrets_leak_detector.id
  labels = {
    application = "secrets_leak_detector"
  }
  transport {
    pubsub {
      topic = "projects/${var.gcp_project_id}/topics/${var.prow_pubsub_topic_name}"
    }
  }
}
# (2025-03-04)