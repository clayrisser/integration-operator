{ nixpkgs ? import <nixpkgs> {} }:

nixpkgs.stdenv.mkDerivation rec {
  name = "packages";
  buildInputs = [
    nixpkgs.gnumake42
    nixpkgs.gnused
    nixpkgs.jq
  ];
}
