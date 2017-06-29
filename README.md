# extract-gzip
extract all file with tar.gz under a folder.

# How to build
```
go build extract.go
```
This will generate a application named extract.exe under the source folder.

# How to use
```
extract --targetdir=<the-folder-where-gz-files-are> --desdir=<the-folder-where-gz-files-will-be-extracted-to>
e.g. 
$ extract --targetdir=D:\gzips --desdir=D:\files
```