ENV?=default

all: init workspace plan apply pause deploy

init:
	cd ./terraform/demo && terraform init

workspace:
	cd ./terraform/demo && terraform workspace new ${ENV}

set_workspace:
	cd ./terraform/demo && terraform workspace select ${ENV}

plan: set_workspace
	cd ./terraform/demo && terraform plan

apply: set_workspace
	cd ./terraform/demo && terraform apply -auto-approve

destroy: set_workspace
	cd ./terraform/demo && terraform destroy -auto-approve

clean:
	cd ./terraform/demo &&  rm -f terraform.tfplan
	cd ./terraform/demo &&  rm -f .terraform.lock.hcl
	cd ./terraform/demo &&  rm -fr terraform.tfstate*
	cd ./terraform/demo &&  rm -fr .terraform/

pause:
	echo "Wait for 60 seconds stupid Yandex Cloud creating a VM's..."
	sleep 60
	echo "May be created? Ok, run an ansible playbook..."
deploy:
	cd ./ansible && source .env.news-app && ansible-playbook -i inventory/demo site.yml

reconfig:
	cd ./ansible && source .env.news-app && ansible-playbook -i inventory/demo site.yml -t config