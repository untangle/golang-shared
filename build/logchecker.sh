#!/bin/bash

# Define the regular expressions for log statements
LOG_REGEX1='[lL]ogger.(Alert|Err|Notice|Info|Debug|Trace|Fatal|Warning|OCTrace|OCWarn|OCDebug|OCErr|OCCrit)\("([^"]*)"\)'
LOG_REGEX2='[lL]ogger.(Alert|Err|Notice|Info|Debug|Trace|Fatal|Warning|OCTrace|OCWarn|OCDebug|OCErr|OCCrit)\("([^"]*)",\s*([^)]*)\)'

if [ -d .git ]; then

   git config --global --add safe.directory /go/untangle-shared
   echo ".git directory exists in the current directory."

   # Get the list of added/modified lines in the most recent commit for the specified directory
   modified_lines=$(git diff --staged --unified=0 | grep -E '^\+' | grep -E "$LOG_REGEX1|$LOG_REGEX2")
   echo "Modified lines: $modified_lines"

    echo $modified_lines | while IFS= read -r line; do
      # # Check if the line is empty or consists of only whitespace characters
      if [[ "$(printf '%s' "$line")" =~ ^[[:space:]]*$ ]]; then
         continue
      fi

      # Check if the line ends with \n
      if [[ ! "$(printf '%s' "$line")" == *$'\n'* ]]; then
         echo "Log statement doesn't end with '\\n': $line"

      fi
   done
else
    echo ".git directory does not exist in the current directory."
fi