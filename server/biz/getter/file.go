package getter

import "onij/infra/mysql"

func FileId(f *mysql.File) int64 {return f.Id}