{ config, lib, pkgs, fetchFromGitHub, ... }:

pkgs.buildGoModule rec {
  pname = "xc";
  version = "v0.0.154";
  subPackages = ["cmd/xc"];
  src = pkgs.fetchFromGitHub {
    owner = "joerdav";
    repo = "xc";
    rev = version;
    sha256 = "GJBSPO0PffGdGAHofd1crEFXJi2xqgd8Vk2/g4ff+E4=";
  };
  vendorSha256 = "XDJdCh6P8ScSvxY55ExKgkgFQqmBaM9fMAjAioEQ0+s=";
}
