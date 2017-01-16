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

#### Dependancies needed to run
* [ffmpeg](https://ffmpeg.org/download.html) is required to be installed as it's used for transcoding different audio types to the target mp3s, suitable for uploading to SAM Broadcaster Cloud.  Install the correct version for your platform.  It uses 128k as it's bitrate.
* The SAM Library Import Utility must be in the "SAM" subdirectory from where you are running this tool.  Since SAM Broadcaster Cloud doesn't have a standalone API We utilize SAM's _ImportUtil_ and _PlaylistUtil_ utilities for actually interacting with your station. A working Java environment is required for these tools and details can be found at SAM's [Library Import Utility](http://spacial.com/library-import-utility/) page.

### Downloads
While it makes sense just to build this yourself, I am providing [Linux and macOS binaries ](https://github.com/gabek/samcloud-updater/releases), though they may not be as up to date as what's in master.  Given the wide variety of Linux distros and configurations I'm not sure if the Linux binary will work for everybody, but let me know if it doesn't.

### Build
1. go get github.com/gabek/samcloud-updater
2. cd $GOPATH/src/github.com/gabek/samcloud-updater
3. go build

Feel free to move the resulting `samcloud-updater` binary anywhere you like, as long as the `SAM` and `conf` directories go with it.

#### Audio storage directories
If not already created the _downloads_ and _uploads_ directories will be created.  This is where the audio files are stored.  This is also how the tool determines if a file has already been previously downloaded.  It's up to your own judgement how often you want to clean up these directories, if ever.  Mixcloud downloads, in particular, take a very long time, so clearing out the _downloads_ directory too often and making the tool think there are new files, when there really aren't, is wasteful.  However we only ever care about the *most recent* installment of content, so it's perfectly safe to get rid of old ones.  In addition, once the files from _uploads_ are uploaded, those files aren't needed anymore unless you want to keep them around.

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
      - title: Funky Track Man
        username: funkytrackman

station:
  username: "your@email.address"
  password: yourPassword
  id: 12345
  playlist: Mixes
  ```
