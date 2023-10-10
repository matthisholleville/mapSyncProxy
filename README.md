<a href="https://www.haproxy.com">
    <img src="https://upload.wikimedia.org/wikipedia/commons/a/ab/Haproxy-logo.png" alt="Pritunl logo" title="Pritunl" align="right" height="100" />
</a>
<a href="https://cloud.google.com/storage">
    <img src="https://imgs.search.brave.com/5lxcXp7DQkSKquKVb6CQapUrgQTRsibWDzbcaLBqfi0/rs:fit:860:0:0/g:ce/aHR0cHM6Ly9zdGF0/aWMtMDAuaWNvbmR1/Y2suY29tL2Fzc2V0/cy4wMC9jbG91ZC1z/dG9yYWdlLWljb24t/MjU2eDIwNC1kb3Z3/cGp5eC5wbmc" alt="GCS logo" title="GCS" align="right" height="100" />
</a>

# Simple proxy to synchronize a GCS file with an HAProxy .map file

This small Golang server exposes a route that allows synchronizing a GCS (Google Cloud Storage) file with an HAProxy .map file using the Dataplane API.

## Use cases

This proxy can be used for various purposes:

1. Implement distributed rate limiting across multiple HAProxy instances without the need for an enterprise license.
2. Apply access rules (allowlisting/denylisting) across all machines in one or more environments.
3. Ensure the reliable updating of Map files.

## Demo

https://github.com/matthisholleville/mapSyncProxy/assets/99146727/63d140d9-c7ce-4d55-a460-7fe9dc4d96e0

## Local development

### 1. Requirements

- Docker
- Make
- Python3
- Gsutil
- Go

### 2. Configuration de HAProxy

To build and run a Docker container with the Dataplane API, execute the following command:

```bash
make setup
```

To verify that the container is running correctly, open a browser and navigate to the following URL: `http://127.0.0.1:5555/v2/docs`. You should see the Dataplane API documentation.

### 3. Starting the server

To start the server, execute the following command:

```bash
make run
```

You should see a list of metrics at the following URL: `http://localhost:8000/metrics`.

### 4. Pushing the Local File to a GCS Bucket

Before pushing the file to a bucket, ensure that it is accessible from your script execution context. To initiate the file push, execute the following command:

```bash
make push
```

### 5. Initiating Synchronization

To start synchronization, execute the following command:

```bash
make synchronize map_name=rate-limits bucket=$MY_BUCKET_NAME
```

If everything is successful, you should see the following message:

```bash
...
{"status":"synchronization success."}
```
