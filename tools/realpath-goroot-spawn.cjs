'use strict';

if (process.env.VERCEL_GOROOT_REALPATH_DEBUG) {
  process.stderr.write('[realpath-goroot-spawn] preload loaded\n');
}

/**
 * Preload for `vercel dev` with @vercel/go.
 *
 * The Vercel Go builder sets GOROOT to project/.vercel/cache/golang, which is
 * often a symlink into ~/Library/Caches/... It also sets GOMODCACHE/GOCACHE
 * under that tree. Go 1.23 can then fail to open stdlib sources (ENOENT) even
 * though the files exist. Resolving GOROOT to a real path before spawning the
 * go binary avoids that.
 */

const cp = require('child_process');
const fs = require('fs');
const path = require('path');

function realpathGOROOT(env) {
  if (!env || typeof env.GOROOT !== 'string' || env.GOROOT === '') {
    return env;
  }
  try {
    const resolved = fs.realpathSync(env.GOROOT);
    if (resolved !== env.GOROOT) {
      return { ...env, GOROOT: resolved };
    }
  } catch (_) {
    /* keep original */
  }
  return env;
}

function isGoCommand(command) {
  if (typeof command !== 'string') {
    return false;
  }
  const base = path.basename(command, path.extname(command));
  return base === 'go';
}

function patchSpawn(orig) {
  return function patchedSpawn(command, args, options) {
    if (options && options.env && isGoCommand(command)) {
      if (process.env.VERCEL_GOROOT_REALPATH_DEBUG) {
        const e = options.env;
        process.stderr.write(
          `[realpath-goroot-spawn] go spawn GOROOT before=${e.GOROOT}\n`
        );
      }
      options = { ...options, env: realpathGOROOT(options.env) };
      if (process.env.VERCEL_GOROOT_REALPATH_DEBUG) {
        process.stderr.write(
          `[realpath-goroot-spawn] go spawn GOROOT after=${options.env.GOROOT}\n`
        );
      }
    }
    return orig.call(cp, command, args, options);
  };
}

cp.spawn = patchSpawn(cp.spawn);
cp.spawnSync = patchSpawn(cp.spawnSync);
