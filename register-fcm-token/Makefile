.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/registerFcmTokens registerFcmTokens/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh

localstack:
	make local-destroy ; TMPDIR=/private$$TMPDIR docker-compose up -d

local-setup: local-s3-mb local-sns-create-platform-application local-dynamodb-create-table local-ssm-put-parameter local-deploy

local-s3-mb:
	awslocal s3 mb s3://register-fcm-token-service-local-deployment-bucket

local-sns-create-platform-application:
	awslocal sns create-platform-application \
	--name fcm_sns \
	--platform GCM \
	--attributes PlatformCredential=test

local-dynamodb-create-table:
	awslocal dynamodb create-table \
	--table-name FcmTable \
	--attribute-definitions AttributeName=user_id,AttributeType=S \
	--key-schema AttributeName=user_id,KeyType=HASH \
	--provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1 \

local-ssm-put-parameter:
	awslocal ssm put-parameter \
		--name /fcm/local/sns/arn \
		--value arn:aws:sns:us-east-1:000000000000:app/GCM/fcm_sns \
		--type String \
	&& awslocal ssm put-parameter \
		--name /fcm/local/dynamodb/table/name \
		--value FcmTable \
		--type String

local-deploy: clean build
	sls deploy --stage local -r us-east-1

local-deploy-func: clean build
	sls deploy --function registerFcmTokens --stage local --region us-east-1

local-destroy:
	docker rm -f $$(docker ps -aqf "name=register-fcm-token_localstack_1")