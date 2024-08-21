# lists
A webserver for keeping lists. The database is a sqlite database where each
list had a name and body text. The site has a create list page and an edit list
page.

# usage

A `go` compiler is required to compile this application. Check `go.mod` for the
oldest supported version of [go](https://go.dev/). Then run `make` to compile
the project.

You will probably want to install this database migration tool in order to
create your sqlite database: https://github.com/golang-migrate/migrate

Once installed you can spin up a new database, or update an existing one with
the following command. This should be run every time you pull new updates for
the server.
```sh
migrate -path=./migrations -database=sqlite://lists.db up
```

Run with the `-help` flag for current options. The `-addr` argument is used to
select a port / address to bind to and the `-dsn` argument is used to provide a
path to your sqlite database file along with any database options you wish to
use.

If you're actually trying to host this on the internet make sure you put it
behind a proxy (caddy, nginx, openbsd httpd, apache, etc) as the application
does not serve https by itself.

# author
Written and maintained by Dakota Walsh.
Up-to-date sources can be found at https://git.sr.ht/~kota/lists/

# license
GNU AGPL version 3 or later, see LICENSE.
