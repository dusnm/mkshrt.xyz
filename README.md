# mkshrt.xyz
A fully featured FOSS, ad & tracker free URL shortener.

Features:

* Can be self-hosted
* Configured with YAML
* Shortens URLs to at most 28 characters. (with the default domain)
* No ads, tracking or other nonsense

Dependencies:

* MariaDB or MySQL

## Why?

1. It's fun to implement
2. Almost every other service just wants to steal your data

## Live Demo
[mkshrt.xyz](https://mkshrt.xyz)

## Building from source

Make sure you have the latest version of go available and in your `$PATH`.

Clone this repository.
```shell
git clone https://github.com/dusnm/mkshrt.xyz.git && cd mkshrt.xyz
```

if the target OS and/or architecture differs from the one you're compiling from, set the GOARCH and GOOS environment variables accordingly.
For example, to build for linux and amd64 you'd set them like this.
```shell
export GOOS="linux"
export GOARCH="amd64"
```

You can list all the available GOOS/GOARCH combinations by running the following command.
```shell
go tool dist list
```

Build the binary while bundling all dependencies.
```shell
export GOOS="linux" && \
export GOARCH="amd64" && \
export CGO_ENABLED=0 && \
go build -o mkshrt .
```

Make the resulting binary executable.
```shell
chmod +x ./mkshrt
```

You can now move the resulting binary to the `/usr/local/bin/` directory.
```shell
mv ./mkshrt /usr/local/bin/mkshrt
```

## Configuration

The application is written to look for a configuration file in the following directories.

1. `$XDG_CONFIG_HOME/mkshrt/config.yml`
2. `/etc/mkshrt/config.yml`

You can change the configuration file location easily, by setting the `$XDG_CONFIG_HOME` environment variable.

Example configuration is provided in `res/config.example.yml`

## Logging

You must create the log directory in whatever location you specified when you modified the configuration file.
The process running the application must have write access to that directory.

For example, if you set the log path as `/var/log/mkshrt/mkshrt.log`, and you're using the `www-data` user to run the process,
You'd set the permissions like this.
```shell
chown -R www-data:www-data /var/log/mkshrt
```

## Preparing the database
The application uses MySQL or MariaDB (recommended) as its database.

The database schema is provided in `res/schema.sql`. To import it, run the `mysql` utility.
```shell
mysql -u username -p database_name < ./res/schema.sql
```

## Running the application with systemd

A `systemd` unit file is provided in `res/mkshrt.service`. Place this file in the `/etc/systemd/system/' directory and reload the systemd daemon to make it aware of the new service.
```shell
cp ./res/mkshrt.service /etc/systemd/system/mkshrt.service && systemctl daemon-reload
```

You can then start and enable the service with this command.
```shell
systemctl enable --now mkshrt.service
```

Check the status of the service and look for errors, if any.
```shell
systemctl status mkshrt.service
```

If you configured everything correctly, the application should now be available at the socket you specified in the config file.
The default location is `http://localhost:6060`.

From here you can configure a reverse proxy, like nginx, to handle stuff like domains and TLS certificates.
