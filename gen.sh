#!/usr/bin/env bash
set -euo pipefail

# Generate *.pb.go and *_grpc.pb.go next to each .proto source.
# Plain protoc only. No buf.
# The module option strips the module prefix from each go_package so the
# output path stays next to the source file.

module="github.com/prothegee/erposfnb-protobuf"
root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

proto_files=(
    "common/common.proto"
    "account/account.proto"
    "tenant/tenant.proto"
    "pos/pos.proto"
    "media/media.proto"
    "warehouse/warehouse.proto"
    "gateway_stream/stream.proto"
)

cd "$root_dir"

protoc \
    -I . \
    --go_out=. \
    --go_opt=module="$module" \
    --go-grpc_out=. \
    --go-grpc_opt=module="$module" \
    "${proto_files[@]}"
