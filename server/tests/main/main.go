package main

import "onij/util/boost/tool"

func main() {
	album, err := FetchAlbum()
	if err != nil {
		panic(err)
	}
	print(tool.ToIndentJson(album.Songs[0]))
}
