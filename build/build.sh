#!/bin/bash

cd $(dirname "$0")

TMP='/tmp'
LDFLAGS='-s -w'
OUTDIR='../output'
MAINNAME='ghfs'
VERSION="$(git describe --abbrev=0 --tags 2> /dev/null || git rev-parse --abbrev-ref HEAD 2> /dev/null)"
LICENSE='../LICENSE'

mkdir -p "$OUTDIR"

for build in "$@"; do
	arg=($build)
	export GOOS="${arg[0]}"
	export GOARCH="${arg[1]}"
	OS_SUFFIX="${arg[2]}"

	BIN="$TMP"/"$MAINNAME"
	if [ "$GOOS" == 'windows' ]; then
	  BIN="${BIN}.exe"
	fi;
	rm -f "$BIN"
	echo "Building: $GOOS $GOARCH"
	go build -ldflags "$LDFLAGS" -o "$BIN" ../src/main.go

	OUT="$OUTDIR"/"$MAINNAME"-"$VERSION"-"$GOOS""$OS_SUFFIX"-"$GOARCH".zip
	zip -j "$OUT" "$BIN" "$LICENSE"
done