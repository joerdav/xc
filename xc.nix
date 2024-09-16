{ config, pkgs, fetchFromGitHub, ... }:

pkgs.buildGoModule rec {
  pname = "xc";
  version = "v0.8.5";
  subPackages = [ "cmd/xc" ];
  src = pkgs.fetchFromGitHub {
    owner = "joerdav";
    repo = "xc";
    rev = version;
    sha256 = "eaFHK7VsfLSgSJehv4urxq8qMPT+zzs2tRypz4q+MLc=";
  };
  vendorHash = "sha256-EbIuktQ2rExa2DawyCamTrKRC1yXXMleRB8/pcKFY5c=";
}
