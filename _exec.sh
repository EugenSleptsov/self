#!/bin/bash

go run .

if [ $? -ne 0 ]; then
    git reset --hard
    git clean -fd
else
    git add .
    git commit -m "Next iteration"
    git push
fi
