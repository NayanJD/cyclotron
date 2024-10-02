module "cyclotron-k6-ec2" {
  source = "terraform-aws-modules/ec2-instance/aws"

  name = "k6-instance"
  
  ami = "ami-0f27560e987a14b7a"

  subnet_id              = element(module.vpc.private_subnets, 0)
  vpc_security_group_ids = [module.security_group_instance.security_group_id]
  instance_type          = "m5a.large"

  create_iam_instance_profile = true
  iam_role_description        = "IAM role for EC2 instance"
  iam_role_policies = {
    AmazonSSMManagedInstanceCore = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
  }
  
  user_data = file("${path.module}/k6.sh")
  user_data_replace_on_change = true

  tags = merge(local.tags, {
    Name = "cyclotron-k6"
  })
}

output "k6_instance_id" {
  value = module.cyclotron-k6-ec2.id
}

module "security_group_instance" {
  source  = "terraform-aws-modules/security-group/aws"
  version = "~> 5.0"

  name        = "${local.name}-ec2"
  description = "Security Group for EC2 Instance Egress"

  vpc_id = module.vpc.vpc_id
    
  ingress_cidr_blocks = [module.vpc.vpc_cidr_block]
  ingress_rules = ["all-all"]
  egress_rules = ["https-443-tcp","all-all"]

  tags = local.tags
}

resource "aws_iam_role_policy" "test_policy" {
  name = "test_policy"
  role = module.cyclotron-k6-ec2.iam_role_name

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
      }
    ]
  })
}

resource "aws_volume_attachment" "jaeger-volume-attachment" {
  device_name = "/dev/sdj"
  volume_id   = aws_ebs_volume.jaeger-volume.id
  instance_id = module.cyclotron-k6-ec2.id
}

resource "aws_ebs_volume" "jaeger-volume" {
  availability_zone = module.cyclotron-k6-ec2.availability_zone
  size              = 10

  tags = local.tags
}

resource "aws_volume_attachment" "prometheus-volume-attachment" {
  device_name = "/dev/sdp"
  volume_id   = aws_ebs_volume.prometheus-volume.id
  instance_id = module.cyclotron-k6-ec2.id
}

resource "aws_ebs_volume" "prometheus-volume" {
  availability_zone = module.cyclotron-k6-ec2.availability_zone
  size              = 10

  tags = local.tags
}
