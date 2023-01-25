{ config, lib, pkgs, fetchFromGitHub, ... }:

pkgs.buildGoModule rec {
  pname = "xc";
  version = "v0.0.148";
  subPackages = ["cmd/xc"];
  src = pkgs.fetchFromGitHub {
    owner = "joerdav";
    repo = "xc";
    rev = version;
    sha256 = "aWtl/ItO/0hPssfkE9o+DX0iFoHXwo2ouCaEfbtx+Nw=";
  };
  vendorSha256 = "14dtguu787VR8/sYA+9WaS6xr/dB6ZcUjOzDEkFDpH4=";
}
