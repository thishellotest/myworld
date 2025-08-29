package biz

const Redis_sync_googlesheet_tasks_queue = "googlesheet_tasks_queue"
const Redis_sync_googlesheet_tasks_processing = "googlesheet_tasks_processing"
const Redis_sync_adobesign_tasks_queue = "adobesign_tasks_queue"
const Redis_sync_adobesign_tasks_processing = "adobesign_tasks_processing"
const Redis_sync_asana_tasks_queue = "asana_tasks_queue"
const Redis_sync_asana_tasks_processing = "asana_tasks_processing"
const Redis_sync_asana_users_queue = "asana_users_queue"
const Redis_sync_asana_users_processing = "asana_users_processing"

const Redis_record_review_tasks_queue = "record_review_tasks_queue"
const Redis_record_review_tasks_processing = "record_review_tasks_processing"

const Redis_client_task_handle_what_gid_queue = "client_task_handle_what_gid_queue"
const Redis_client_task_handle_what_gid_processing = "client_task_handle_what_gid_processing"
const Redis_client_task_handle_who_gid_queue = "client_task_handle_who_gid_queue"
const Redis_client_task_handle_who_gid_processing = "client_task_handle_who_gid_processing"

const Redis_client_name_change_job_queue = "client_name_change_job_queue"
const Redis_client_name_change_job_processing = "client_name_change_job_processing"

type RedisUsecase struct {
}

func NewRedisUsecase() *RedisUsecase {
	return &RedisUsecase{}
}
