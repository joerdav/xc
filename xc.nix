{ config, lib, pkgs, fetchFromGitHub, ... }:

pkgs.buildGoModule rec {
  pname = "xc";
  version = "v0.6.0";
  subPackages = ["cmd/xc"];
  src = pkgs.fetchFromGitHub {
    owner = "joerdav";
    repo = "xc";
    rev = version;
    sha256 = "0Er8MqAqKCyz928bdbYRO3D9sGZ/JJBrCXhlq9M2dEA=";
  };
  vendorSha256 = "J4/a4ujM7A6bDwRlLCYt/PmJf6HZUmdYcJMux/3KyUI=";
}
