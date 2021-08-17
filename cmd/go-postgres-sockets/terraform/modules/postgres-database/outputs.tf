output "postgres_server_ip" {
  value = google_sql_database_instance.instance.ip_address
}

output "postgres_server_instance_connection_name" {
  value = google_sql_database_instance.instance.connection_name
}