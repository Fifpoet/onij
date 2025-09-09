package getter

import "onij/infra/mysql"

func FileId(f *mysql.File) int64 { return f.Id }
func FileHash(f *mysql.File) string { return f.Hash }