# kbackup
Quick and easy incremental backups. A reimplementation of the ideas in `rsnapshot` built with the composibility and
simplicity of `rsync` in-mind.

## Installation
Grab the right binary for your OS and ARCH from [Github](https://github.com/thekashifmalik/kbackup). Place this anywhere
in your path and make sure it is executable. For `linux` and `amd64` we do:
```
curl -L https://github.com/thekashifmalik/kbackup/releases/latest/download/kbackup-linux-amd64 > kbackup
chmod +x kbackup
mv kbackup ...
```

## Quickstart
`kbackup` supports a similar interface to `rsync`:

```
kbackup [HOST:]SRC... [HOST:]DEST
```

Create an incremental backup of `~/mydata` at the remote location `myserver:~/backups/mydata`:
```
kbackup ~/mydata myserver:~/backups
```

When this command is run the last backup (if any) is rotated and new changes are synced. Hard links and SSH is used
where possible.

If you also want to clean up older backups, pass the `--prune` option:
```
kbackup --prune ~/mydata myserver:~/backups
```

This will apply the default retention rules and delete extra backups. Pruned backups will be printed out. The default
retention rules are:
- 24 hourly backups
- 30 daily backups
- 12 monthly backups
- 10 yearly backups

Historical backups are kept in `myserver:~/backups/mydata/.kbackup`. They are fully browsable and take no extra space
for unchanged files.


We can also back up locally from remote locations:
```
kbackup server1:~/mydata server2:~/otherdata ~/backups
```
This will create incremental backups in `~/backups/mydata` and `~/backups/otherdata`.

## About
A utility to compliment `rsync` for backups. While `rsync` can be used to make _full_ backups, i.e exact clones of the
target data, `kbackup` can be used to make _incremental_ backups which consist of snapshots of the target data. These
snapshots can be taken as often as needed and they are cheap both in time and space. `kbackup` can also manage these
snapshots, restore files to earlier versions (NOT IMPLEMENTED YET) and encrypt backups (NOT IMPLEMENTED YET).

## Features
`kbackup` is built on top of `rsync` and borrows heavily from `rsnapshot`.

### Incremental
Whenever `kbackup` is run a copy of the last backup is stored in `.kbackup`. Backups are fully browsable.

### Differential
`kbackup` uses `rsync` so only the actual differences between files are sent over-the-wire.

### Deduplicated
`kbackup` uses the same hard-link mechanism that `rsnapshot` uses so unchanged files between snapshots do not use any
storage space.

### Bidirectional
`kbackup` can backup from (pull) or to (push) a remote location, just like `rsync`.

### Encrypted
`kbackup` uses `rsync` and `ssh` so all communication is encrypted. Addtionally backups can be encrypted on disk
(NOT IMPLEMENTED YET).


## Demo

```
$ ls -A ~/mydata
README.md new-file.txt

$ ssh desktop-1 ls ~/Backups/mydata
README.md .kbackup

$ ssh desktop-1 ls ~/Backups/mydata/.kbackup
2024-01-12T09-12-32 last

$ kbackup ~/mydata desktop-1:~/Backups
...

$ ssh desktop-1 ls ~/Backups/mydata
README.md new-file.txt .kbackup

$ ssh desktop-1 ls ~/Backups/mydata/.kbackup
2024-01-12T09-12-32 2024-01-23T18-26-10 last
```


## Unimplemented

```
# Manual pruning
kbackup prune myserver:~/backups/mydata

# Overriding default rentention
kbackup --rentention 30d ~/mydata myserver:~/backups

# Restoring files
kbackup restore --from 1w myserver:~/backups/mydata/file ~/mydata/file

# File encryption - supports deduplication but less secure
kbackup --encrypt ~/mydata myserver:~/backups

# Folder encryption - no deduplication but completely secure
kbackup --encrypt-folder ~/mydata myserver:~/backups

```

## References
- [rsync](https://rsync.samba.org/)
- [rsnapshot](https://rsnapshot.org/)
- [Easy Automated Snapshot-Style Backups with Linux and Rsync](http://www.mikerubel.org/computers/rsync_snapshots/)
