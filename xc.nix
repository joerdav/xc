{ config, lib, pkgs, fetchFromGitHub, ... }:

pkgs.buildGoModule rec {
  pname = "xc";
  version = "v0.0.152";
  subPackages = ["cmd/xc"];
  src = pkgs.fetchFromGitHub {
    owner = "joerdav";
    repo = "xc";
    rev = version;
    sha256 = "za+ehaSLD3NHqQhj/OCUMOnVYVrRn40AhyoOeVx//QM=";
  };
  vendorSha256 = "XDJdCh6P8ScSvxY55ExKgkgFQqmBaM9fMAjAioEQ0+s=";
}
