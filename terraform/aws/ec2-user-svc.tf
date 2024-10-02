data "aws_ami" "user_svc_ubuntu" {
  most_recent      = true
  owners           = ["amazon"]

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-*"]
  }

  filter {
    name   = "root-device-type"
    values = ["ebs"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

module "cyclotron_user_service_ec2" {
  source = "terraform-aws-modules/ec2-instance/aws"

  name = "user-service-instance"
  
  # bitnami-postgresql-15.4.0-7-r13-linux-debian-11-x86_64-hvm-ebs-nami
  ami = data.aws_ami.user_svc_ubuntu.image_id
  
  # associate_public_ip_address = true

  subnet_id              = element(module.vpc.private_subnets, 0)
  vpc_security_group_ids = [module.user_svc_security_group_instance.security_group_id]
  instance_type          = "c5a.large"

  create_iam_instance_profile = true
  iam_role_description        = "IAM role for EC2 instance"
  iam_role_policies = {
    AmazonSSMManagedInstanceCore = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
  }
  
  user_data = file("${path.module}/user-svc.sh")
  user_data_replace_on_change = true

  tags = merge(local.tags, {
    Name = "cyclotron-user-svc"
    POSTGRES_URL = module.cyclotron-postgres-ec2.private_ip
    JAEGER_URL = module.cyclotron-k6-ec2.private_ip
  })
}

output "usersvc_instance_id" {
  value = module.cyclotron_user_service_ec2.id
}

module "user_svc_security_group_instance" {
  source  = "terraform-aws-modules/security-group/aws"
  version = "~> 5.0"

  name        = "${local.name}-user-svc-ec2"
  description = "Security Group for EC2 Instance Egress"

  vpc_id = module.vpc.vpc_id
    
  ingress_cidr_blocks = [module.vpc.vpc_cidr_block]
  ingress_rules = ["all-all"]
  egress_rules = ["all-all"]

  tags = local.tags
}

resource "aws_iam_role_policy" "user-svc-policy" {
  name = "user-svc-policy"
  role = module.cyclotron_user_service_ec2.iam_role_name

  # Terraform's "jsonencode" function converts a
  # Terraform expression result to valid JSON syntax.
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "secretsmanager:GetResourcePolicy",
          "secretsmanager:GetSecretValue",
          "secretsmanager:DescribeSecret",
          "secretsmanager:ListSecretVersionIds"
        ]
        Effect   = "Allow"
        Resource = [
                "arn:aws:secretsmanager:${local.region}:${data.aws_caller_identity.current.account_id}:secret:cyclotron/commons-T2SS5b",
            ]
      },
      {
        Action = "secretsmanager:ListSecrets",
        Effect = "Allow",
        Resource = "*"
      },
      {
        Action = "ec2:DescribeTags",
        Effect = "Allow",
        Resource = "*"
      }
    ]
  })
}

resource "aws_ec2_tag" "postgres-user-svc-tag" {
  resource_id = module.cyclotron-postgres-ec2.id
  key         = "USER_SVC_URL"
  value       = module.cyclotron_user_service_ec2.private_ip
}
