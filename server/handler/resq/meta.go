package resq

type GetMetaResp struct {
	MetaEnumCode int `json:"meta_enum_code"`
	MetaList     []MetaModel
}

type MetaModel struct {
	Value int    `json:"value"`
	Name  string `json:"name"`
}
