package getter

func FileId(f *model.File) int64    { return f.Id }
func FileHash(f *model.File) string { return f.Hash }
