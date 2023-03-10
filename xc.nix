{ config, lib, pkgs, fetchFromGitHub, ... }:

pkgs.buildGoModule rec {
  pname = "xc";
  version = "v0.1.180";
  subPackages = ["cmd/xc"];
  src = pkgs.fetchFromGitHub {
    owner = "joerdav";
    repo = "xc";
    rev = version;
    sha256 = "nRqCNFfQAgI0WXfcX9JRSzRPnxRFjFRCwGfzqQPKw8g=";
  };
  vendorSha256 = "cySflcTuAzbFZbtXmzZ98nfY8HUq1UedONTtKP4EICs=";
}
