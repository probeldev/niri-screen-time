{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-24.11";
    flake-utils.url = "github:numtide/flake-utils";
  };
  outputs = {
    self,
    nixpkgs,
    flake-utils,
  }:
    flake-utils.lib.eachDefaultSystem (system: let
      pkgs = import nixpkgs {
        inherit system;
      };
      niri-screen-time-package = pkgs.callPackage ./package.nix {};
    in {
      packages = rec {
        niri-screen-time = niri-screen-time-package;
        default = niri-screen-time;
      };

      apps = rec {
        niri-screen-time = flake-utils.lib.mkApp {
          drv = self.packages.${system}.niri-screen-time;
        };
        default = niri-screen-time;
      };

      devShells.default = pkgs.mkShell {
        packages = (with pkgs; [
          go
        ]);
      };
    });
}
