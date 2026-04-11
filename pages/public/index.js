function main() {
	const player = document.getElementById("player");

	/**
	 * @type {HTMLButtonElement}
	 */
	const playerStartButton = document.getElementById("start-player");
	playerStartButton.addEventListener("click", () => {
		playerStartButton.remove();
		startSongReporting();
		setInterval(() => {
			htmx.trigger(player, "playerReloadEvent");
		}, 5000);
	});


	function startSongReporting() {
		/**
		 * @param duration {number}
		 */
		function reportSongLength(duration) {
			console.debug("Reporting song duration");
			fetch(`/reportSongDuration/{{.CurrentSong.ID}}?duration=${duration}`, {
				method: "POST",
				credentials: "same-origin",
			}).catch((error) => {
				console.error("Failed to report song duration", error);
			});
		}

		setInterval(() => {
			/**
			 * @type {HTMLAudioElement | null}
			 */
			const audioPlayer = document.querySelector("#audio-player");
			if (!audioPlayer) {
				console.debug("No audio player, probably because no current song.");
				return;
			}


			if (audioPlayer.duration > 0) {
				reportSongLength(audioPlayer.duration);
			}
		}, 5000);
	}
}

if (document.readyState === "loading") {
	// Loading hasn't finished yet
	document.addEventListener("DOMContentLoaded", main);
} else {
	// `DOMContentLoaded` has already fired
	main();
}
