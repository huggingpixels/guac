{
  description = "guac";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-compat.url = "https://flakehub.com/f/edolstra/flake-compat/1.tar.gz";
    flake-parts.url = "github:hercules-ci/flake-parts";
  };

  outputs = inputs @ {flake-parts, ...}:
    flake-parts.lib.mkFlake {inherit inputs;} {
      systems = ["aarch64-linux" "x86_64-linux" "aarch64-darwin" "x86_64-darwin"];
      perSystem = {
        config,
        self',
        inputs',
        pkgs,
        system,
        ...
      }: {
        devShells.default = pkgs.mkShell {
          nativeBuildInputs = with pkgs; [
            colima
            docker
            docker-compose
            gcc
            go
            go-outline
            go-tools
            golangci-lint
            gopls
            goreleaser
            gotests
            jq
            nats-server
            protobuf
            protoc-gen-go
            protoc-gen-go-grpc
          ];
        };
        formatter = pkgs.alejandra;
      };
    };
}
