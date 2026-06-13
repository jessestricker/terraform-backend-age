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

        terraform-backend-age = pkgs.buildGoModule {
          pname = "terraform-backend-age";
          version = "0.0.0-dev";
          src = ./.;
          vendorHash = "sha256-vcMWw5iL9UD7QihrMWu5Irh+n8ss5zD8V9MocJLZ2O4=";
        };
      in
      {
        packages = {
          inherit terraform-backend-age;
          default = terraform-backend-age;
        };

        devShells.default = pkgs.mkShell {
          inputsFrom = [ terraform-backend-age ];
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
