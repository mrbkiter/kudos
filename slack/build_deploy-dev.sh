#!/bin/bash
echo "running go test"
testResult=$(go test . | grep FAIL)
if [[ $testResult == *"FAIL"* ]]; then
  echo "Run test failed. Please fix all issues below"
  echo $testResult
else 
    GOOS=linux go build -o slack .
    zip -r main.zip ./slack ./config/*
    echo DONE ZIP
    aws lambda update-function-code --function-name  dev-kudos-slack --zip-file fileb://./main.zip
fi

