/**
 * @param songId {number}
 */
function startSongReporting(songId) {
	let interval = 0;

	/**
	 * @param duration {number}
	 */
	function reportSongLength(duration) {
		console.debug("Reporting song duration");
		fetch(`/reportSongDuration/${songId}?duration=${duration}`, {
			method: "POST",
			credentials: "same-origin",
		}).then(() => {
			console.debug("Succeeded reporting, stopping interval");
			clearInterval(interval);
		}).catch((error) => {
			console.error("Failed to report song duration", error);
		});
	}

	interval = setInterval(() => {
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

function initPlayer() {
	/**
	 * @type {HTMLAudioElement | null}
	 */
	const audioPlayer = document.querySelector("#audio-player");
	if (!audioPlayer) {
		console.debug("No audio player, probably because no current song.");
		return;
	}

	if (audioPlayer.dataset.initialized === "true") {
		return;
	}
	audioPlayer.dataset.initialized = "true";

	audioPlayer.volume = getVolume();
	void audioPlayer.play();
	audioPlayer.onended = () => {
		console.debug("Song ended, reloading");
		htmx.trigger(audioPlayer, "playerReloadEvent");
	};

	const currentTime = Number(audioPlayer.dataset.currentTime);
	if (Math.abs(audioPlayer.currentTime - currentTime) < 2) {
		// not that far off, ignore
	} else {
		if (currentTime) {
			audioPlayer.currentTime = currentTime;
		}
	}

	const knownLength = Number(audioPlayer.dataset.duration);
	if (!knownLength || isNaN(knownLength)) {
		startSongReporting(Number(audioPlayer.dataset.songId));
	}
}

function getVolume() {
	const defaultVolume = 0.10;

	const storedVolume = localStorage.getItem("volume");
	if (storedVolume == null) {
		return defaultVolume;
	}
	const volume = Number(storedVolume);
	if (isNaN(volume)) {
		return defaultVolume;
	}

	return volume;
}


function initVolume() {
	/**
	 * @type {HTMLInputElement}
	 */
	const volumeSlider = document.getElementById("volumeSlider");
	/**
	 * @type {HTMLOutputElement}
	 */
	const volumeSliderOutput = document.getElementById("volumeSliderOutput");


	volumeSlider.addEventListener("input", (event) => {
		const volume = event.target.value;
		const volumeNumeric = Number(volume) / 100;
		if (isNaN(volumeNumeric)) {
			return;
		}

		localStorage.setItem("volume", volumeNumeric.toString());

		const audioPlayer = document.querySelector("#audio-player");
		if (audioPlayer) {
			audioPlayer.volume = volumeNumeric;
		}

		if (volumeSliderOutput) {
			volumeSliderOutput.textContent = volume;
		}
	});
}

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

	setInterval(() => {
		initPlayer();
	}, 1000);
}

if (document.readyState === "loading") {
	// Loading hasn't finished yet
	document.addEventListener("DOMContentLoaded", main);
} else {
	// `DOMContentLoaded` has already fired
	main();
	initVolume();
}
