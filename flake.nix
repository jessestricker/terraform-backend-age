{
  outputs =
    { nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;
          config.allowUnfree = true;
        };
      in
      rec {
        packages.default = pkgs.buildGoModule {
          pname = "terraform-backend-age";
          version = "0.0.0-dev";
          src = ./.;
          vendorHash = "sha256-vcMWw5iL9UD7QihrMWu5Irh+n8ss5zD8V9MocJLZ2O4=";
        };

        devShells.default = pkgs.mkShell {
          inputsFrom = [ packages.default ];
          packages = with pkgs; [
            age
            delve
            gopls
            terraform
          ];
        };
      }
    );
}
