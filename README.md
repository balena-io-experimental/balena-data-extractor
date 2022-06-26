# Balena Data Extractor

[![balena deploy button](https://www.balena.io/deploy.svg)](https://dashboard.balena-cloud.com/deploy?repoUrl=https://github.com/maggie0002/balena-data-extractor)

This is an experimental project for extracting various forms of data (such as device information or logs) from Balena devices and uploading it to [PrivateBin](https://privatebin.info/directory/); a secure online version of PasteBin. 

All of the extracted content is encrypted on the device before being uploaded to PrivateBin meaning PrivateBin cannot see any of the stored content. You can retrieve data by using the URL returned by this container, which includes your key for decrypting the content. For added security, you can apply a password required when accessing the content on the PrivateBin website.

As an added bonus and for even greater security and privacy, you can run your own instance of the open source [PrivateBin](https://privatebin.info/directory/) project on your own server and then pass the URL of your server to this container via an environment variable or to the executable (see below).

 
# Basic Usage:

Run the container and the default mode will extract the following content and create an individual URL for each:

- Device info (via the Balena Supervisor)
- OS release info for the container being used by this app
- Environment variables (API_KEY variables are filtered out)
- JournalCtl logs
- A List of available network interfaces


### With the Docker Compose file:

Add the `balena-data-extractor` section of docker-compose.yml file in this repository to your own docker-compose file.

### With a `run` command on a device (not compatible with processes that require Balena Supervisor access):

`balena run bcr.io/maggie0002/balena-data-extractor`

# Advanced Usage

You can change the default PrivateBin instance used by modifying the `PRIVATEBIN_URL` in the Docker Compose file or by [passing the env variable](https://docs.docker.com/engine/reference/run/#env-environment-variables) to the `balena run` command.

You can edit the cmds.yaml file to amend the current commands or add your own commands to execute and pass to PrivateBin. Options include:

### Request to the Balena Supervisor
```
balena_supervisor_device_info:
    name: Device Info # Name that precedes the export URL
    cmd_type: api # Specifies that this is an API request 
    url: /v1/device # End of the path for the requested content. The Supervisor URL is retrieved automatically. 
    supervisor: true # Indicate that the Supervisor is being requested
```

### Request to an API that returns JSON content
```
request_from_api:
    name: Api Request # Name that precedes the export URL
    cmd_type: api # Specifies that this is an API request 
    url: http://0.0.0.0/path/ # Full path for the requested content.
    supervisor: false # Indicates that this is not a Supervisor request
```

### Request that returns the content of a file (for example a log file)
```
os_file:
    name: OS Release Info # Name that precedes the export URL
    cmd_type: file # Specifies that this is a request to read a local file 
    path: /etc/os-release # File path
```

### Request that returns the output of a shell command
```
network_interfaces:
    name: Network Interfaces # Name that precedes the export URL
    cmd_type: shell # Specifies that this is a command to execute and return the output
    cmd: ls -lah # shell command to execute
```

### Additional options

You can also set additional options by passing them in the Docker Compose command field or by putting them at the end of your `balena run` command:

```
  -burn
        Burn all data after being read once
  -expire string
        Delete all data after specified time. Options are: 'hour', 'day', 'week' or 'month' (default "day")
  -password string
        Set a password for accessing the uploaded content
  -url string
        Override the default data host with the passed URL
  -help
        Show this content
```
