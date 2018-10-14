# ENV required
# We also need to have 
# AWS CLI credential setup
# AWS_ACCESS_KEY_ID
# AWS_SECRET_ACCESS_KEY
#
# and 
# S3_BUCKET 

# Vars
PACKAGE_NAME = adlmhelper
PACKAGE_DIR = ./
PACKAGE_TEMPLATE = template.yaml
PACKAGE_OUTPUT_TEMPLATE = packaged.yaml
STACK_NAME = sam-adlm-helper

.PHONY: default clean build package deploy
default: clean build package deploy	

clean: build_clean template_clean 

build_clean: 
	rm -rf ./build

template_clean:
	rm -rf ${PACKAGE_OUTPUT_TEMPLATE}

test:
	go test -v ./...

build:
	GOOS=linux GOARCH=amd64 go build -o ./build/${PACKAGE_NAME} ${PACKAGE_DIR}

package:
	sam package --template-file ${PACKAGE_TEMPLATE} --output-template-file ${PACKAGE_OUTPUT_TEMPLATE} --s3-bucket ${S3_BUCKET}

deploy:
	sam deploy --template-file ${PACKAGE_OUTPUT_TEMPLATE} --stack-name ${STACK_NAME} --capabilities CAPABILITY_IAM --capabilities CAPABILITY_NAMED_IAM

binary: build_clean build

sam: template_clean package deploy
