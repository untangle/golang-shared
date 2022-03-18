# Working with golang modules
## [Intranet article link](https://intranet.untangle.com/display/MF/Working+with+Golang+Modules)
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
## Updating golang vendor modules
1. ### Verify the commit in golang-shared is tagged with a version
    ```
    cd golang-shared
    git tag vX.X.X
    git push
    git push --tags
    ```
2. ### In the dependent package (ie: packetd, reportd, etc) use Go get with the -u flag and include the version hash until the PR is approved
    ```
    cd packetd
    GOPRIVATE=github.com/untangle/golang-shared go get -u github.com/untangle/golang-shared@9eb30a3d32ae74ac7e925824a5601817e895067b
    ```

3. ### Verify the version in go.mod has been updated
    ```
    grep golang-shared go.mod
        github.com/untangle/golang-shared v0.2.1
    ```
4. ### Then we need to tidy and vendor again
    ```
    go mod vendor
    go mod tidy
    ```