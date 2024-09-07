#!/bin/bash

FILE_TO_WATCH=$1
TEMP_FILE=$(mktemp)

SORTED1=$(mktemp)
SORTED2=$(mktemp)

# Create an initial snapshot of the file
cp "$FILE_TO_WATCH" "$TEMP_FILE"

cleanup() {
    exit 0
}

# Set a trap to capture SIGINT (CTRL-C) and run the cleanup function
trap cleanup SIGINT

# Tail the file and detect changes
while true; do
    sleep 1  # Check every 1 second

    if ! cmp -s "$FILE_TO_WATCH" "$TEMP_FILE"; then

        sort "$TEMP_FILE" > "$SORTED1"
        sort "$FILE_TO_WATCH" > "$SORTED2"

        # Show the changes using diff
        printf "\nChanges detected:"
        diff --color "$SORTED1" "$SORTED2"

        # Update the snapshot
        cp "$FILE_TO_WATCH" "$TEMP_FILE"

    fi
done
