# kbackup
Quick and easy incremental backups.

`kbackup` is a reimplementation of the ideas in `rsnapshot` built with the composibility and simplicity of `rsync`
in-mind. It can be used by-hand, in scripts or as part of a software system.

```
git clone git@github.com:thekashifmalik/kbackup.git
cd kbackup

bin/kbackup ~/mydata myserver:~/backups
```

This will create a mirror of `~/mydata` at the remote location `myserver:~/backups/mydata`. It will also rotate any
pre-existing snapshots and associated metadata in `myserver:~/backups/mydata/.kbackup`.

## About
A utility written to compliment `rsync` for backups. Specifically, while `rsync` can be used to make _mirrored_ backups,
i.e exact clones of the target data, `kbackup` can be used to make _incremental_ backups which consist of snapshots of
the target data. These snapshots can be taken as often as needed and they are cheap both in time and space. `kbackup`
can optionally manage these snapshots for you, restore files to earlier versions and encrypt backups.

## Features
`kbackup` is built on top of `rsync` and borrows heavily from `rsnapshot`.

### Incremental
Whenever `kbackup` is run a copy of the last backup is stored in `.kbackup`. Backups are fully browsable.

### Differential
`kbackup` uses `rsync` so only the actual differeces b/w files are sent over-the-wire.

### Encrypted
`kbackup` users `rsync` which uses `ssh` so all communication is encrypted. Addtionally backups can be encrypted on
disk.

### Deduplicated
`kbackup` uses the hard-link mechanism from `rsnapshot` so unchanged files b/w snapshots do not use any space.

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
# Move off prototype
kbackup ~/mydata myserver:~/backups

# Manual pruning
kbackup prune --rentention 30d myserver:~/backups/mydata

# Automatic pruning
kbackup --rentention 30d ~/mydata myserver:~/backups

# Restoring files
kbackup restore --from 1w myserver:~/backups/mydata/file ~/mydata/file

# File encryption
kbackup --encrypt ~/mydata myserver:~/backups

# Folder encryption
kbackup --encrypt-folder ~/mydata myserver:~/backups

```
