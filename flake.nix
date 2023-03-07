{
  description = "A flake for go dev env";
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let pkgs = nixpkgs.legacyPackages.${system};
      in {
        devShell = pkgs.mkShell {
          buildInputs = with pkgs; [
            buf
            protobuf
          ];
          postShellHook = ''
            echo "Welcome to buf env build by nix"
            export PATH="~/go/bin:$PATH"
            export PATH=$PATH:/usr/local/go/bin
          '';
        };
      });
}
