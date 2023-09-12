#!/bin/bash

BASE_URL="http://localhost:8080/"


clean_and_terminate(){
  rm -rf testOutput
  exit 1
}

mkdir testOutput

for file in "./testData"/*; do
  if [ -f "$file" ]; then

    file_name=$(basename "$file")
    url_path="${file#$TEST_DATA_FOLDER/}"

    curl --request POST --data-binary "@$file" -s "$BASE_URL$url_path" >/dev/null

    # get the saved file from file server
    curl --request GET -s "$BASE_URL$url_path" >testOutput/$file_name

    # compare the file from server with original file
    if diff -q testData/$file_name testOutput/$file_name >/dev/null; then
      echo "Test [$file_name] is PASSED"
    else
      echo "Test [$file_name] is FAILED"
      clean_and_terminate
    fi

    curl --request DELETE -s "$BASE_URL$url_path" >/dev/null

    # make sure the file is not hosted on the server
    status_code=$(curl --request GET -o /dev/null -s -I -w "%{http_code}" "$BASE_URL$url_path")
    if [[ $status_code != "404" ]]; then
      echo "ERROR! status_code isn't 404"
      clean_and_terminate
    fi
  fi
done
rm -rf testOutput
