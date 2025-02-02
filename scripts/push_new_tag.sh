#!/bin/bash

# Check if we're in the scripts directory, if not try to enter it
if [[ ! $(basename "$PWD") == "scripts" ]]; then
    if [[ -d "scripts" ]]; then
        cd scripts || exit 1
    else
        echo "Error: Must be executed from the root directory or scripts directory"
        exit 1
    fi
fi

# Check for uncommitted changes
if [[ -n $(git status -s) ]]; then
    echo "Error: There are uncommitted changes. Please commit or stash them first."
    exit 1
fi

# Get the latest tag
latest_tag=$(git describe --tags --abbrev=0 2>/dev/null)

if [[ -z "$latest_tag" ]]; then
    echo "No existing tags found. Creating v0.0.1"
    new_tag="v0.0.1"
else
    echo "Latest tag: $latest_tag"

    if [[ "$1" == "-f" ]]; then
        # Force update the previous tag instead of creating a new one
        git tag -fa "$latest_tag" -m "Release $latest_tag"
        git push --force origin "$latest_tag"
        echo "Successfully force updated tag: $latest_tag"
        exit 0
    fi
    
    # Extract version numbers
    IFS='.' read -r major minor patch <<< "${latest_tag#v}"
    
    # Ensure we're reading the numbers as integers
    major=${major:-0}
    minor=${minor:-0}
    patch=${patch:-0}
    
    # Increment patch version
    new_patch=$((patch + 1))
    new_tag="v${major}.${minor}.${new_patch}"
    
    echo "Creating new tag: $new_tag"
fi

# Create and push the new tag
git tag -a "$new_tag" -m "Release $new_tag"
git push origin "$new_tag"
echo "Successfully created and pushed tag: $new_tag"
