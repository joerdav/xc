{ config, pkgs, fetchFromGitHub, ... }:

pkgs.buildGoModule rec {
  pname = "xc";
  version = "v0.8.6";
  subPackages = [ "cmd/xc" ];
  src = pkgs.fetchFromGitHub {
    owner = "joerdav";
    repo = "xc";
    rev = version;
    sha256 = "Q17ldwHp1Wp/u0BkUZiA1pRJaFpo/5iDW011k9qkIEA=";
  };
  env = {
    CGO_ENABLED = "0";
  };
  vendorHash = "sha256-EbIuktQ2rExa2DawyCamTrKRC1yXXMleRB8/pcKFY5c=";
}
