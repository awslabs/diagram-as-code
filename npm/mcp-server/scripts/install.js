#!/usr/bin/env node
'use strict';

/**
 * Postinstall script for the awsdac-mcp-server npm package.
 * Downloads the correct pre-compiled binary from GitHub Releases
 * based on the current OS and CPU architecture.
 */

const https = require('https');
const fs = require('fs');
const path = require('path');
const os = require('os');
const { execSync } = require('child_process');

const pkg = require('../package.json');
const VERSION = `v${pkg.version}`;

const PLATFORM_MAP = {
  linux: 'linux',
  darwin: 'darwin',
  win32: 'windows',
};

const ARCH_MAP = {
  x64: 'amd64',
  ia32: '386',
  arm: 'arm',
  arm64: 'arm64',
};

const platform = PLATFORM_MAP[process.platform];
const arch = ARCH_MAP[process.arch];

if (!platform || !arch) {
  console.warn(
    `awsdac-mcp-server: Unsupported platform "${process.platform}/${process.arch}". ` +
    `Download the binary manually from https://github.com/fernandofatech/diagram-as-code/releases`
  );
  process.exit(0);
}

const isWindows = platform === 'windows';
const binaryName = isWindows ? 'awsdac-mcp-server.exe' : 'awsdac-mcp-server';
const nativeBinaryName = isWindows ? 'awsdac-mcp-server-bin.exe' : 'awsdac-mcp-server-bin';
const zipBaseName = `awsdac-mcp-server-${VERSION}_${platform}-${arch}`;
const zipFileName = `${zipBaseName}.zip`;
const downloadUrl = `https://github.com/fernandofatech/diagram-as-code/releases/download/${VERSION}/${zipFileName}`;

const binDir = path.join(__dirname, '..', 'bin');
const tmpDir = path.join(os.tmpdir(), `awsdac-mcp-install-${Date.now()}`);
const zipPath = path.join(tmpDir, zipFileName);
const nativeBinaryPath = path.join(binDir, nativeBinaryName);

// Skip if binary already exists
if (fs.existsSync(nativeBinaryPath)) {
  console.log(`awsdac-mcp-server: Binary already installed at ${nativeBinaryPath}`);
  process.exit(0);
}

fs.mkdirSync(binDir, { recursive: true });
fs.mkdirSync(tmpDir, { recursive: true });

console.log(`awsdac-mcp-server: Downloading ${downloadUrl} ...`);

function download(url, dest, callback) {
  const file = fs.createWriteStream(dest);

  function get(url) {
    https
      .get(url, (res) => {
        if (res.statusCode === 301 || res.statusCode === 302) {
          get(res.headers.location);
          return;
        }
        if (res.statusCode !== 200) {
          fs.unlink(dest, () => {});
          callback(new Error(`HTTP ${res.statusCode} downloading ${url}`));
          return;
        }
        res.pipe(file);
        file.on('finish', () => file.close(callback));
      })
      .on('error', (err) => {
        fs.unlink(dest, () => {});
        callback(err);
      });
  }

  get(url);
}

function cleanup() {
  try {
    fs.rmSync(tmpDir, { recursive: true, force: true });
  } catch (_) {}
}

download(zipPath, zipPath, (err) => {
  if (err) {
    cleanup();
    console.warn(`awsdac-mcp-server: Failed to download binary — ${err.message}`);
    console.warn('Install manually from: https://github.com/fernandofatech/diagram-as-code/releases');
    process.exit(0); // non-fatal
  }

  console.log(`awsdac-mcp-server: Extracting ${zipFileName} ...`);

  try {
    if (isWindows) {
      execSync(
        `powershell -NoProfile -Command "Expand-Archive -Force '${zipPath}' '${tmpDir}'"`,
        { stdio: 'pipe' }
      );
    } else {
      execSync(`unzip -o "${zipPath}" -d "${tmpDir}"`, { stdio: 'pipe' });
    }

    const extractedBinary = path.join(tmpDir, zipBaseName, binaryName);

    if (!fs.existsSync(extractedBinary)) {
      throw new Error(`Expected binary not found at ${extractedBinary}`);
    }

    fs.copyFileSync(extractedBinary, nativeBinaryPath);

    if (!isWindows) {
      fs.chmodSync(nativeBinaryPath, 0o755);
    }

    console.log(`awsdac-mcp-server: Binary installed at ${nativeBinaryPath}`);
    console.log(`awsdac-mcp-server: Run "awsdac-mcp-server --help" to verify.`);
  } catch (e) {
    console.warn(`awsdac-mcp-server: Failed to extract binary — ${e.message}`);
    console.warn('Install manually from: https://github.com/fernandofatech/diagram-as-code/releases');
  } finally {
    cleanup();
  }
});
