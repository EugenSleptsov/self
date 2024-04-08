#!/bin/bash

# Each time we go run . we are receiving new main.go content

# Running current version which we know is working
go run .

# After run we have received new main.go content
# making a local commit to save the current state
git add .
git commit -m "Next iteration of the project, time: $(date +'%Y-%m-%d %H:%M:%S')"

# Running the new version
go run .

# Checking if the new version is working
if [ $? -ne 0 ]; then
  # If the new version is not working, we are reverting to the previous version (discarding the last commit)
  git reset --hard HEAD~1
  git clean -fd
else
  # If the new version is working, we are pushing the changes to the remote repository, but we need to return to the saved local commit
  git reset --hard
  git push
fi
