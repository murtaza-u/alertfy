{
  description = "Webhook to hook alertmanager to ntfy";
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-24.05";
    flake-utils.url = "github:numtide/flake-utils";
  };
  outputs = { nixpkgs, flake-utils, ... }@inputs:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
        version = "0.1.0";
      in
      {
        formatter = pkgs.nixpkgs-fmt;
        packages = rec {
          amify = pkgs.buildGoModule {
            pname = "amify";
            version = version;
            src = ./.;
            vendorHash = null;
            CGO_ENABLED = 0;
            subPackages = [ "cmd/amify" ];
          };
          dockerImage = pkgs.dockerTools.buildImage {
            name = "murtazau/amify";
            tag = version;
            copyToRoot = with pkgs.dockerTools; [
              caCertificates
            ];
            config = {
              Cmd = [ "${amify}/bin/ellipsis" ];
              WorkingDir = "/data";
            };
          };
          default = amify;
        };
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            nixd
            nixpkgs-fmt
            go
            go-tools
            gopls
          ];
        };
      });
}
