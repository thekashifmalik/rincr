# kbackup
Quick and easy incremental backups.

A reimplementation of the ideas in `rsnapshot` built with the composibility and simplicity of `rsync` in-mind. It can be
used by-hand, in scripts or as part of a software system.

## Quickstart
Grab the right binary for your OS and ARCH from [Github](https://github.com/thekashifmalik/kbackup). You can place this
anywhere on your $PATH. For `linux` and `amd64` we do:
```
curl -L https://github.com/thekashifmalik/kbackup/releases/latest/download/kbackup-linux-amd64 > kbackup
chmod +x kbackup
mv kbackup ...
```

Create an incremental backup of `~/mydata` at the remote location `myserver:~/backups/mydata`:
```
kbackup ~/mydata myserver:~/backups
```

Historical snapshots are stored in `myserver:~/backups/mydata/.kbackup`.

If you also want to clean up older snapshots, pass the `--prune` option:
```
kbackup --prune ~/mydata myserver:~/backups
```

This will apply the default retention rules and keep:
- 24 hourly backups
- 30 daily backups
- 12 monthly backups
- 10 yearly backups

Any pruned backups will be printed out.

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

### Encrypted
`kbackup` uses `rsync` which uses `ssh` so all communication is encrypted. Addtionally backups can be encrypted on
disk (NOT IMPLEMENTED YET).

### Deduplicated
`kbackup` uses the hard-link mechanism from `rsnapshot` so unchanged files between snapshots do not use any space.

### Bidirectional
`kbackup` can backup from/to a local or remote location, just like `rsync` (NOT IMPLEMENTED YET).

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
