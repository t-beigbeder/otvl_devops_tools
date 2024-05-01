output "ami" {
  value       = data.aws_ami.this
  description = "The latest filtered AMI"
}
