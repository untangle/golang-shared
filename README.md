This package contains shared service components between untangle's restd, packetd, and reportd packages

# Repo Layout:
* ## services/*
  * ### These are services that are imported into daemon projects
* ## structs/*
  * ### These are common struct types that may need to be used across repositories
* ## protobuffersrc/*
  * ### These are the .proto source files that are compiled into *.pb.go files
  * ### To rebuild these, use the included docker files:
    ```
    docker-compose -f build/docker-compose.build.yml up --build musl-local
    ```

# Versioning and CI/CD

Every time you merge to master, a new version must be created. This is
done semi-automatically by CI, but you need to tell it which type of
version -- bug, minor, major -- is being created by the merge. To
indicate this, when merging to master, add a message to the bottom
line of your merge message (which is created when you merge the PR
from the github UI) saying `version: bug` or `version: minor` or
`version: major`.

The version message needs to be on a line by itself. *Make sure it's
in your merge message.* It cannot be in some message that was commited
to the branch that is being merged, it must be in the merge message
itself, which will be created when you merge the PR on github.

When things go wrong and it didn't version correctly, you can push to
master with an empty commit like:

```
git commit --allow-empty -m "version: bug"
```

If the code is being merged on a **different branch** than master,
there is an option to provide the branch name using the **branch** parameter.
`version.py --branch {branch_name}`
The final tag that will be applied will contain the latest tag number available on the branch provided. 
If no `branch` parameter is provided, latest tag will default from master.

Example: 

Latest tag on master: `v.1.60.1`
Latest tag on eft: `v.1.59.0`
Merge message:
```
version: bug
```
Branch input: `eft`
Tag that will be applied: `v1.59.1`

The script can be tested by running the following command:
```
echo -e "version: minor" | python3 version.py --fetch -- branch eft
```

## Versioning Strategy

### Major Version Update
- Increment when a new release branch is created.
- Master branch uses the latest major version.
- Release branch uses the previous major version.

### Minor Version Increment
- Increment for new features or functionality.

### Patch ('Bug') Version Increment
- Increment for bug fixes.

##  [Working with modules](./MODULES.md)
