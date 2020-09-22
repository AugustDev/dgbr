# Dgraph Backup Restore (dgbr)

## Details
This tool allows to export/import Dgraph database to/from AWS S3 bucket. Backups are created by requesting Dgraph to export existing database, archiving the data and uploading it to the S3 bucket. Backups can be configured to be made periodically. Database restore allows to download existing archive from S3 bucket and import it to the database using live loader. `dgbr` is currently used in production and is planned to be maintained and improved.

## Install
You can either download binary release or build `dgbr` from source.
### From source
`dgbr` is written in Go, therefore you have to instal the compiler to build it from source.
```
git clone https://github.com/AugustDev/dgbr.git
go build
go install
```
### Release
If you choose to download the release, you may add it to your bin path like this.
```
wget https://github.com/AugustDev/dgbr/releases/download/0.9.1/dgbr-linux-amd64.gz
gunzip -c dgbr-linux-amd64.gz > dgbr
chmod +x dgbr
sudo mv dgbr /usr/bin/dgbr
```

## Usage
Both Dgraph backup and restore actions require specifying: `AWS_ACCESS_KEY`, `AWS_SECRET_KEY`, `bucket` and `region`.
### Backup
In adition to aforementioned variables `dgbr backup` requires to specify Dgraph `export` path which is specified when starting `dgraph alpha` (Read in Notes). Make sure Dgraph has permission to write to `export` path and that the user calling `dgbr` has permission to access and then delete files from `export` path.

```bash
dgbr backup \
--AWS_ACCESS_KEY=X \
--AWS_SECRET_KEY=Y \
--bucket=my-dgraph-backups \
--region=eu-west-1 \
--export=/exports
```
### Restore
Restoring database requires specifying name of the `zip` file in S3 bucket.

```bash
dgbr restore \
--AWS_ACCESS_KEY=X \
--AWS_SECRET_KEY=Y \
--bucket=views-dgraph-backups-development \
--region=eu-west-1 \
--name=dg-my_backup.zip
```

## Periodic backups
To schedule perodic (daily, hourly etc.) backups simply create a script and add it to your cron list. Make sure that `dgbr` is in the appropriate `bin` folder.

```bash
#!/bin/bash
export PATH="/usr/local/bin:/usr/bin:/bin"

AWS_ACCESS_KEY=X
AWS_SECRET_KEY=Y
BUCKET=views-dgraph-backups-development
REGION=eu-west-1
EXPORT_PATH=/Users/august/exports

dgbr backup \
--AWS_ACCESS_KEY=${AWS_ACCESS_KEY} \
--AWS_SECRET_KEY=${AWS_SECRET_KEY} \
--bucket=${BUCKET}  \
--region=${REGION}  \
--export=${EXPORT_PATH} 
```

An example of hourly export cronjob would be
```bash
0 * * * * /Users/august/backup.sh >> /Users/august/log.txt 2>&1
```
where the logs are saved to `/Users/august/log.txt`.

## CLI help

## Notes
##### What is my export path?

If you start `dgraph alpha` using
```
dgraph alpha --lru_mb 2048 -p /var/run/dgraph/p -w /var/run/dgraph/w --export /exports
```
your export path is `/exports`.

##### Origin

Project was inspired by `https://github.com/xanthous-tech/dgraph-backup-cli`, yet I found it to have multiple issues and lacking documentation.