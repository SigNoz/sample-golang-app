#!/usr/bin/env bash
set -e

# OpenTelemetry / SigNoz 相关环境变量（可按需覆盖）
export SERVICE_NAME="${SERVICE_NAME:goapp}"
export OTEL_EXPORTER_OTLP_ENDPOINT="${OTEL_EXPORTER_OTLP_ENDPOINT:-localhost:4317}"
export INSECURE_MODE="${INSECURE_MODE:-true}"

# 可选：OTEL 服务名（与 SERVICE_NAME 一致即可）
export OTEL_SERVICE_NAME="${OTEL_SERVICE_NAME:-$SERVICE_NAME}"

cd "$(dirname "$0")"
#go run .

go build -v -gcflags 'all=-N -l' -o sample-golang-app
./sample-golang-app

