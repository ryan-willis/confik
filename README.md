# confik

Insert configuration files into Docker volumes using labels.

## Why?

Honestly, this is one of the silliest tools I've built. I manage my own personal docker servers remotely over ssh using `DOCKER_HOST=ssh://remote_server_hostname_or_ip`, and I don't like storing configuration files or checking out a git repo on the actual host itself. Therefore, I don't like using bind mounts for configuration. I prefer to have all of the configuration for a service stored in a git repo that doesn't get checked out on the host that I can make changes to and simply run `docker compose up -d` to apply them.

## Usage

To get the most out of `confik`, you should use it with Docker Compose. I personally use it to configure [watchtower](https://github.com/containrrr/watchtower). Here's an example:

```yaml

services:
  watchtower:
    image: containrrr/watchtower
    command:
      - --label-enable
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - watchtower-confik:/conf/watchtower
    environment:
      DOCKER_CONFIG: /conf/watchtower
    labels:
      # not necessary for watchtower, but an option if you want
      # to restart the container when the configuration changes
      # defaults to false
      - confik.restart=true

      # the path to the target file in the confik container's volume mount
      - confik.file=/conf/watchtower/config.json

      # the contents of the configuration file
      - confik.contents={"auths":{"ghcr.io":{"auth":"$GHCR_TOKEN"}}}
  confik:
    image: ghcr.io/ryan-willis/confik
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - watchtower-config:/conf/watchtower
      - confik-state:/confik_state

volumes:
  watchtower-config:
    external: false
  confik-state:
    external: false

```

## Notes

Currently, `confik` only supports one configuration file per service.

The value of the `confik.contents` label will be stored in raw form, so you can't really use formats that require newlines (yet). I'll probably add support for converting inline JSON to various formats like YAML in the future via some `confik.format` label, or allowing for escaped characters like newlines to be included in the `confik.contents` label. TBD.

## License

MIT
