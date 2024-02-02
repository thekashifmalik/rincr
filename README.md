# kbackup
A reimplementation of `rsnapshot` built with the composibility and simplicity of `rsync` in-mind.

```
kbackup --rentention 30d ~/mydata myserver:~/Backups
```

## About
A utility written to compliment `rsync` for backups. Specifically, while `rsync` can be used to make _mirrored_ backups,
i.e exact clones of the target data, `kbackup` can be used to make _incremental_ backups which consist of snapshots of
the target data. These snapshots can be taken as often as needed and they are cheap both in time and space. `kbackup`
can optionally manage these snapshots for you, restore files to any stored versions and encrypt backups.

## Features
`kbackup` is built on top of `rsync` and borrows heavily from `rsnapshot`.

### Incremental
Whenever `kbackup` is run a copy of the last backup is stored in `.kbackup`. Backups are fully browsable.

### Differential
`kbackup` uses `rsync` so only the actual differeces b/w files are sent over-the-wire.

### Encrypted
`kbackup` users `rsync` which uses `ssh` so all communication is encrypted. Addtionally backups can be encrypted.

### Deduplicated
`kbackup` uses the same hard-link mechanism from `rsnapshot` so duplicate files b/w snapshots do not use extra space.


## Demo
```
$ ls -A ~/mydata
README.md new-file.txt

$ ssh desktop-1 ls ~/Backups/mydata
README.md .rincr

$ ssh desktop-1 ls ~/Backups/mydata/.rincr
2024-01-12T09-12-32 last

$ kbackup ~/mydata desktop-1:~/Backups
...

$ ssh desktop-1 ls ~/Backups/mydata
README.md new-file.txt .rincr

$ ssh desktop-1 ls ~/Backups/mydata/.rincr
2024-01-12T09-12-32 2024-01-23T18-26-10 last

$ kbackup --rentention 3d ~/mydata desktop-1:~/Backups
...

$ ssh desktop-1 ls ~/Backups/mydata/.rincr
2024-01-23T18-26-10 2024-01-23T18-28-09 last

```
