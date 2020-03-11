output "url" {
  value = "http://${aws_instance.web_server.public_ip}:8080"
}