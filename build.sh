set -euxo pipefail

GOBIN=$(pwd)/functions go install ./...
chmod +x "$(pwd)"/functions/*
go env