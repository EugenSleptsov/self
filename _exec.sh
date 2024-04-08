#!/bin/bash

go run .

if [ $? -ne 0 ]; then
    git reset --hard
    git clean -fd
else
    git add .
    git commit -m "Next iteration of the project, time: $(date +'%Y-%m-%d %H:%M:%S')"
    git push
fi
