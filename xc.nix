{ config, lib, pkgs, fetchFromGitHub, ... }:

pkgs.buildGoModule rec {
  pname = "xc";
  version = "v0.8.0";
  subPackages = ["cmd/xc"];
  src = pkgs.fetchFromGitHub {
    owner = "joerdav";
    repo = "xc";
    rev = version;
    sha256 = "vTyCS85xbJnAgbasWD6LFxij9EezzlJ1pyvCJptqmOU=";
  };
  vendorSha256 = "EbIuktQ2rExa2DawyCamTrKRC1yXXMleRB8/pcKFY5c=";
}
