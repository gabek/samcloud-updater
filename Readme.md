## SAM Broadcaster Cloud Updater


#### Goals
* Keep the station content fresh by supplying it with updated material.
* Utilize the services that people are already using to share content.
* Make it easy to use for broadcasters.
* Replace old versions of the content with new versions on the station.

#### Supported Services
* Podcasts (RSS + Enclosures)
* Mixcloud

#### Dependancies to run
* [ffmpeg](https://www.ffmpeg.org/download.html) is required to be installed as it's used for transcoding different audio types to the target mp3s, suitable for uploading to SAM Broadcaster Cloud.  Install the correct version for your platform.
* The SAM Library Import Utility must be in the "SAM" subdirectory from where you are running this tool.  Since SAM Broadcaster Cloud doesn't have a standalone API We utilize SAM's _ImportUtil_ and _PlaylistUtil_ utilities for actually interacting with your station. A working Java environment is required for these tools and details can be found at SAM's [Library Import Utility](http://spacial.com/library-import-utility/) page.

#### Audio storage directories
If not already created the _audio_ and _transcode_ directories will be created.  This is where the audio files are stored.  This is also how the tool determines if a file has already been previously downloaded.  It's up to your own judgement how often you want to clean up these directories, if ever.  Mixcloud downloads, in particular, take a very long time, so clearing out the directory too often and making the tool think there are new files, when there really aren't, is wasteful.  However we only ever care about the *most recent* installment of content, so it's perfectly safe to get rid of old ones.

#### Configuration
A config file is used to specify the details of the content you want to import, and how to access your station.  The station ID is found in the URL when logging into your station on the web and the playlist is the name of your Playlist you want the content added to.
```
podcasts:
    - title: The Requiem Podcast
      url: http://therequiem.libsyn.com/rss
    - title: Communion After Dark
      url: http://communionafterdark.com/rss.xml

mixcloud:
    - title: DJ Dark Machine
      username: djdarkmachine
    - title: DJ Led Manville
      username: djledmanville

station:
  username: "your@email.address"
  password: yourPassword
  id: 12345
  playlist: Mixes
  ```
