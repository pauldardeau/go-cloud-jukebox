package jukebox


type SongDownloader struct {
   jukebox *Jukebox
   list_songs []*SongMetadata
}


func NewSongDownloader(jukebox *Jukebox,
                       list_songs []*SongMetadata) *SongDownloader {
    var sd SongDownloader;
    sd.jukebox = jukebox
    sd.list_songs = list_songs
    return &sd
}

func (sd *SongDownloader) run() {
    if sd.jukebox != nil && sd.list_songs != nil {
        sd.jukebox.batch_download_start()
        for _, song := range sd.list_songs {
            if sd.jukebox.exit_requested {
                break
            } else {
                sd.jukebox.download_song(song)
            }
        }
        sd.jukebox.batch_download_complete()
    }
}
