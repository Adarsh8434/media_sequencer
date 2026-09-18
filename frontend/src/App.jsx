import { useEffect, useState } from "react";

const API_URL = import.meta.env.VITE_API_URL;

function MediaWindow({ windowData, syncState, media }) {
  const [playlist, setPlaylist] = useState([]);
  const [currentIndex, setCurrentIndex] = useState(0);
  const [selectedMediaId, setSelectedMediaId] = useState("");

  useEffect(() => {
    loadPlaylist();
  }, [windowData.id]);

  const loadPlaylist = async () => {
    try {
      const response = await fetch(
        `${API_URL}/windows/${windowData.id}/playlist`
      );

      if (!response.ok) {
        throw new Error("Failed to load playlist");
      }

      const data = await response.json();
      setPlaylist(data);
    } catch (error) {
      console.error("Failed to load playlist:", error);
    }
  };

  const addMedia = async () => {
    if (!selectedMediaId) {
      return;
    }

    try {
      const response = await fetch(`${API_URL}/playlist`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          windowId: windowData.id,
          mediaId: Number(selectedMediaId),
          position: playlist.length + 1,
        }),
      });

      if (!response.ok) {
        throw new Error("Failed to add media");
      }

      setSelectedMediaId("");
      await loadPlaylist();
    } catch (error) {
      console.error("Failed to add media:", error);
    }
  };

  // Current item from the normal playlist
  const currentPlaylistItem = playlist[currentIndex];

  // Find the actual media object
  const normalMedia = media.find(
    (item) => item.id === currentPlaylistItem?.mediaId
  );

  // Media selected for synchronization
  const syncMedia = syncState?.active
    ? media.find((item) => item.id === syncState.mediaId)
    : null;

  // During sync show sync media.
  // Otherwise show the normal playlist media.
  const currentMedia = syncMedia || normalMedia;

  // Normal playlist playback
  useEffect(() => {
    if (!normalMedia || syncState?.active || playlist.length === 0) {
      return;
    }

    const timer = setTimeout(() => {
      setCurrentIndex((previousIndex) => {
        const nextIndex = previousIndex + 1;

        // Restart playlist when the last item finishes
        if (nextIndex >= playlist.length) {
          return 0;
        }

        return nextIndex;
      });
    }, normalMedia.durationSeconds * 1000);

    return () => clearTimeout(timer);
  }, [
    currentIndex,
    normalMedia,
    playlist.length,
    syncState?.active,
  ]);

  // Prevent undefined media error
  if (!currentMedia) {
    return (
      <div className="media-window">
        <h2>{windowData.name}</h2>

        <div className="blank-screen">
          No media configured
        </div>

        <div className="playlist-controls">
          <select
            value={selectedMediaId}
            onChange={(e) => setSelectedMediaId(e.target.value)}
          >
            <option value="">Select Media</option>

            {media.map((item) => (
              <option key={item.id} value={item.id}>
                {item.name}
              </option>
            ))}
          </select>

          <button onClick={addMedia}>
            Add Media
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="media-window">
      <h2>{windowData.name}</h2>

      {currentMedia.type === "image" && (
        <img
          src={currentMedia.url}
          alt={currentMedia.name}
          className="media-content"
        />
      )}

      {currentMedia.type === "video" && (
        <video
          src={currentMedia.url}
          className="media-content"
          autoPlay
          muted
          playsInline
        />
      )}

      {currentMedia.type === "blank" && (
        <div className="blank-screen">
          Blank
        </div>
      )}

      <p>
        Playing: <strong>{currentMedia.name}</strong>
      </p>

      <div className="playlist-controls">
        <select
          value={selectedMediaId}
          onChange={(e) => setSelectedMediaId(e.target.value)}
        >
          <option value="">Select Media</option>

          {media.map((item) => (
            <option key={item.id} value={item.id}>
              {item.name}
            </option>
          ))}
        </select>

        <button onClick={addMedia}>
          Add Media
        </button>
      </div>

      {syncState?.active && (
        <p>
          <strong>SYNC PLAYBACK</strong>
        </p>
      )}
    </div>
  );
}

function App() {
  const [windows, setWindows] = useState([]);
  const [media, setMedia] = useState([]);
  const [syncState, setSyncState] = useState(null);

  const [syncMediaId, setSyncMediaId] = useState("");
  const [syncDuration, setSyncDuration] = useState(10);

  const [loading, setLoading] = useState(true);

  // Load windows and media
  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      const [windowsResponse, mediaResponse] = await Promise.all([
        fetch(`${API_URL}/windows`),
        fetch(`${API_URL}/media`),
      ]);

      if (!windowsResponse.ok || !mediaResponse.ok) {
        throw new Error("Failed to load application data");
      }

      const windowsData = await windowsResponse.json();
      const mediaData = await mediaResponse.json();

      setWindows(windowsData);
      setMedia(mediaData);

      if (mediaData.length > 0) {
        setSyncMediaId(mediaData[0].id);
      }
    } catch (error) {
      console.error("Failed to load data:", error);
    } finally {
      setLoading(false);
    }
  };

  // Check sync state every second
  useEffect(() => {
    const checkSync = async () => {
      try {
        const response = await fetch(`${API_URL}/sync`);

        if (!response.ok) {
          throw new Error("Failed to get sync state");
        }

        const data = await response.json();

        setSyncState(data);
      } catch (error) {
        console.error("Failed to get sync state:", error);
      }
    };

    checkSync();

    const interval = setInterval(checkSync, 1000);

    return () => clearInterval(interval);
  }, []);

  // Start synchronization
  const startSync = async () => {
    if (!syncMediaId) {
      return;
    }

    try {
      const response = await fetch(`${API_URL}/sync`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          mediaId: Number(syncMediaId),
          durationSeconds: Number(syncDuration),
        }),
      });

      if (!response.ok) {
        throw new Error("Failed to start sync");
      }

      const data = await response.json();

      console.log(data);
    } catch (error) {
      console.error("Failed to start sync:", error);
    }
  };

  if (loading) {
    return <h2>Loading...</h2>;
  }

  return (
    <div className="app">
      <h1>Media Sequencer</h1>

      {/* Sync controls */}
      <div className="sync-controls">
        <select
          value={syncMediaId}
          onChange={(e) => setSyncMediaId(e.target.value)}
        >
          {media.map((item) => (
            <option key={item.id} value={item.id}>
              {item.name}
            </option>
          ))}
        </select>

        <input
          type="number"
          min="1"
          value={syncDuration}
          onChange={(e) => setSyncDuration(e.target.value)}
        />

        <button onClick={startSync}>
          Start Sync
        </button>
      </div>

      {/* Four media windows */}
      <div className="windows-container">
        {windows.map((windowData) => (
          <MediaWindow
            key={windowData.id}
            windowData={windowData}
            syncState={syncState}
            media={media}
          />
        ))}
      </div>
    </div>
  );
}

export default App;