#!/bin/bash

# Define the regular expressions for log statements
LOG_REGEX1='[lL]ogger.(Alert|Err|Notice|Info|Debug|Trace|Fatal|Warning|OCTrace|OCWarn|OCDebug|OCErr|OCCrit)\("([^"]*)"\)'
LOG_REGEX2='[lL]ogger.(Alert|Err|Notice|Info|Debug|Trace|Fatal|Warning|OCTrace|OCWarn|OCDebug|OCErr|OCCrit)\("([^"]*)",\s*([^)]*)\)'

# Check if the current directory is a Git repository
if [ -d .git ]; then

   # Configure a global Git setting for the safe directory
   git config --global --add safe.directory /go/untangle-shared
   echo ".git directory exists in the current directory."

   # Get the list of added/modified files in the most recent commit
   modified_files=$(git status -uno | grep "modified:" | awk '{print $2}')

   # Iterate over the list of modified files
   for file in $modified_files; do
      # Use grep with a negative lookahead assertion to exclude lines ending with \n
      modified_lines=$(grep -HnE "$LOG_REGEX1|$LOG_REGEX2" "$file" | grep -vE '.*\\n.*')

      # Check if the modified_lines variable is empty or consists of only whitespace characters
      if [[ "$(printf '%s' "$modified_lines")" =~ ^[[:space:]]*$ ]]; then
         continue  # Skip to the next file if there are no relevant log statements
      fi

      # Check if there are any log statements that don't end with \n
      if [ -n "$modified_lines" ]; then
         echo -e '\nLog statement doesn'"'"'t end with \\n\n'"$modified_lines"
      fi
   done
else
    echo ".git directory does not exist in the current directory."
fi
