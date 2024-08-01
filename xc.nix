{ config, pkgs, fetchFromGitHub, ... }:

pkgs.buildGoModule rec {
  pname = "xc";
  version = "v0.8.2";
  subPackages = [ "cmd/xc" ];
  src = pkgs.fetchFromGitHub {
    owner = "joerdav";
    repo = "xc";
    rev = version;
    sha256 = "JHu/+Uey9iR0KkFvm75yQebO42je5lEOTXk4VCOO0gg=";
  };
  vendorHash = "sha256-EbIuktQ2rExa2DawyCamTrKRC1yXXMleRB8/pcKFY5c=";
}
