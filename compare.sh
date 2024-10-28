#!/bin/bash

if [ "$(git rev-parse HEAD)" = "$(git rev-parse origin/master)" ]; then
    echo "Local and remote are the same."
else
    echo "Local and remote differ."
fi
