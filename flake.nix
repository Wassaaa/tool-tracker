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
            # Backend (Go)
            go # Latest Go version
            gopls # Go language server
            gotools # Additional Go tools (goimports, etc.)
            air # Hot reload tool (since you're using it)

            # Frontend (Node.js)
            nodejs_22 # Latest LTS Node.js version
            nodePackages.pnpm # Latest pnpm package manager

            # Container tools
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
            echo "Node version: $(node --version)"
            echo "pnqpm version: $(pnpm --version)"
            echo "Available commands: go, air, docker, docker-compose, node, pnpm"

            # Set up Go environment
            export GOPATH="$PWD/.go"
            export GOCACHE="$PWD/.go/cache"
            mkdir -p "$GOPATH" "$GOCACHE"

            # Add Go bin directory to PATH for tools like mockgen
            export PATH="$GOPATH/bin:$PATH"

            # Install Go tools if not already present
            if [ ! -f "$GOPATH/bin/mockgen" ]; then
              echo "Installing mockgen..."
              go install github.com/golang/mock/mockgen@latest
            fi

            if [ ! -f "$GOPATH/bin/swag" ]; then
              echo "Installing swag..."
              go install github.com/swaggo/swag/cmd/swag@latest
            fi

            # Set up Node.js environment
            export NODE_ENV="development"

            # VS Code will inherit these environment variables
            echo "GOPATH set to: $GOPATH"
            echo "Go tools directory added to PATH"
            echo "NODE_ENV set to: $NODE_ENV"
          '';
        };
      }
    );
}
