{ config, lib, pkgs, fetchFromGitHub, ... }:

pkgs.buildGoModule rec {
  pname = "xc";
  version = "v0.0.180";
  subPackages = ["cmd/xc"];
  src = pkgs.fetchFromGitHub {
    owner = "joerdav";
    repo = "xc";
    rev = version;
    sha256 = "LjcN6cxRehhNvJ7yi80FTF/jblB7+24nTfgLawSs0DE=";
  };
  vendorSha256 = "cySflcTuAzbFZbtXmzZ98nfY8HUq1UedONTtKP4EICs=";
}
