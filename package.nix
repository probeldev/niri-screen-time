{
  buildGoModule
}:
buildGoModule {
  name = "niri-screen-time";
  src = ./.;
  vendorHash = "sha256-9y1F2ZrmpiQJ9ZTq9SoRE2PxR65DDNCeBKf4M0HUQC4=";
  meta.mainProgram = "niri-screen-time";
}
