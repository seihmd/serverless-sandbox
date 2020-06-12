.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/app app/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh

package:
	sls package

#local-invoke:
#	SLS_DEBUG=* sls invoke local --function app

#local-s3-mb:
#	awslocal s3 mb s3://fcm-delivery-service-local-bucket

localstack:
	make local-destroy ; TMPDIR=/private$$TMPDIR docker-compose up -d

local-s3-mb:
	awslocal s3 mb s3://fcm-delivery-service-local-bucket

local-sqs-createQueue:
	awslocal sqs create-queue --queue-name=FcmDeliveryPlatformTargetQueue-local

local-deploy: clean build
	sls deploy \
	  --stage=local \
	  --region=us-east-1 \
	  --bucket=fcm-delivery-service-local-bucket

local-deploy-func: clean build
	sls deploy \
	  --function app \
	  --stage local \
	  --region us-east-1 \
	  --bucket=fcm-delivery-service-local-bucket

local-destroy:
	docker rm -f $$(docker ps -aqf "name=fcm-delivery-platform_localstack_1")

local-invoke:
	awslocal sqs send-message --queue-url=http://localhost:4576/queue/FcmDeliveryPlatformTargetQueue-local --message-body="{}"