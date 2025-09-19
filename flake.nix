{
  description = "Tool Tracker Go Development Environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go # Latest Go version
            gopls # Go language server
            gotools # Additional Go tools (goimports, etc.)
            air # Hot reload tool (since you're using it)
            docker
            docker-compose
            # Development tools
            git
            curl
            jq
          ];

          shellHook = ''
            echo "Tool Tracker development environment loaded!"
            echo "Go version: $(go version)"
            echo "Available commands: go, air, docker, docker-compose"

            # Set up Go environment
            export GOPATH="$PWD/.go"
            export GOCACHE="$PWD/.go/cache"
            mkdir -p "$GOPATH" "$GOCACHE"

            # Add Go bin directory to PATH for tools like mockgen
            export PATH="$GOPATH/bin:$PATH"

            # VS Code will inherit these environment variables
            echo "GOPATH set to: $GOPATH"
            echo "Go tools directory added to PATH"
          '';
        };
      }
    );
}
