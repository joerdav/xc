{ config, pkgs, fetchFromGitHub, ... }:

pkgs.buildGoModule rec {
  pname = "xc";
  version = "v0.8.4";
  subPackages = [ "cmd/xc" ];
  src = pkgs.fetchFromGitHub {
    owner = "joerdav";
    repo = "xc";
    rev = version;
    sha256 = "AttYWpqD8m+CaIaX7WVNlG507uVykIVgArLebRNEUpM=";
  };
  vendorHash = "sha256-EbIuktQ2rExa2DawyCamTrKRC1yXXMleRB8/pcKFY5c=";
}
