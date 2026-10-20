#!/usr/bin/env bash

set -e

if command -v make >/dev/null 2>&1; then
	exec make up
fi

if command -v podman >/dev/null 2>&1 && podman compose version >/dev/null 2>&1; then
	exec podman compose up -d --build
elif command -v docker >/dev/null 2>&1 && docker compose version >/dev/null 2>&1; then
	exec docker compose up -d --build
else
	echo "Erro: make, podman compose ou docker compose não encontrado." >&2
	exit 1
fi
