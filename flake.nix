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
          alertfy = pkgs.buildGoModule {
            pname = "alertfy";
            version = version;
            src = ./.;
            vendorHash = "sha256-D8x1xSSyEgJprRpenB/peI6C1YsrG0pQSDmklW4BIKc=";
            CGO_ENABLED = 0;
            subPackages = [ "cmd/alertfy" ];
          };
          dockerImage = pkgs.dockerTools.buildImage {
            name = "murtazau/alertfy";
            tag = version;
            copyToRoot = with pkgs.dockerTools; [
              caCertificates
            ];
            config = {
              Cmd = [ "${alertfy}/bin/ellipsis" ];
              WorkingDir = "/data";
            };
          };
          default = alertfy;
        };
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            nixd
            nixpkgs-fmt
            go
            go-tools
            gopls
            kubernetes-helm
            kind
            kubectl
          ];
        };
      });
}
