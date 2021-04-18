# Shitposta GO

This a project was greatly inspired by [this one](https://github.com/alecs297/shitposta). The main reason why i made my own was that i wanted to learn Golang (and the fact that i hate node.js) and add a few feature to his project that weren't there at the time.

## Requirement

- [Golang](https://golang.org/) 1.16 or above
- [ffmpeg](https://www.ffmpeg.org/) (for gifs transcoding)

## Usage

Once evrything is installed just build it using `go build` and a folder data should be created. This folder should contain 3 sub folder (`img`, `video` and `new`) and a `shit.db` file. All new media added to the server should be added in the `new` folder wich is scaned every 5 minutes for new media. It will then move, rename and transcode them if necessary automaticly.

The web server is started on port `8090` by default
