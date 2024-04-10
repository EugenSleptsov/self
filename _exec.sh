#!/bin/bash
# Get the directory of the script
SCRIPT_DIR=$(dirname "$(readlink -f "$0")")

# Change to the script directory
cd "$SCRIPT_DIR" || exit

# In case of some problems with current git status, we are resetting the repository
git reset --hard HEAD
git pull

# Each time we go run . we are receiving new main.go content

# Running current version which we know is working
go run .

# After run we have received new main.go content
# making a local commit to save the current state
git add .
git commit -m "Next iteration of the project, time: $(date +'%Y-%m-%d %H:%M:%S')"

# Running the new version
go run .

# Running the version provided by new version (sometimes code changes model to davinci codex and it is not working)
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
