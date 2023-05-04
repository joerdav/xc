{ config, lib, pkgs, fetchFromGitHub, ... }:

pkgs.buildGoModule rec {
  pname = "xc";
  version = "v0.4.0";
  subPackages = ["cmd/xc"];
  src = pkgs.fetchFromGitHub {
    owner = "joerdav";
    repo = "xc";
    rev = version;
    sha256 = "pKsttrdXZQnWgJocGtyk7+qze1dpmZTclsUhwun6n8E=";
  };
  vendorSha256 = "hCdIO377LiXFKz0GfCmAADTPfoatk8YWzki7lVP3yLw=";
}
