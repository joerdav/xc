{ config, pkgs, fetchFromGitHub, ... }:

pkgs.buildGoModule rec {
  pname = "xc";
  version = "v1";
  subPackages = [ "cmd/xc" ];
  src = pkgs.fetchFromGitHub {
    owner = "joerdav";
    repo = "xc";
    rev = version;
    sha256 = "p6rXzNiDm2uBvO1MLzC5pJp/0zRNzj/snBzZI0ce62s=";
  };
  vendorHash = "sha256-EbIuktQ2rExa2DawyCamTrKRC1yXXMleRB8/pcKFY5c=";
}
