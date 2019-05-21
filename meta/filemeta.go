package meta

import (
	"fmt"
	"github.com/pibigstar/go-cloudstore/db"
	"sort"
)

// 文件元信息结构
type FileMeta struct {
	FileSha1 string
	FilePath string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetes map[string]FileMeta

func init() {
	// 初始化
	fileMetes = make(map[string]FileMeta)
}

// 新增/更新文件元信息
func UpdateFileMeta(fileMeta FileMeta) {
	fileMetes[fileMeta.FileSha1] = fileMeta
}

// 通过sha1值获取文件的元信息
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetes[fileSha1]
}

// 返回最后上传的文件元数据
func GetLastFileMeta(count int) []FileMeta {
	fMetaArray := make([]FileMeta, len(fileMetes))
	for _, v := range fileMetes {
		fMetaArray = append(fMetaArray, v)
	}
	sort.Sort(ByUploadTime(fMetaArray))
	return fMetaArray[0:count]
}

func RemoveFileMeta(filehash string) {
	delete(fileMetes, filehash)
}

// 上传文件后，将文件信息保存到mysql中
func UpdateFileMetaDB(meta FileMeta) bool {
	return db.InsertFile(meta.FileSha1, meta.FileName, meta.Location, meta.FileSize)
}

// 从数据库中获取文件元数据
func GetFileMetaDB(hash string) *FileMeta {
	file, err := db.GetFileMeta(hash)
	if err != nil {
		fmt.Printf("Failed to get file meta from db，err:%s\n", err.Error())
		return nil
	}
	fmeta := FileMeta{
		FileName: file.FileName.String,
		FileSha1: file.FileHash,
		Location: file.FileAddr.String,
		FileSize: file.FileSize.Int64,
	}
	return &fmeta
}
