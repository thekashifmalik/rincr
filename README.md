# rincr
Tool to make **r**emote-**incr**emental backups which consist of a mirror and historical snapshots of the target data.
These snapshots are created using hard links and can be taken as often as needed as they are cheap both in time and
space. `rincr` can manage these snapshots, clean up unnecessary data and also restore files to earlier versions.

> **Note**: This project was previously called *kbackup*. See this [issue](https://github.com/thekashifmalik/kbackup/issues/2)
> for more information.

## Quickstart
Create incremental backups with a simple command:

```bash
# Backup local folder to remote server
rincr ~/mydata user@server:~/backups

# Backup remote folder to local directory
rincr user@server:~/mydata ~/backups

# Backup multiple sources
rincr ~/docs ~/photos user@server:~/backups
```

`rincr` creates space-efficient snapshots using hard links. Each backup is a complete copy that you can browse and restore from directly. Historical snapshots are stored in `.rincr/` subdirectories at the destination.

## Installation
Download the right binary for your OS and Architecture (all binaries available on [Github](https://github.com/thekashifmalik/rincr/releases)):
- [Linux on x86](https://github.com/thekashifmalik/rincr/releases/latest/download/rincr-linux-amd64)
- [MacOS on ARM](https://github.com/thekashifmalik/rincr/releases/latest/download/rincr-darwin-arm64)

Make sure it is executable and then place it anywhere in your `$PATH`:

```bash
chmod +x rincr
mv rincr ~/bin/rincr # Use whatever location you put your user binaries in.
```

Test that it works:

```bash
rincr --version
```

> **Note**: This software is not yet stable; there may be backwards-incompatible changes before v1. Use at your own
> risk.


## Features

### Backup Management

Create backups with automatic pruning to manage disk space:
```bash
# Backup with default retention (24h, 30d, 12m, 10y)
rincr --prune ~/mydata myserver:~/backups

# Custom retention rules
rincr --prune --hourly=12 --daily=15 --monthly=6 --yearly=5 ~/mydata myserver:~/backups
```

Prune existing backups without syncing new data:
```bash
rincr prune myserver:~/backups/mydata
```

> **Note**: When pruning, provide the full path including the target directory name.

### File Restoration

Restore files from any backup snapshot:

```bash
# Restore from latest backup
rincr restore --latest myserver:~/backups/mydata path/to/file ~/restored/

# Restore from specific time point
rincr restore --from=12h myserver:~/backups/mydata path/to/file ~/restored/

# Restore multiple files
rincr restore --latest myserver:~/backups/mydata file1 file2 dir/ ~/restored/
```

Files are searched from newest to oldest backups. Only changed files are transferred using rsync.



## Example

```bash
# Initial state - source has files
$ ls -A ~/mydata
README.md config.json

# Destination has previous backup
$ ssh server ls ~/backups/mydata
README.md .rincr

$ ssh server ls ~/backups/mydata/.rincr
2024-01-12T09-12-32/ last

# Run backup - new file gets synced
$ rincr ~/mydata server:~/backups
Rotating backup: 2024-01-12T09-12-32
Syncing changes...
Created snapshot: 2024-01-23T18-26-10

# New snapshot created, old one preserved
$ ssh server ls ~/backups/mydata
README.md config.json .rincr

$ ssh server ls ~/backups/mydata/.rincr
2024-01-12T09-12-32/ 2024-01-23T18-26-10/ last
```

## References
- [rsync](https://rsync.samba.org/) - The underlying sync engine
- [rsnapshot](https://rsnapshot.org/) - Inspiration for snapshot-style backups
- [Easy Automated Snapshot-Style Backups with Linux and Rsync](http://www.mikerubel.org/computers/rsync_snapshots/) - Technical foundation
- [Development Blog](blog) - Design decisions and implementation notes
