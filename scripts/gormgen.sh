#!/bin/bash

shellExit()
{
if [ $1 -eq 1 ]; then
    printf "\nfailed!!!\n\n"
    exit 1
fi
}

printf "\nRegenerating file\n\n"
time go run -v ./cmd/mysqlmd/main.go  -addr $1 -user $2 -pass $3 -name $4 -tables $5
shellExit $?

printf "\ncreate curd code : \n"
time go build -o gormgen ./cmd/gormgen/main.go
shellExit $?

if [ ! -d $GOPATH/bin ];then
   mkdir -p $GOPATH/bin
fi

mv gormgen $GOPATH/bin
shellExit $?

time go install golang.org/x/tools/cmd/goimports@latest
#如果tables != all 则只生成指定的表
if [ $5 != "all" ];then
    printf "\nGenerating code for tables: $5\n\n"
    # 根据，分割
    IFS=',' read -r -a tables <<< "$5"
    for table in "${tables[@]}"
    do
        printf "\nGenerating code for table: $table\n\n"
        result="${table#t_}"  # 去掉前缀
        result="${result%_tab}"  # 再去掉后缀
        printf "\nGenerating code for table without prefix and suffix: $result\n\n"
        go generate ./internal/repository/mysql/$result/...
        goimports -w ./internal/repository/mysql/$result
        shellExit $?
    done
else
    printf "\nGenerating code for all tables\n\n"
    go generate ./...
    goimports -w .
    shellExit $?
fi

printf "\nFormatting code\n\n"
time go run -v ./cmd/mfmt/main.go
shellExit $?

printf "\nDone.\n\n"
