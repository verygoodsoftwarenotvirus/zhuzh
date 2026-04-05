locals {
  api_database_username                   = "api_db_user"
  async_message_handler_database_username = "async_message_handler"
  db_cleaner_username                     = "db_cleaner"
  search_data_index_scheduler_username    = "search_data_index_scheduler"
  mobile_notification_scheduler_username  = "mobile_notification_scheduler"
  queue_test_username                     = "queue_test"
}

# api_database_username

resource "random_password" "api_user_database_password" {
  length           = 64
  special          = true
  override_special = "#$*-_=+[]"
}

resource "google_sql_user" "api_user" {
  name     = local.api_database_username
  instance = google_sql_database_instance.prod.name
  password = random_password.api_user_database_password.result
}

# async_message_handler_database_username

resource "random_password" "async_message_handler_database_user_database_password" {
  length           = 64
  special          = true
  override_special = "#$*-_=+[]"
}

resource "google_sql_user" "async_message_handler_database_user" {
  name     = local.async_message_handler_database_username
  instance = google_sql_database_instance.prod.name
  password = random_password.async_message_handler_database_user_database_password.result
}

# db_cleaner_username

resource "random_password" "db_cleaner_user_database_password" {
  length           = 64
  special          = true
  override_special = "#$*-_=+[]"
}

resource "google_sql_user" "db_cleaner_user" {
  name     = local.db_cleaner_username
  instance = google_sql_database_instance.prod.name
  password = random_password.db_cleaner_user_database_password.result
}

# search_data_index_scheduler_username

resource "random_password" "search_data_index_scheduler_user_database_password" {
  length           = 64
  special          = true
  override_special = "#$*-_=+[]"
}

resource "google_sql_user" "search_data_index_scheduler_user" {
  name     = local.search_data_index_scheduler_username
  instance = google_sql_database_instance.prod.name
  password = random_password.search_data_index_scheduler_user_database_password.result
}

# mobile_notification_scheduler_username

resource "random_password" "mobile_notification_scheduler_user_database_password" {
  length           = 64
  special          = true
  override_special = "#$*-_=+[]"
}

resource "google_sql_user" "mobile_notification_scheduler_user" {
  name     = local.mobile_notification_scheduler_username
  instance = google_sql_database_instance.prod.name
  password = random_password.mobile_notification_scheduler_user_database_password.result
}

# queue_test_username

resource "random_password" "queue_test_user_database_password" {
  length           = 64
  special          = true
  override_special = "#$*-_=+[]"
}

resource "google_sql_user" "queue_test_user" {
  name     = local.queue_test_username
  instance = google_sql_database_instance.prod.name
  password = random_password.queue_test_user_database_password.result
}
