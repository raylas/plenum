{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachSystem [
      "x86_64-linux"
      "aarch64-linux"
      "aarch64-darwin"
    ] (system:
      let
        graft = pkgs: pkg:
          pkg.override { buildGoModule = pkgs.buildGo122Module; };
        pkgs = import nixpkgs {
          inherit system;
          overlays = [
            (final: prev: {
              go = prev.go_1_22;
              go-tools = graft prev prev.go-tools;
              gotools = graft prev prev.gotools;
              gopls = graft prev prev.gopls;
            })
          ];
        };
      in
      rec {
        devShell = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            go-tools
            gotools
            gopls
          ];
        };
      });
}
