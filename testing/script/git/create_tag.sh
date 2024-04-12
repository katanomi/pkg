#!/bin/bash

set -e
repoDir=$1
branchName=$2
message=$3
tag=$4

cd ${repoDir}

# ignore all uncommit chanages
git stash


# create a new branch based on main branch
if [ "${branchName}" != "" ]; then
    branchName="main"
fi

git tag -a ${tag} -m "${message}"
git push origin ${tag}
