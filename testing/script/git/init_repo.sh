#!/bin/bash

set -e
repoHttpCloneUrl=$1
username=$2
token=$3

# create repo dir
temp_dir=$(mktemp -d);cd ${temp_dir}
git init

# setup git credential
protocol=$(echo "${repoHttpCloneUrl}" | awk -F: '{print $1}')
host=$(echo "${repoHttpCloneUrl}" | awk -F/ '{print $3}')
echo "${protocol}://${username}:${token}@${host}" > .git/git_credential

# setup git config
git config user.name "e2e-tester"
git config user.email "e2e-tester@katanomi.dev"
git config credential.helper "store --file "$(pwd)"/.git/git_credential"
git remote add origin ${repoHttpCloneUrl}

# create some files
mkdir -p a/b/c
# WARNING: Do not change the file path and contents, it will break the test
echo -n "a/b/c/hello.txt" > a/b/c/hello.txt
echo -n "README.md" > README.md

# make a commit
git add .
git commit -am "chore: initial commit"
git branch -m "main"
git push -u origin "main"

# output the repo dir
echo "##output##$(pwd)"
