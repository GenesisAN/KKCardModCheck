#!/bin/bash

REPO_NAME=$(basename $(git rev-parse --show-toplevel))

# 获取最新提交的哈希值
LATEST_COMMIT_HASH=$(git rev-parse HEAD)

# 获取最新提交的标签
LATEST_COMMIT_TAG=$(git tag --contains $LATEST_COMMIT_HASH)

# 获取最近的v开头的标签
LATEST_V_TAG=$(git describe --tags --abbrev=0 `git rev-list --tags --max-count=1` 2>/dev/null)

# 需要一起打包的外部文件
FILES_TO_PACKAGE= "" #"../.env ../cors_config.yaml ../email_*.html ../conf/locales/*.yaml"

if [ -z "$LATEST_COMMIT_TAG" ]; then
    # 当最新的提交没有tag时
    if [[ "$LATEST_V_TAG" == v* ]]; then
        # 如果最近的标签以v开头，则取最近的v开头的tag然后在Tag后面加上-dev-最新的提交的短哈希
        VERSION="$LATEST_V_TAG-dev-$(git rev-parse --short $LATEST_COMMIT_HASH)"
    else
        # 如果没有找到v开头的标签，设置默认值
        VERSION="1.0.0"
    fi
else
    # 当最新的提交存在tag时
    if [[ "$LATEST_COMMIT_TAG" == v* ]]; then
        # 并且v开头时，则取Tag的内容为version
        VERSION=$LATEST_COMMIT_TAG
    else
        # 如果标签不是以v开头，设置默认值
        VERSION="1.0.0"
    fi
fi


# 其他变量
GIT_HASH=$(git rev-parse HEAD)
BUILD_DATE=$(date +%Y-%m-%dT%H:%M:%S)
LD_FLAGS="-X main.version=$VERSION -X main.buildDate=$BUILD_DATE -X main.gitHash=$GIT_HASH"

echo "REPO_NAME=$REPO_NAME"
echo "VERSION=$VERSION"
echo "BUILD_DATE=$BUILD_DATE"
echo "GIT_HASH=$GIT_HASH"

echo "Building $REPO_NAME-linux-arm"



go mod tidy

# 编译 ARM、x64 的 Linux 版本
GOOS=linux GOARCH=arm go build -ldflags "$LD_FLAGS" -o "output/$REPO_NAME-linux-arm" .

echo "Building $REPO_NAME-linux-amd64"
GOOS=linux GOARCH=amd64 go build -ldflags "$LD_FLAGS" -o "output/$REPO_NAME-linux-amd64" .

echo "Building $REPO_NAME-windows-amd64.exe"
# 编译 x64 的 Windows 版本并附加 .exe 扩展名
GOOS=windows GOARCH=amd64 go build -ldflags "$LD_FLAGS" -o "output/$REPO_NAME-windows-amd64.exe" .

cd "output"
rm -rf *.tar.gz

echo "Packaging $REPO_NAME-linux-arm.tar.gz"
tar -czf "$REPO_NAME-linux-arm.tar.gz" "$REPO_NAME-linux-arm" $FILES_TO_PACKAGE

echo "Packaging $REPO_NAME-linux-amd64.tar.gz"
tar -czf "$REPO_NAME-linux-amd64.tar.gz" "$REPO_NAME-linux-amd64" $FILES_TO_PACKAGE

echo "Packaging $REPO_NAME-windows-amd64.tar.gz"
tar -czf "$REPO_NAME-windows-amd64.tar.gz" "$REPO_NAME-windows-amd64.exe" $FILES_TO_PACKAGE

echo "Packaging complete!"
