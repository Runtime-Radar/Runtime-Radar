#!/usr/bin/env bash

set -e

# Go to the project root directory
cd $(dirname ${0})/../..

publishPackage() {

    packageName="RuntimeRadar.UI"

    buildVersion=$(cat package.json \
          | grep version \
          | head -1 \
          | awk -F: '{ print $2 }' \
          | sed 's/[",]//g' \
          | tr -d '[[:space:]]')

    branchName=$(git branch | sed -n '/\* /s///p' | awk -F'/' '{print $2}')


    commitSha=$(git rev-parse --short HEAD)
    commitMessage=$(git log --oneline -n 1)

    buildVersionName="${buildVersion}-${commitSha}"

    echo "Starting publish process of ${packageName} for ${buildVersionName} into ${branchName}.."
    echo "with commitMessage ${commitMessage}"

    packageFullName="${packageName}.${buildVersion}"

    echo "Package fullname ${packageFullName}"
    tar -czf ${packageFullName}.tar.gz dist/apps/runtime-radar

    # its workaround for crossbuilder
    echo "${packageFullName}.tar.gz" > packageExportName
}

(publishPackage)
