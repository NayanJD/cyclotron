tmux:
	K6_INSTANCE_ID=$(terraform output -raw k6_instance_id)
	POSTGRES_INSTANCE_ID=$(terraform output -raw postgres_instance_id)
	USERSVC_INSTANCE_ID=$(terraform output -raw usersvc_instance_id)
	
	tmuxinator start cyclotron
