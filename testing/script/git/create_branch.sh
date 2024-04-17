#!/bin/bash

set -e
repoDir=$1
branchName=$2

cd ${repoDir}

# ignore all uncommit chanages
git stash

# create a new branch based on main branch
git checkout -b ${branchName} main
git push -u origin ${branchName}
