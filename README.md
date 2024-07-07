# rincr
Tool to make **r**emote-**incr**emental backups which consist of a mirror and historical snapshots of the target data.
These snapshots can be taken as often as needed and are cheap both in time and space. `rincr` can also manage these
snapshots, restore files to earlier versions and encrypt backups (NOT IMPLEMENTED YET).

> **Note**: This project was previously called *kbackup*. See this [issue](https://github.com/thekashifmalik/kbackup/issues/2)
> for more information.

## Quickstart
`rincr` supports a similar interface to `rsync`:

```
rincr [[USER@]HOST:]SRC... [[USER@]HOST:]DEST
```

When this command is run `rincr` checks to see if there are existing backups at the destination. If there are, a copy of
the latest backup is created using hard links. Then new changes (if any) are synced from the source.

Incremental backups are kept at `DEST/SRC/.rincr`. They are fully browsable and take no extra space for files that have
not changed between versions. Each file acts as a full-backup and can be copied back out manually to restore data to an
older version.

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

### Push Backups

Create an incremental backup of `~/mydata` at the remote location `myserver:~/backups/mydata`:
```bash
rincr ~/mydata myserver:~/backups
```

### Pull Backups

We can also back up locally from remote locations:
```bash
rincr server1:~/mydata server2:~/otherdata ~/backups
```

This will create incremental backups in `~/backups/mydata` and `~/backups/otherdata`.

### Pruning

If you also want to clean up older backups, pass the `--prune` option:
```bash
rincr --prune ~/mydata myserver:~/backups
```

This will apply the default retention rules and delete extra backups. Pruned backups will be printed out. The default
retention rules are:
- 24 hourly backups
- 30 daily backups
- 12 monthly backups
- 10 yearly backups


If you just want to prune backups without syncing new data, you can use:

```bash
rincr prune myserver:~/backups/mydata
```
> **Note**: When pruning data you have to provide the path to the actual backup destination which includes the target
> name `mydata`.

If you want to override the default retention rules, you can use the options:

```bash
rincr --prune --hourly=12 --daily=15 --monthly=6 --yearly=5 ~/mydata myserver:~/backups
```

### Restoring Files

To restore files from a backed-up repository, we can use:

```bash
rincr restore --latest myserver:backups/mydata path-to-restore output-path
```

`rincr` will check backups from latest to oldest until it finds a matching path. It will then copy that path recursively
into the output location. We can restore single files or full directory trees this way. Since `rsync` is used for the
the underlying transfer, only necessary files are transferred. Mutiple paths can be restored in one go:


```bash
rincr restore --latest myserver:backups/mydata file-1 file-2 output-path
```

Any paths that are not found will skipped. We can also configure how old of a backup we want to fetch:


```bash
rincr restore --from=12h myserver:backups/mydata path-to-restore output-path
```

It will find the closest backup to the given duration (from current time) and copy the path if it finds one.



## Demo

```bash
$ ls -A ~/mydata
README.md new-file.txt

$ ssh desktop-1 ls ~/Backups/mydata
README.md .rincr

$ ssh desktop-1 ls ~/Backups/mydata/.rincr
2024-01-12T09-12-32 last

$ rincr ~/mydata desktop-1:~/Backups
...

$ ssh desktop-1 ls ~/Backups/mydata
README.md new-file.txt .rincr

$ ssh desktop-1 ls ~/Backups/mydata/.rincr
2024-01-12T09-12-32 2024-01-23T18-26-10 last
```

## Unimplemented

```bash
# File encryption - supports deduplication but less secure
rincr --encrypt ~/mydata myserver:~/backups

# Folder encryption - no deduplication but completely secure
rincr --encrypt-folder ~/mydata myserver:~/backups

```

## References
- [rsync](https://rsync.samba.org/)
- [rsnapshot](https://rsnapshot.org/)
- [Easy Automated Snapshot-Style Backups with Linux and Rsync](http://www.mikerubel.org/computers/rsync_snapshots/)
- [WIP Blog](blog)
