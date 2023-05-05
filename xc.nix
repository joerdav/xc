{ config, lib, pkgs, fetchFromGitHub, ... }:

pkgs.buildGoModule rec {
  pname = "xc";
  version = "v0.4.1";
  subPackages = ["cmd/xc"];
  src = pkgs.fetchFromGitHub {
    owner = "joerdav";
    repo = "xc";
    rev = version;
    sha256 = "Dc7MVn9hF2HtXqMvWQ5UsLQW5ZKcFKt7AHcXdiWDs1I=";
  };
  vendorSha256 = "hCdIO377LiXFKz0GfCmAADTPfoatk8YWzki7lVP3yLw=";
}
