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
task completion-mac
```

To Linux

```sh
task completion-linux
```

### To generate code to talk to Strava API

```sh
git clone https://github.com/swagger-api/swagger-codegen
cd swagger-codegen
git checkout master
./run-in-docker.sh mvn package
./run-in-docker.sh generate -i https://developers.strava.com/swagger/swagger.json -l go -o generated/go
rm -rf ../strava-cli/swagger
mv generated/go ../strava-cli/swagger
```
