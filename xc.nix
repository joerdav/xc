{ config, lib, pkgs, fetchFromGitHub, ... }:

pkgs.buildGoModule rec {
  pname = "xc";
  version = "v0.0.169";
  subPackages = ["cmd/xc"];
  src = pkgs.fetchFromGitHub {
    owner = "joerdav";
    repo = "xc";
    rev = version;
    sha256 = "8iy9003fU3lY7JY2cri1OCzky4muMsVUfRof7LHyzy4=";
  };
  vendorSha256 = "cySflcTuAzbFZbtXmzZ98nfY8HUq1UedONTtKP4EICs=";
}
