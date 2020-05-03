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

setup:
	TMPDIR=/private$TMPDIR docker-compose up -d

setup-local:
	awslocal s3 mb s3://register-fcm-token-service-local-deployment-bucket \
	&& awslocal sns create-platform-application --name fcm_sns --platform GCM --attributes PlatformCredential=test \
	&& awslocal dynamodb create-table \
		--table-name FcmTable \
		--attribute-definitions AttributeName=user_id,AttributeType=S\
		--key-schema AttributeName=user_id,KeyType=HASH \
		--provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1

deploy-local: clean build
	sls deploy --function registerFcmTokens --stage local
