#!/bin/bash

set -e
repoDir=$1
branchName=$2
message=$3
quantity=$4

cd ${repoDir}

# ignore all uncommit chanages
git stash

echo "==> fetching origin and remote branches..."
git fetch origin
remote_branch=$(git ls-remote --heads origin refs/heads/${branchName})

# create a new branch based on main branch
# or checkout existing branch
if [ "$remote_branch" == "" ]; then
    echo "=> will create branch ${branchName} from main..."
    git checkout -b ${branchName} main
else
    echo "=> branch ${branchName} exists and will be checked-out"
    git checkout ${branchName}
fi

echo "==> creating commits in file..."
for i in $(seq 1 $quantity)
do
    echo "create commit ${i}" >> commit_change.txt
    git add .
    git commit -m "${message} ${i}"
done
echo "==> pushing commits..."
git push -u origin ${branchName}

# output the last created commit id
echo "##output##$(git rev-parse HEAD)"
