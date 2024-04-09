#!/bin/bash
git reset --hard HEAD
git pull

# Get the directory of the script
SCRIPT_DIR=$(dirname "$(readlink -f "$0")")

# Change to the script directory
cd "$SCRIPT_DIR" || exit

# Each time we go run . we are receiving new main.go content

# Running current version which we know is working
go run .

# After run we have received new main.go content
# making a local commit to save the current state
git add .
git commit -m "Next iteration of the project, time: $(date +'%Y-%m-%d %H:%M:%S')"

# Running the new version
go run .

# Running the version provided by new version
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
