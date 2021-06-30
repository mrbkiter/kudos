#!/bin/bash
echo "running go test"
testResult=$(go test . | grep FAIL)
if [[ $testResult == *"FAIL"* ]]; then
  echo "Run test failed. Please fix all issues below"
  echo $testResult
else 
    GOOS=linux go build -o report .
    zip -r main.zip ./report ./config/*
    echo DONE ZIP
    aws lambda update-function-code --function-name  dev-kudos-report --zip-file fileb://./main.zip
fi

