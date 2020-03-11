resource "aws_instance" "web_server" {
  ami                    = "ami-43a15f3e"
  instance_type          = "t2.micro"
  vpc_security_group_ids = [aws_security_group.web_server.id]
  # Run a "Hello, World" web server on port 8080
  user_data = <<-EOF
              #!/bin/bash
              echo "Hello, World" > index.html
              nohup busybox httpd -f -p 8080 &
              EOF 
}

resource "aws_security_group" "web_server" {
  ingress {
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
