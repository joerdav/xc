{ config, lib, pkgs, fetchFromGitHub, ... }:

pkgs.buildGoModule rec {
  pname = "xc";
  version = "v0.0.175";
  subPackages = ["cmd/xc"];
  src = pkgs.fetchFromGitHub {
    owner = "joerdav";
    repo = "xc";
    rev = version;
    sha256 = "Uc9MTxl32xQ7u6N0mocDAoD9tgv/YOPCzhonsavX9Vo=";
  };
  vendorSha256 = "cySflcTuAzbFZbtXmzZ98nfY8HUq1UedONTtKP4EICs=";
}
