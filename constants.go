package main

const (
	UPLOAD_FILE_BUTTON     = "upload_file"
	DOWNLOAD_FILE_BUTTON   = "download_file"
	AUTH_USER              = "auth_user"
	WORKING_WITH_FOLDERS   = "working_with_folders"
	VIEW_FOLDERS_AND_FILES = "view_folders_and_files"

	UPDATE_CONFIG_TIMEOUT = 60

	CREATE_FOLDER_URL_TEMPLATE              = "https://cloud-api.yandex.net/v1/disk/resources?path="
	DELETE_FOLDER_URL_TEMPLATE              = "https://cloud-api.yandex.net/v1/disk/resources?path="
	COPY_FOLDER_OR_FILE_URL_TEMPLATE        = "https://cloud-api.yandex.net/v1/disk/resources/copy?"
	MOVE_FOLDER_OR_FILE_URL_TEMPLATE        = "https://cloud-api.yandex.net/v1/disk/resources/move?"
	FETCH_FOLDER_CONTENT_URL_TEMPLATE       = "https://cloud-api.yandex.net/v1/disk/resources?path="
	DOWNLOAD_FILE_URL_TEMPLATE              = "https://cloud-api.yandex.net/v1/disk/resources/download?path="
	SHOW_INFORMATION_DISK_URL_TEMPLATE      = "https://cloud-api.yandex.net/v1/disk/"
	CLEAN_TRASH_URL_TEMPLATE                = "https://cloud-api.yandex.net/v1/disk/trash/resources"
	RESTORE_CONTENT_FROM_TRASH_URL_TEMPLATE = "https://cloud-api.yandex.net/v1/disk/trash/resources/restore"
	FETCH_USER_FILES_URL_TEMPLATE           = "https://cloud-api.yandex.net/v1/disk/resources/files"

	TELEGRAM_API_TEMPLATE = "https://api.telegram.org/bot"

	KMS_CREATE_KEY_TEMPLATE = "https://kms.api.cloud.yandex.net/kms/v1/keys"
	KMS_GET_KEY_TEMPLATE    = "https://kms.api.cloud.yandex.net/kms/v1/keys"
	KMS_LIST_KEYS_TEMPLATE  = "https://kms.api.cloud.yandex.net/kms/v1/keys"

	LOCKBOX_CREATE_SECRET_TEMPLATE = "https://lockbox.api.cloud.yandex.net/lockbox/v1/secrets"
	LOCKBOX_ADD_VERSION_TEMPLATE   = "https://lockbox.api.cloud.yandex.net/lockbox/v1/secrets"
	LOCKBOX_LIST_SECRETS_TEMPLATE  = "https://lockbox.api.cloud.yandex.net/lockbox/v1/secrets"
	LOCKBOX_GET_PAYLOAD_TEMPLATE   = "https://payload.lockbox.api.cloud.yandex.net/lockbox/v1/secrets"

	GET_IAM_TOKEN_TEMPLATE = "https://iam.api.cloud.yandex.net/iam/v1/tokens"

	RUN_LOG_EXPORT_TEMPLATE = "https://logging.api.cloud.yandex.net/logging/v1/run-export"
)
