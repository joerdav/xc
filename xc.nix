{ config, lib, pkgs, fetchFromGitHub, ... }:

pkgs.buildGoModule rec {
  pname = "xc";
  version = "v0.2.0";
  subPackages = ["cmd/xc"];
  src = pkgs.fetchFromGitHub {
    owner = "joerdav";
    repo = "xc";
    rev = version;
    sha256 = "ACOi9MdsHZ7Td8pcYGj6Oy7uE/g/Mx7B+Uqw6K9hIrE=";
  };
  vendorSha256 = "YETiKG8/+DFd9paqy58YYzCKZ47cWlhnUWHjrw5CrDI=";
}
