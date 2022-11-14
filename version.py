#!/usr/bin/env python3
"""
This is a script to output the new version tag on stdout. It assumes
it is run in the git repo, and the latest commit message text is provided to
the program's STDIN. We search the stdin for a pattern saying whether
to increment version and if so, which one (major, minor, bug).
"""
import subprocess
import re
import sys

tag_reg = r'^v([0-9]+)\.([0-9]+)\.([0-9]+)$'
patterns = {
    'major': (r'^version: major\s*$', 0),
    'minor': (r'^version: minor\s*$', 1),
    'bug': (r'^version: bug\s*$', 2),
}


def tag2version(tag: str):
    """Convert a tag string to a version tuple,
    returns a three-tuple of integers or None if the tag string wasn't
    a proper tag."""
    output = re.match(tag_reg, tag)
    if output is not None:
        return (int(output.group(1)),
                int(output.group(2)),
                int(output.group(3)))
    return None


def find_latest_tag(fetch: bool):
    """
    find the most recent/highest numbered tag. If fetch is True, run
    git fetch origin --tags first so we know what the tags are.
    """
    if fetch:
        subprocess.run(['git', 'fetch', 'origin', '--tags'], check=True)
    tags = [tag2version(i) for i in
            (subprocess.run(['git', 'tag'], stdout=subprocess.PIPE,
                            check=True)
             .stdout.decode().splitlines())]
    latest = sorted([i for i in tags if i])[-1]
    return latest


do_fetch = '--fetch' in sys.argv[1:]
latest = find_latest_tag(do_fetch)

# Now, loop through each line of the commit message, looking for what
# version.
msg = sys.stdin.read().splitlines()
for type, (pattern, index) in patterns.items():
    # if any line matches the version increment pattern, increment the
    # version and exit.
    if any((re.match(pattern, i) for i in msg)):
        new = list(latest)
        new[index] += 1
        for i in range(index+1, len(new)):
            new[i] = 0
        print(f'bumping with a {type} version to: {new}', file=sys.stderr)
        print('v' + '.'.join([str(i) for i in new]))
        break
else:
    print(
        f'Could not find a version message in commit msg: {msg}!\n'
        'Please make sure your commit has version: (bug|minor|major)\n'
        'on a line by itself somewhere. See the README.md Versioning section.',
        file=sys.stderr)
    raise SystemExit(1)
