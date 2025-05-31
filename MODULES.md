# Working with golang modules
## [Intranet article link](https://awakesecurity.atlassian.net/wiki/spaces/MF/pages/2076150225/Working+with+Golang+Modules)
## Adding golang vendor modules
1. ### Project must have the go.mod created
    ```
    go mod init
    ```
2. ### Use go get to add the package
    ```
    GOPRIVATE=github.com/untangle/golang-shared go get -d github.com/untangle/golang-shared
    ```
3. ### Use go tidy to cleanup dependencies
    ```
    go mod tidy
    ```
4. ### Use go vendor to add the vendor files into vendor/*
    ```
    go mod vendor
    ```
