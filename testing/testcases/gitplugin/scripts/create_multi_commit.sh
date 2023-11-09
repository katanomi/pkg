#!/bin/bash

repoDir=$1
branchName=$2
message=$3
quantity=$4

cd ${repoDir}

# ignore all uncommit chanages
git stash
# create a new branch based on main branch
git checkout -b ${branchName} main

for i in $(seq 1 $quantity)
do
    echo "create commit ${i}" >> commit_change.txt
    git add .
    git commit -m "${message} ${i}"
done

git push -u origin ${branchName}

