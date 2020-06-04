IAM_ID := 334927132216

.PHONY: enter
enter:
	docker exec -it symfony-fargate_php_1 /bin/sh

.PHONY: compose-local
compose-local:
	docker-compose -f docker-compose.yml up -d --build

.PHONY: build-nginx
build-nginx:
	cd docker/nginx && \
	docker build . -t symfony-fargate-nginx -f ./Dockerfile && \
	docker tag symfony-fargate-nginx:latest $(IAM_ID).dkr.ecr.ap-northeast-1.amazonaws.com/symfony-fargate-nginx:latest && \
	docker push $(IAM_ID).dkr.ecr.ap-northeast-1.amazonaws.com/symfony-fargate-nginx:latest

.PHONY: build-php
build-php:
	docker build ./app -t symfony-fargate-php -f ./app/Dockerfile && \
	docker tag symfony-fargate-php:latest $(IAM_ID).dkr.ecr.ap-northeast-1.amazonaws.com/symfony-fargate-php:latest && \
	docker push $(IAM_ID).dkr.ecr.ap-northeast-1.amazonaws.com/symfony-fargate-php:latest

.PHONY: compose-dev
compose-dev:
	IAM_ID=$(IAM_ID) ecs-cli compose -f docker-compose.dev.yml up;

.PHONY: ecs-setup
ecs-setup:
	ecs-cli up --keypair symfony_fargate_dev  --instance-type t2.micro --capability-iam --force

.PHONY: ecs-auth
ecs-auth:
	aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin $(IAM_ID).dkr.ecr.ap-northeast-1.amazonaws.com