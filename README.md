# Dgraph Backup-Restore 

## Details
This tool allows to export and import Dgraph database to AWS S3 bucket. Backups are created by requesting Dgraph to export existing database, archiving the data and uploading it to the S3 bucket. Backups can be configured to be made periodically. Database restore allows to download existing archive from S3 bucket and import it to the database using live loader.

## Usage
Both Dgraph export and import actions require specifying: `AWS_ACCESS_KEY`, `AWS_SECRET_KEY`, `bucket` and `region`.
### Backup
In adition to aforementioned variables export requires to specify Dgraph `export` path which is specified when starting `dgraph alpha`. 

```bash
dgbr backup \
--AWS_ACCESS_KEY=X \
--AWS_SECRET_KEY=Y \
--bucket=my-dgraph-backups \
--region=eu-west-1 \
--export=/exports
```
### Restore
Importing database requires specifying name of the `zip` file in S3 bucket.

```bash
dgbr restore \
--AWS_ACCESS_KEY=AKIASEJMBX3ZVP6IVP4L \
--AWS_SECRET_KEY=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 \
--bucket=views-dgraph-backups-development \
--region=eu-west-1 \
--name=dg-my_backup.zip
```

## Periodic backups
TBD

### FAQ
##### What is my export path?

If you start `dgraph alpha` using
```
dgraph alpha --lru_mb 2048 -p /var/run/dgraph/p -w /var/run/dgraph/w --export /exports
```
your export path is `/exports`.