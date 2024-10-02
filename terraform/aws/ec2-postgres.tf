module "cyclotron-postgres-ec2" {
  source = "terraform-aws-modules/ec2-instance/aws"

  name = "postgres-instance"
  
  # bitnami-postgresql-15.4.0-7-r13-linux-debian-11-x86_64-hvm-ebs-nami
  ami = "ami-03b0a3601c699360f"
  
  key_name = "cyclotron"
  # associate_public_ip_address = true

  subnet_id              = element(module.vpc.private_subnets, 0)
  vpc_security_group_ids = [module.postgres_security_group_instance.security_group_id]
  instance_type          = "r5a.large"

  create_iam_instance_profile = true
  iam_role_description        = "IAM role for EC2 instance"
  iam_role_policies = {
    AmazonSSMManagedInstanceCore = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
  }
  
  user_data = file("${path.module}/postgres.sh")
  user_data_replace_on_change = true

  tags = merge(local.tags, {
    Name = "cyclotron-postgres"
  })
}

output "postgres_instance_id" {
  value = module.cyclotron-postgres-ec2.id
}

module "postgres_security_group_instance" {
  source  = "terraform-aws-modules/security-group/aws"
  version = "~> 5.0"

  name        = "${local.name}-postgres-ec2"
  description = "Security Group for EC2 Instance Egress"

  vpc_id = module.vpc.vpc_id

  ingress_cidr_blocks = [module.vpc.vpc_cidr_block]
  ingress_rules = ["postgresql-tcp"]
  egress_rules = ["https-443-tcp", "all-all"]
  # ingress_with_cidr_blocks = [
  #   {
  #     from_port = 22
  #     protocol  = "tcp"
  #     to_port = 22
  #     cidr_blocks = "0.0.0.0/0"
  #     description = "SSH"
  #   }
  # ]
  # egress_with_cidr_blocks = [
  # {
  #     from_port   = 8080
  #     to_port     = 8090
  #     protocol    = "tcp"
  #     description = "User-service ports"
  #     cidr_blocks = "10.10.0.0/16"
  #   }
  # ]
  tags = local.tags
}

# Not used yet
# resource "aws_volume_attachment" "postgres-volume-attachment" {
#   device_name = "/dev/xvdj"
#   volume_id   = aws_ebs_volume.postgres-volume.id
#   instance_id = module.cyclotron-postgres-ec2.id
# }

# resource "aws_ebs_volume" "postgres-volume" {
#   availability_zone = module.cyclotron-postgres-ec2.availability_zone
#   size              = 5

#   tags = local.tags
# }

resource "aws_iam_role_policy" "cycltoron-postgres-instance-policy" {
  name = "test_policy"
  role = module.cyclotron-postgres-ec2.iam_role_name

  # Terraform's "jsonencode" function converts a
  # Terraform expression result to valid JSON syntax.
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "secretsmanager:GetResourcePolicy",
                "secretsmanager:*",
        ]
        Effect   = "Allow"
        Resource = [
                "arn:aws:secretsmanager:${local.region}:${data.aws_caller_identity.current.account_id}:secret:cyclotron/commons-??????",
            ]
      },
      {
        Action = "secretsmanager:ListSecrets",
        Effect = "Allow",
        Resource = "*"
      }
    ]
  })
}
