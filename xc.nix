{ config, lib, pkgs, fetchFromGitHub, ... }:

pkgs.buildGoModule rec {
  pname = "xc";
  version = "v0.3.0";
  subPackages = ["cmd/xc"];
  src = pkgs.fetchFromGitHub {
    owner = "joerdav";
    repo = "xc";
    rev = version;
    sha256 = "e/cJ1kVFUs2chOpqu2tHsK9MBhlMx9u6mIa5uUwvVg8=";
  };
  vendorSha256 = "hCdIO377LiXFKz0GfCmAADTPfoatk8YWzki7lVP3yLw=";
}
