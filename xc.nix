{ config, lib, pkgs, fetchFromGitHub, ... }:

pkgs.buildGoModule rec {
  pname = "xc";
  version = "v0.0.159";
  subPackages = ["cmd/xc"];
  src = pkgs.fetchFromGitHub {
    owner = "joerdav";
    repo = "xc";
    rev = version;
    sha256 = "5Vw/UStMtP5CHbSCOzeD4LMJccPG34Rxw9VHc9Ut3oM=";
  };
  vendorSha256 = "XDJdCh6P8ScSvxY55ExKgkgFQqmBaM9fMAjAioEQ0+s=";
}
