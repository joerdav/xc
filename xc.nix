{ config, lib, pkgs, fetchFromGitHub, ... }:

pkgs.buildGoModule rec {
  pname = "xc";
  version = "v0.5.0";
  subPackages = ["cmd/xc"];
  src = pkgs.fetchFromGitHub {
    owner = "joerdav";
    repo = "xc";
    rev = version;
    sha256 = "cVTa2ot95Hcm+1V1QXnlxSL9OjmoQNR9nVUgW/rZhl0=";
  };
  vendorSha256 = "J4/a4ujM7A6bDwRlLCYt/PmJf6HZUmdYcJMux/3KyUI=";
}
