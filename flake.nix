{
	description = "A simple Debugger for OpenID Connect IdPs.";

	inputs = {
		nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
    flake-parts = {
      url = "github:hercules-ci/flake-parts";
      inputs.nixpkgs-lib.follows = "nixpkgs";
    };
		systems.url = "github:nix-systems/default";
	};

	outputs = inputs@{ flake-parts, ... }: flake-parts.lib.mkFlake { inherit inputs; } {
		systems = import inputs.systems;

		perSystem = { pkgs, self', ... }: {

			devShells.default = pkgs.mkShellNoCC {
				buildInputs = [ pkgs.go ];
			};

      packages.default = pkgs.buildGoModule {
        pname = "openid-connect-debugger";
        version = "1.1.0";
        src = inputs.self;
        vendorHash = null;
      };

			apps.default = {
				type = "app";
				program = pkgs.writeShellApplication {
					name = "openid-connect-debugger";
					text = "${self'.packages.default}/bin/openid-connect-debugger";
				};
			};
		};
	};
}
