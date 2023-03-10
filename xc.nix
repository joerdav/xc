{ config, lib, pkgs, fetchFromGitHub, ... }:

pkgs.buildGoModule rec {
  pname = "xc";
  version = "v0.1.181";
  subPackages = ["cmd/xc"];
  src = pkgs.fetchFromGitHub {
    owner = "joerdav";
    repo = "xc";
    rev = version;
    sha256 = "C6qZdO6+n9BWm69y09kvnEBF45sB6bfOfmteNO2x68I=";
  };
  vendorSha256 = "cySflcTuAzbFZbtXmzZ98nfY8HUq1UedONTtKP4EICs=";
}
