variable "REGISTRY" {
  default = "docker.io/library"
}

variable "NAME" {
  default = ""
}

variable "VERSION" {
  default = ""
}

variable "GIT_COMMIT" {
  default = ""
}

target "default" {
  context    = ".."
  dockerfile = "docker/Dockerfile"
  platforms  = ["linux/amd64", "linux/arm64"]
  output     = ["type=registry"]
  tags = [
    "${REGISTRY}/${NAME}:${GIT_COMMIT}",
    "${REGISTRY}/${NAME}:${VERSION}",
    "${REGISTRY}/${NAME}:${split(".", VERSION)[0]}",
    "${REGISTRY}/${NAME}:${split(".", VERSION)[0]}.${split(".", VERSION)[1]}",
    "${REGISTRY}/${NAME}:latest",
  ]
}
