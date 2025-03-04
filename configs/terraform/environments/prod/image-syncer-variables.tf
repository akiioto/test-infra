variable "image_syncer_reusable_workflow_ref" {
  type        = string
  description = "The value of GitHub OIDC token job_workflow_ref claim of the image-syncer reusable workflow in the test-infra repository. This is used to identify token exchange requests for image-syncer reusable workflow."
  default     = "kyma-project/test-infra/.github/workflows/image-syncer.yml@refs/heads/main"
}

variable "image_syncer_reader_service_account_name" {
  type        = string
  description = "The service account name of the image-syncer service account."
  default     = "image-syncer-reader"
}

variable "image_syncer_writer_service_account_name" {
  type        = string
  description = "The service account name of the image-syncer service account."
  default     = "image-syncer-writer"
}# (2025-03-04)