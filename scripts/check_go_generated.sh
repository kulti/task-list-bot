#!/bin/bash

set -e

go generate ./...
DIFF_FILES=$(git ls-files --modified --other)
if [ -n "${DIFF_FILES}" ]; then
    echo "Modified and new files:"
    echo "${DIFF_FILES}"
    echo ""
    echo "Git Diff:"
    git diff
    exit 1
fi
