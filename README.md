# Balena Data Extractor

[![balena deploy button](https://www.balena.io/deploy.svg)](https://dashboard.balena-cloud.com/deploy?repoUrl=https://github.com/maggie0002/balena-data-extractor)

This is an experimental project for extracting various forms of data (such as device information or logs) from Balena devices and uploading it to [PrivateBin](https://privatebin.info/directory/); a secure online version of PasteBin. 

All of the extracted content is encrypted on the device before being uploaded to PrivateBin meaning PrivateBin cannot see any of the stored content. You can retrieve data by using the URL returned by this container, which includes your key for decrypting the content. For added security, you can apply a password required when accessing the content on the PrivateBin website.

As an added bonus and for even greater security and privacy, you can run your own instance of the open source [PrivateBin](https://privatebin.info/directory/) project on your own server and then pass the URL of your server to this container via an environment variable or to the executable (see below).

 
# Basic Usage:

Run the container and the default mode will extract the following content and create an individual URL for each:

```
Device Info (via the Balena Supervisor)
Environment Variables (API_KEY variables are filtered out)
JournalCtl Logs
A List of Network Interfaces
```

## With the Docker Compose file:

Add the `balena-data-extractor` section of docker-compose.yml file in this repository to your own docker-compose file.

## With a `run` command on a device:

`balena run bcr.io/maggie0002/balena-data-extractor`

# Advanced Usage

You can change the default PrivateBin instance used by modifying the PRIVATEBIN_URL in the Docker Compose file or by [passing the env variable](https://docs.docker.com/engine/reference/run/#env-environment-variables) to the `balena run` command.

You can use a cmds.yaml file to generate your own commands to execute and pass to PrivateBin (see cmds.example.yaml file) by passing `-data yaml` and copying a cmds.yml file in to the container with the Dockerfile.

You can also set additional options by passing them in the Docker Compose command field or by putting them at the end of your `balena run` command:

  -burn
        Burn all data after being read once
  -expire string
        Delete all data after specified time. Options are: 'hour', 'day', 'week' or 'month' (default "day")
  -data string
        Choose which data to export. Options are: 'all', 'deviceinfo', 'envvars', 'journalctl', 'networkinterfaces', 'yaml') (default "all")
  -password string
        Set a password for accessing the uploaded content
  -url string
        Override the default data host with the passed URL
  -help
        Show this content

* Only the last 10000 lines of the JournalCtl logs are returned otherwise the browser struggles to decrypt it. 