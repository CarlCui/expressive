#!/bin/bash

echo "Starting e2e tests"

LLI_PATH=./llvm-prebuilt/lli
# LLI_PATH=lli # for local development

for file in ./e2e/*.exp; do
    echo $file

    filename=$(basename $file)
    filenameWithoutExt="${filename%.*}"

    expressive -d ./e2e -f $filename --outDir ./dist

    result=$($LLI_PATH ./dist/$filenameWithoutExt.s)

    for expectedFile in ./e2e/*.txt; do

        expectedFileBasename=$(basename $expectedFile)
        expectedFilenameWithoutExt="${expectedFileBasename%.*}"

        if [ "$filenameWithoutExt" = "$expectedFilenameWithoutExt" ]; then
            echo $expectedFile
            expected=$(cat $expectedFile)

            if [ "$result" != "$expected" ]; then
                echo Expected is $expected, but result is $result
                exit 1
            fi
        fi
    done
done

echo "Finished e2e tests"

exit 0
