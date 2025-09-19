package main

// AlbumResponse 对应 album.json 顶层结构
type AlbumResponse struct {
	ResourceState bool   `json:"resourceState"`
	Songs         []Song `json:"songs"`
	Code          int    `json:"code"`
	Album         Album  `json:"album"`
}

// Song 对应 songs 数组每一项
type Song struct {
	RtUrls          []string    `json:"rtUrls"`
	Ar              []Artist    `json:"ar"`
	Al              AlbumInfo   `json:"al"`
	St              int         `json:"st"`
	NoCopyrightRcmd interface{} `json:"noCopyrightRcmd"`
	SongJumpInfo    interface{} `json:"songJumpInfo"`
	DjID            int         `json:"djId"`
	No              int         `json:"no"`
	Fee             int         `json:"fee"`
	Mv              int         `json:"mv"`
	Cd              string      `json:"cd"`
	T               int         `json:"t"`
	V               int         `json:"v"`
	Rt              string      `json:"rt"`
	Mst             int         `json:"mst"`
	Cp              int         `json:"cp"`
	Crbt            interface{} `json:"crbt"`
	Cf              string      `json:"cf"`
	Dt              int         `json:"dt"`
	H               *Quality    `json:"h"`
	Sq              *Quality    `json:"sq"`
	Hr              interface{} `json:"hr"`
	L               *Quality    `json:"l"`
	RtUrl           *string     `json:"rtUrl"`
	Ftype           int         `json:"ftype"`
	Rtype           int         `json:"rtype"`
	Rurl            interface{} `json:"rurl"`
	Pst             int         `json:"pst"`
	Alia            []string    `json:"alia"`
	Pop             int         `json:"pop"`
	A               interface{} `json:"a"`
	M               *Quality    `json:"m"`
	Name            string      `json:"name"`
	ID              int         `json:"id"`
	Privilege       Privilege   `json:"privilege"`
}

// Quality 对应 h/m/l/sq 等音质对象
type Quality struct {
	Br   int `json:"br"`
	Fid  int `json:"fid"`
	Size int `json:"size"`
	Vd   int `json:"vd"`
	Sr   int `json:"sr"`
}

// Artist 对应 ar 中的歌手
type Artist struct {
	ID   int      `json:"id"`
	Name string   `json:"name"`
	Alia []string `json:"alia"`
}

// AlbumInfo 对应 al 中的专辑信息
type AlbumInfo struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	PicStr string `json:"pic_str"`
	Pic    int64  `json:"pic"`
}

// Privilege 权限信息
type Privilege struct {
	ID                 int                `json:"id"`
	Fee                int                `json:"fee"`
	Payed              int                `json:"payed"`
	St                 int                `json:"st"`
	Pl                 int                `json:"pl"`
	Dl                 int                `json:"dl"`
	Sp                 int                `json:"sp"`
	Cp                 int                `json:"cp"`
	Subp               int                `json:"subp"`
	Cs                 bool               `json:"cs"`
	Maxbr              int                `json:"maxbr"`
	Fl                 int                `json:"fl"`
	Toast              bool               `json:"toast"`
	Flag               int                `json:"flag"`
	PreSell            bool               `json:"preSell"`
	PlayMaxbr          int                `json:"playMaxbr"`
	DownloadMaxbr      int                `json:"downloadMaxbr"`
	MaxBrLevel         string             `json:"maxBrLevel"`
	PlayMaxBrLevel     string             `json:"playMaxBrLevel"`
	DownloadMaxBrLevel string             `json:"downloadMaxBrLevel"`
	PlLevel            string             `json:"plLevel"`
	DlLevel            string             `json:"dlLevel"`
	FlLevel            string             `json:"flLevel"`
	Rscl               interface{}        `json:"rscl"`
	FreeTrialPrivilege FreeTrialPrivilege `json:"freeTrialPrivilege"`
	RightSource        int                `json:"rightSource"`
	ChargeInfoList     []ChargeInfo       `json:"chargeInfoList"`
	Code               int                `json:"code"`
	Message            interface{}        `json:"message"`
	PlLevels           interface{}        `json:"plLevels"`
	DlLevels           interface{}        `json:"dlLevels"`
	IgnoreCache        interface{}        `json:"ignoreCache"`
	Bd                 interface{}        `json:"bd"`
}

// FreeTrialPrivilege 对应 privilege.freeTrialPrivilege
type FreeTrialPrivilege struct {
	ResConsumable      bool        `json:"resConsumable"`
	UserConsumable     bool        `json:"userConsumable"`
	ListenType         interface{} `json:"listenType"`
	CannotListenReason interface{} `json:"cannotListenReason"`
	PlayReason         interface{} `json:"playReason"`
	FreeLimitTagType   interface{} `json:"freeLimitTagType"`
}

// ChargeInfo 对应 privilege.chargeInfoList 数组项
type ChargeInfo struct {
	Rate          int         `json:"rate"`
	ChargeURL     interface{} `json:"chargeUrl"`
	ChargeMessage interface{} `json:"chargeMessage"`
	ChargeType    int         `json:"chargeType"`
}

// Album 顶层 album 字段
type Album struct {
	Songs           []interface{}  `json:"songs"`
	Paid            bool           `json:"paid"`
	OnSale          bool           `json:"onSale"`
	Mark            int            `json:"mark"`
	AwardTags       interface{}    `json:"awardTags"`
	DisplayTags     interface{}    `json:"displayTags"`
	Artists         []AlbumArtist  `json:"artists"`
	CopyrightID     int            `json:"copyrightId"`
	PicID           int64          `json:"picId"`
	Artist          AlbumArtist    `json:"artist"`
	PublishTime     int64          `json:"publishTime"`
	Company         string         `json:"company"`
	BriefDesc       string         `json:"briefDesc"`
	PicURL          string         `json:"picUrl"`
	CommentThreadID string         `json:"commentThreadId"`
	BlurPicURL      string         `json:"blurPicUrl"`
	CompanyID       int            `json:"companyId"`
	Pic             int64          `json:"pic"`
	Alias           []string       `json:"alias"`
	Status          int            `json:"status"`
	SubType         string         `json:"subType"`
	Description     string         `json:"description"`
	Tags            string         `json:"tags"`
	Name            string         `json:"name"`
	ID              int            `json:"id"`
	Type            string         `json:"type"`
	Size            int            `json:"size"`
	PicIDStr        string         `json:"picId_str"`
	Info            AlbumInfoBlock `json:"info"`
}

// AlbumArtist 对应 album.artists 数组项 以及 album.artist 对象
type AlbumArtist struct {
	Img1v1ID    int64    `json:"img1v1Id"`
	TopicPerson int      `json:"topicPerson"`
	PicID       int64    `json:"picId"`
	MusicSize   int      `json:"musicSize"`
	AlbumSize   int      `json:"albumSize"`
	BriefDesc   string   `json:"briefDesc"`
	PicURL      string   `json:"picUrl"`
	Img1v1URL   string   `json:"img1v1Url"`
	Followed    bool     `json:"followed"`
	Trans       string   `json:"trans"`
	Alias       []string `json:"alias"`
	Name        string   `json:"name"`
	ID          int      `json:"id"`
	Img1v1IDStr string   `json:"img1v1Id_str"`
	PicIDStr    string   `json:"picId_str,omitempty"`
}

// AlbumInfoBlock 对应 album.info
type AlbumInfoBlock struct {
	CommentThread    CommentThread `json:"commentThread"`
	LatestLikedUsers interface{}   `json:"latestLikedUsers"`
	Liked            bool          `json:"liked"`
	Comments         interface{}   `json:"comments"`
	ResourceType     int           `json:"resourceType"`
	ResourceID       int           `json:"resourceId"`
	CommentCount     int           `json:"commentCount"`
	LikedCount       int           `json:"likedCount"`
	ShareCount       int           `json:"shareCount"`
	ThreadID         string        `json:"threadId"`
}

// CommentThread 对应 album.info.commentThread
type CommentThread struct {
	ID               string          `json:"id"`
	ResourceInfo     CommentResource `json:"resourceInfo"`
	ResourceType     int             `json:"resourceType"`
	CommentCount     int             `json:"commentCount"`
	LikedCount       int             `json:"likedCount"`
	ShareCount       int             `json:"shareCount"`
	HotCount         int             `json:"hotCount"`
	LatestLikedUsers interface{}     `json:"latestLikedUsers"`
	ResourceOwnerID  int             `json:"resourceOwnerId"`
	ResourceTitle    string          `json:"resourceTitle"`
	ResourceID       int             `json:"resourceId"`
}

// CommentResource 对应 album.info.commentThread.resourceInfo
type CommentResource struct {
	ID        int         `json:"id"`
	UserID    int         `json:"userId"`
	Name      string      `json:"name"`
	ImgURL    string      `json:"imgUrl"`
	Creator   interface{} `json:"creator"`
	EncodedID interface{} `json:"encodedId"`
	SubTitle  interface{} `json:"subTitle"`
	WebURL    interface{} `json:"webUrl"`
}

// LyricResponse 对应 lyric.json 顶层结构
type LyricResponse struct {
	Sgc       bool       `json:"sgc"`
	Sfy       bool       `json:"sfy"`
	Qfy       bool       `json:"qfy"`
	LyricUser LyricUser  `json:"lyricUser"`
	Lrc       LyricBlock `json:"lrc"`
	KLyric    LyricBlock `json:"klyric"`
	TLyric    LyricBlock `json:"tlyric"`
	RomaLrc   LyricBlock `json:"romalrc"`
	Code      int        `json:"code"`
}

// LyricUser 对应 lyricUser
type LyricUser struct {
	ID       int    `json:"id"`
	Status   int    `json:"status"`
	Demand   int    `json:"demand"`
	UserID   int    `json:"userid"`
	Nickname string `json:"nickname"`
	Uptime   int64  `json:"uptime"`
}

// LyricBlock 对应 lrc/klyric/tlyric/romalrc
type LyricBlock struct {
	Version int    `json:"version"`
	Lyric   string `json:"lyric"`
}
