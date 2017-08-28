## SAM Broadcaster Cloud Updater


#### Goals
* Keep the station content fresh by supplying it with updated material.
* Utilize the services that people are already using to share content.
* Make it easy to use for broadcasters.
* Replace old versions of the content with new versions on the station.

#### Supported Services
* ##### Podcasts (RSS + Enclosures)
    The most recent episode will download and overwrite the previous episode on your station.

* ##### Mixcloud Users
    The most recent mix by each user listed will be downloaded and overwrite the previous mix by this user on your station.

* ##### Mixcloud Tags
    The most recent _"popular"_ mix that is listed under the tag specified will be downloaded and overwrite the previous mix for this tag on your station.

* ##### Soundcloud Users
    The most recent track by each user listed will be downloaded and overwrite the previous track by this user on your station.

    You'll need the user id number for each Soundcloud user in your config file. You can search for it [here](https://helgesverre.com/soundcloud/).

#### Dependancies needed to run
* [ffmpeg](https://ffmpeg.org/download.html) is required to be installed as it's used for transcoding different audio types to the target mp3s, suitable for uploading to SAM Broadcaster Cloud.  Install the correct version for your platform.  It uses 128k as it's bitrate.
* The SAM Library Import Utility must be in the "SAM" subdirectory from where you are running this tool.  Since SAM Broadcaster Cloud doesn't have a standalone API We utilize SAM's _ImportUtil_ and _PlaylistUtil_ utilities for actually interacting with your station. A working Java environment is required for these tools and details can be found at SAM's [Library Import Utility](http://spacial.com/library-import-utility/) page.

### Build
1. go get github.com/gabek/samcloud-updater
2. cd $GOPATH/src/github.com/gabek/samcloud-updater
3. go build

Feel free to move the resulting `samcloud-updater` binary anywhere you like, as long as the `SAM` and `conf` directories go with it.

#### Audio storage directories
If not already created the _downloads_ and _uploads_ directories will be created.  This is where the original and transcoded audio files are stored.  Feel free to clear out the contents of both of these directories after a run of this utility, however if you want to keep
the original files, the transcoded versions, or both for your audio library, then they're available there for you.

#### Configuration
A config file is used to specify the details of the content you want to import, and how to access your station.  The station ID is found in the URL when logging into your station on the web and the playlist is the name of your Playlist you want the content added to.
```
podcasts:
    - title: The Requiem Podcast
      url: http://therequiem.libsyn.com/rss
    - title: Communion After Dark
      url: http://communionafterdark.com/rss.xml

mixcloudusers:
    - title: DJ Dark Machine
      username: djdarkmachine
    - title: DJ Led Manville
      username: djledmanville

mixcloudtags:
    - futurepop

soundcloudusers:
      - username: Funky Track Man
        userId: 12345

station:
  username: "your@email.address"
  password: yourPassword
  id: 12345
  playlist: Mixes

userAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.101 Safari/537.36"
downloadLog: "./previouslyDownloaded.db"
  ```
