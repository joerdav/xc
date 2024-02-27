{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/23.11";
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
          xc = xc;
        };
        devShells = {
          default = pkgs.mkShell {
            packages = [ xc ];
          };
          xc = pkgs.mkShell {
            packages = [ xc ];
          };
        };
      }
    );
}
