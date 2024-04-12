#!/bin/bash

set -e
repoDir=$1
branchName=$2
message=$3

cd ${repoDir}

# ignore all uncommit changes
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

echo "==> creating commit in file..."
echo -n "create commit" >> commit_change.txt
git add .
git commit -m "${message}"
echo "==> pushing commit..."
git push -u origin ${branchName} -f

# output the newly created commit id
echo "##output##$(git rev-parse HEAD)"
