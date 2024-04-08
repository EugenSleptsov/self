#!/bin/bash

# Attempt to start the Go program
go run .

# Check the exit code of the program
if [ $? -ne 0 ]; then
    # If program failed to start, reset to the last commit
    git reset --hard HEAD^
else
    # If program started successfully, commit changes and push to the server
    git add .
    git commit -m "Next iteration"
    git push origin master
fi
