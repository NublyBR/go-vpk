# go-vpk
[![Go Report Card](https://goreportcard.com/badge/github.com/NublyBR/go-vpk)](https://goreportcard.com/report/github.com/NublyBR/go-vpk)
Golang implementation of Valve's Pak format

# Installation
```
$ go get -u github.com/NublyBR/go-vpk
```

# Examples
List all entries inside a `.vpk` dir file:
```go
package main

import (
	"fmt"

	"github.com/NublyBR/go-vpk"
)

func main() {
	// Open the VPK dir file
	pak, err := vpk.OpenDir(`C:\Program Files (x86)\Steam\steamapps\common\Half-Life 2\hl2\hl2_pak_dir.vpk`)
	if err != nil {
		panic(err)
	}
	defer pak.Close()

	// Iterate through all files in the VPK
	for _, file := range pak.Entries() {
		// Print the file size and full file name
		fmt.Printf("% 8d %s\n", file.Length(), file.Filename())
	}
}
```
Example output:
```
    1517 gamepadui/schemetab.res
    2056 gamepadui/schemesavebutton.res
    3348 gamepadui/schemepanel.res
    1233 gamepadui/schemeoptions_wheelywheel.res
    1440 gamepadui/schemeoptions_slideyslide.res
    3206 gamepadui/schemeoptions_skillyskill.res
    1150 gamepadui/schemeoptions_sectiontitle.res
    2038 gamepadui/schemeoptions_checkybox.res
    2691 gamepadui/schememainmenu.res
    2051 gamepadui/schemechapterbutton.res
    2356 gamepadui/schemeachievement.res
    8817 gamepadui/options.res
    1074 gamepadui/mainmenu.res
   18563 whitelist.cfg
     301 unusedcontent.cfg
    5595 shader_cache.cfg
  997033 scenes/scenes.image
   81317 scene.cache
      16 modelsounds.cache
    6062 maps/graphs/intro.ain
    7350 maps/graphs/d3_citadel_05.ain
       ...
```
