{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.11";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          system = system;
        };
        xc = pkgs.callPackage ./xc.nix { };
      in
      {
        defaultPackage = xc;
        packages = {
          default = xc;
          xc = xc;
        };
        devShells = {
          default = pkgs.mkShell {
            packages = [
              pkgs.go_1_25
              pkgs.golangci-lint
            ];
          };
        };
      }
    );
}
