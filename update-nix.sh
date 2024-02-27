#!/usr/bin/env bash

sed -e "s/__VERSION__/$(git describe --tags --abbrev=0)/g" xc.nix.tmpl > xc.nix

# We first try to build and it fails with hash mismatch, and we use it to populate sha256.
nix-build -E 'with import <nixpkgs> { }; callPackage ./xc.nix { }'
SRC_SHA256="$(nix-build -E 'with import <nixpkgs> { }; callPackage ./xc.nix { }' 2>&1 | grep -oE 'got:\s+sha256-\S+' | cut -d "-" -f 2)"
sed -i -e "s|sha256 = \"\";|sha256 = \"$SRC_SHA256\";|g" xc.nix

# We try again to build and it fails with hash mismatch, and we use it to populate vendorSha256.
VENDOR_SHA256="$(nix-build -E 'with import <nixpkgs> { }; callPackage ./xc.nix { }' 2>&1 | grep -oE 'got:\s+sha256-\S+' | cut -d "-" -f 2)"
sed -i -e "s|vendorHash = \"\";|vendorHash = \"sha256-$VENDOR_SHA256\";|g" xc.nix

