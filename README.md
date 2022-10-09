# strava-cli

## Commands

### Login

```sh
strava-cli login
```

### Show profile

```sh
strava-cli profile
```

### List activities

```sh
strava-cli list activity
```

### List starred segments

```sh
strava-cli list segment
```

### List recent efforts of a segment

```sh
strava-cli list segment-effort --id 631238
```

### Bash completion

To Mac

```sh
strava-cli completion bash > /opt/homebrew/etc/bash_completion.d/strava-cli
```

To Linux

```sh
strava-cli completion bash | sudo tee /etc/bash_completion.d/strava-cli
```

