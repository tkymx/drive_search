mkdir builds
ROOT=`pwd`

cd $ROOT/drive_search/
go build drive_search.go
mv drive_search ../builds/

cd $ROOT/drive_download/
go build drive_download.go
mv drive_download ../builds