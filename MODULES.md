# Working with golang modules
## [Intranet article link](https://intranet.untangle.com/display/MF/Working+with+Golang+Modules)
## Adding golang vendor modules
1. ### Project must have the go.mod created
    ```
    go mod init
    ```
2. ### Use go get to add the package
    ```
    go get -d github.com/untangle/golang-shared
    ```
3. ### Use go tidy to cleanup dependencies
    ```
    go mod tidy
    ```
4. ### Use go vendor to add the vendor files into vendor/*
    ```
    go mod vendor
    ```
## Updating golang vendor modules
1. ### Go get with the -u flag will update to latest tag
    ```
    go get -u github.com/untangle/golang-shared
    ```
    ### If you are relying on a specific commit, you can use the hash also
    ```
    go get -u github.com/untangle/golang-shared@19fa40e
    ```
2. ### Then we need to tidy and vendor again
    ```
    go mod tidy
    go mod vendor
    ```