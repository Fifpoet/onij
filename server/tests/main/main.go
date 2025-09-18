package main

import "onij/util/boost/tool"


func main() {
	album, err := GetOriginAlbum()
	if err != nil {
		panic(err)
	}
	print(tool.ToIndentJson(album.Songs[0]))
}
