{ buildGoModule, fetchFromGitHub }:
{
  default = buildGoModule {
    pname = "muxie";
    version = "1.1.0";
    subPackages = [ "cmd/muxie" ];
    src = fetchFromGitHub {
      owner = "phanorcoll";
      repo = "muxie";
      rev = "65c7e7101f8a2c83fae9181907e1b01094dd9be2";
      hash = "sha256-JO6b1Nbm82FLnTSK7sKdZVTNIUznwwFcFU/dxpHu5pM=";
    };
    meta.mainProgram = "muxie";
    vendorHash = "sha256-CXd2j180T9ln21RTBCqCqdO32aeNIXHwiPRNcFPFt2I=";
  };
}
