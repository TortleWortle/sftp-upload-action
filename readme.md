# sftp-upload-action
Made this so I don't have to use sftp for school.

## Usage

```yaml
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: TortleWortle/sftp-upload-action@master
        with:
          server: 127.0.0.0:22
          username: ${{ secrets.FTP_USERNAME }}
          password: ${{ secrets.FTP_PASSWORD }}
          local_dir: dist
          remote_dir: /home/user/website

```

## TODO

- [x] actually write instructions for this thing
- [x] multi stage docker file
- [ ] deletion of files that don't belong to git and aren't in gitignore