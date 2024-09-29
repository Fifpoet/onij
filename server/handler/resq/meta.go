package resq

type GetMetaResp struct {
	MetaEnumCode int    `json:"meta_enum_code"`
	MetaName     string `json:"meta_name"`
	MetaList     []MetaModel
}

type MetaModel struct {
	Value int    `json:"value"`
	Name  string `json:"name"`
}
