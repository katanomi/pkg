#!/bin/bash

repoDir=$1
branchName=$2
message=$3

cd ${repoDir}

# ignore all uncommit chanages
git stash

echo "create commit" >> commit_change.txt

# create a new branch based on main branch
git checkout -b ${branchName} main
git push -u origin ${branchName}
git add .
git commit -m "${message}"
git push -f
