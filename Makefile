# https://serverless.com/framework/docs/providers/aws/cli-reference/create/

.PHONY: create-goawsmod
create-aws-gomod:
	serverless create -t aws-go-mod -p newservice

.PHONY: create-aws-typescript
create-aws-typescript:
	serverless create -t aws-nodejs-typescript -p newservice
