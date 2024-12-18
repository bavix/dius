# Disk Utility (dius)

![ezgif-1-33af454190](https://user-images.githubusercontent.com/5111255/154798365-5f598167-5865-4352-99f0-f2c4ec974166.gif)

## Install

```bash
go install github.com/bavix/dius@latest
```

#### Alternative installation

```bash
cd /tmp
git clone https://github.com/bavix/dius.git
cd dius
go install -ldflags "-s -w"
```

## Usage

```bash
# pwd dir
dius

# user folder
dius ~

# mounts
dius /mnt/md0
```
