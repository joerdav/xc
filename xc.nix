{ config, lib, pkgs, fetchFromGitHub, ... }:

pkgs.buildGoModule rec {
  pname = "xc";
  version = "v0.0.110";
  subPackages = ["cmd/xc"];
  src = pkgs.fetchFromGitHub {
    owner = "joerdav";
    repo = "xc";
    rev = version;
    sha256 = "Yn+9FPYPwx1BfDf45uaNZ5fLEIu6LgN6ihw4eDAT9mY=";
  };
  vendorSha256 = "14dtguu787VR8/sYA+9WaS6xr/dB6ZcUjOzDEkFDpH4=";
}
