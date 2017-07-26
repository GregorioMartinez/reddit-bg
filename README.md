# reddit-bg
**Only works using gnome desktop**

Grab images from reddit and set as background.

## Installing
`go get github.com/GregorioMartinez/reddit-bg`

## Building
`go get -u github.com/GregorioMartinez/reddit-bg`
`cd $GOPATH/src/github.com/GregorioMartinez/reddit-bg`
`go build`

## Example Usage
`reddit-bg -w=1920 -subreddits=space,accidentalwesanderson,earthporn`

### Options
| Option | Description |
| --- | --- |
| subreddits | Comma separated list of subreddits to pull images from |
| t | Time frame to search. Can either be hour, day, week, month, year, or all |
| sort | Default sorting method of posts. Can be hot, new, rising, top, controversial, gilded, or promoted |
| w | Min width of image to use |
| h | Min height of image to use |
| limit | Number of posts to search |
| verbose | Verbose mode. Causes debugging information to be output |
| d | Directory to save images |
| r | - Randomly select an image from response |