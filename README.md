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

# Versioning

When merging to master, add a message to the bottom line of your merge
message saying `version: bug` or `version: minor` or `version: major`.

The version message needs to be on a line by itself. Make sure it's
your merge message.

##  [Working with modules](./MODULES.md)
