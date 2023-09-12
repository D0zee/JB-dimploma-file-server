# File Server

### Methods:

- `POST /<file-system-path>`
- `GET /<file-system-path>`
- `DELETE /<file-system-path>`

## How to run:

1) Build image from Dockerfile: `docker build -t file-server .`
2) Run container: `docker run -p 8080:8080 file-server`

You can do these steps manually or just run ```launchFromDocker.sh```

Application listens on `8080` port by default, working directory is `/workDir`.
You can change these parameters if you want. Sources of application are placed in `/app` folder inside a container.

## How to use:

- To get a file from server: `curl --request GET http://localhost:8080/<path_to_file>`
- To load a file on
  server: `curl --request POST --data-binary "@<path_to_file_on_your_pc>" http://localhost:8080/<path_to_file>`
- To delete a file from server: `curl --request DELETE http://localhost:8080/<path_to_file>`

`path_to_file_on_your_pc` might be an absolute or relative path. Don't forget about `@` sign before a path to file on
your pc.

You can test this application by running the ``testScript.sh``. **Important note: before launching you must launch the
server.**
