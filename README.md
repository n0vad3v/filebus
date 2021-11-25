# Filebus

> A simply upload fileserver for fun

## Build

Binary:

```
go build -o filebus .
```

Docker

```
docker build . -t n0vad3v/filebus:lastest
```

## Usage

Filebus will by default listen on `0.0.0.0:3000`, so using it with Docker is **strongly recommended.**

It supports two ways of uploading.

### Upload method 1

You must fill in `file` and `filepath` variables

```bash
➜  ~ curl -F file=@main.py -F filepath=path/to/main.py http://localhost:3000/upload
{"md5":"5f3962dfab1b52e49f58b4fdaa22dc27","url":"http://filebus.nova.moe/path/to/main.py"}% 
```

### Upload method 2

```bash
➜  ~ curl -F test/main.py=@main.py http://localhost:3000/upload
{"md5":"5f3962dfab1b52e49f58b4fdaa22dc27","url":"http://filebus.nova.moe/test/main.py"}%           
```

### Download

```bash
➜  ~ curl http://localhost:3000/path/to/main.py
import json
import requests
import os
...
```

## Deployment

### Without DB

Filebus can be deployed without using any DB, and it's now a stateless service, you can create a `docker-compose.yml` file like this to quickly use Filebus:

```yml
version: '3'

services:

  filebus:
    image: n0vad3v/filebus:lastest
    restart: always
    environment:
      FILEBUS_URL: "https://filebus.nova.moe/"
    ports:
      - '0.0.0.0:3000:3000'
    volumes:
      - ./data:/data
```


### With DB

First you need to spin up a working TiDB(or MySQL), create database called `filebus`, and initialize table called `upload_logs` with:

```sql
CREATE TABLE upload_logs (
    filename varchar(255),
    filepath varchar(255),
    filesize BIGINT,
    filehash varchar(255),
    uploader_ip varchar(255),
    uploaded_at DATETIME
);
```

`docker-compose.yml` example when using DB:

```yaml
version: '3'

services:

  filebus:
    image: n0vad3v/filebus:latest
    restart: always
    environment:
      FILEBUS_URL: "https://filebus.nova.moe/"
      ENABLE_LOG: "TRUE"
      DB_HOST: "db"
      DB_USERNAME: "root"
      DB_PASSWORD: "password"
      DB_DBNAME: "filebus"
    ports:
      - '0.0.0.0:3000:3000'
    volumes:
      - ./data:/data

  db:
    image: mysql:5.7
    restart: always
    environment:
      MYSQL_DATABASE: 'filebus'
      MYSQL_USER: 'user'
      MYSQL_PASSWORD: 'password'
      MYSQL_ROOT_PASSWORD: 'password'
    ports:
      - '0.0.0.0:3306:3306'
    volumes:
      - ./db_data:/var/lib/mysql

```

## License

Filebus is under the GPLv3. See the [LICENSE](./LICENSE) file for details.