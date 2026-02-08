#!/bin/bash

# Get the project root directory
ROOT_DIR=$(git rev-parse --show-toplevel)

echo "Setting up Git hooks..."

# Ensure scripts have execution permissions
chmod +x "$ROOT_DIR/scripts/git/pre-commit"
chmod +x "$ROOT_DIR/scripts/git/commit-msg"

# Create symbolic links to the .git/hooks directory
ln -sf "$ROOT_DIR/scripts/git/pre-commit" "$ROOT_DIR/.git/hooks/pre-commit"
ln -sf "$ROOT_DIR/scripts/git/commit-msg" "$ROOT_DIR/.git/hooks/commit-msg"

echo "Git hooks installed successfully!"
