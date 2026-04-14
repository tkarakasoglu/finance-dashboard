#!/usr/bin/env sh
# Run `vercel dev` with a Node preload that resolves GOROOT to a real path before
# each `go` subprocess. @vercel/go uses a symlink at .vercel/cache/golang; Go 1.23
# can then fail to read stdlib sources (terminal ENOENT, 502 in the browser).

set -e
ROOT="$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)"
REQ="--require $ROOT/tools/realpath-goroot-spawn.cjs"
case "$NODE_OPTIONS" in
*realpath-goroot-spawn.cjs*) ;;
*)
	if [ -n "$NODE_OPTIONS" ]; then
		NODE_OPTIONS="$NODE_OPTIONS $REQ"
	else
		NODE_OPTIONS="$REQ"
	fi
	export NODE_OPTIONS
	;;
esac

# On macOS, the default @vercel/go dev build flags include `-ldflags -s -w`,
# which can yield a Mach-O missing LC_UUID and crash under dyld. Override the
# flags for local dev to keep the binary runnable.
export GO_BUILD_FLAGS="${GO_BUILD_FLAGS:--ldflags=}"

exec vercel dev "$@"
