{
  buildGoModule,
}:
{
  default = buildGoModule {
    pname = "muxie";
    version = "1.1.0";
    subPackages = [ "cmd/muxie" ];
    src = ./.;
    meta.mainProgram = "muxie";
    vendorHash = "sha256-CXd2j180T9ln21RTBCqCqdO32aeNIXHwiPRNcFPFt2I=";
  };
}
