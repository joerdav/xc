{ config, lib, pkgs, fetchFromGitHub, ... }:

pkgs.buildGoModule rec {
  pname = "xc";
  version = "v0.7.0";
  subPackages = ["cmd/xc"];
  src = pkgs.fetchFromGitHub {
    owner = "joerdav";
    repo = "xc";
    rev = version;
    sha256 = "ndaffdU+DYuILZzAwsjLTNWFWbq7CrTcAYBA0j3T3gA=";
  };
  vendorSha256 = "AwlXX79L69dv6wbFtlbHAeZRuOeDy/r6KSiWwjoIgWw=";
}
