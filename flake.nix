{
  description = "Extensible chatbot";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs =
    { self, nixpkgs }:
    let
      pname = "supa8";
      user = "supa8-bot";
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };
      goServer = pkgs.buildGoModule {
        pname = pname;
        version = "git-${self.rev or "dirty"}";
        src = ./.;
        vendorHash = "sha256-izsIDqXlHH54cujpkrsAvgPGed6g2sjLU61MpLIRrJg=";
      };
    in
    {
      packages.${system}.default = goServer;

      nixosModules.default =
        {
          config,
          lib,
          pkgs,
          ...
        }:

        let
          cfg = config.services.${pname};
        in
        {
          options.services.${pname} = {
            authFile = lib.mkOption {
              type = lib.types.path;
              description = "Path to the bot's authentication file.";
              example = "/etc/${pname}/auth.yml";
            };
          };

          config.assertions = [
            {
              assertion = cfg.authFile != null;
              message = "You must set services.${pname}.authFile.";
            }
          ];

          users.users.${user} = {
            isSystemUser = true;
            home = "/var/lib/${user}";
            createHome = true;
          };

          systemd.services.${pname} = {
            after = [ "network.target" ];
            wantedBy = [ "multi-user.target" ];
            path = [
              pkgs.zbar
              pkgs.ffmpeg-full
            ];
            serviceConfig = {
              ExecStart = "${goServer}/bin/${pname}";
              Restart = "always";
              User = "${user}";
              WorkingDirectory = "/var/lib/${user}";
            };
          };

          environment.etc."/var/lib/${user}/auth.yml".source = cfg.authFile;
        };
    };
}
