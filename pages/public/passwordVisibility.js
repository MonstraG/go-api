/**
 * @param {HTMLButtonElement} button
 */
function initVisibilityButton(button) {
	const targetId = button.dataset["targetId"];
	if (!targetId) {
		console.error("Target id missing for password visibility button");
		return;
	}

	/** @type {HTMLInputElement} */
	const target = document.getElementById(targetId);
	if (!target) {
		console.error(`Target not found by id ${targetId}`);
		return;
	}

	/**
	 * @returns {HTMLImageElement}
	 */
	function createImg() {
		const showImage = document.createElement("img");
		showImage.width = 32;
		showImage.height = 32;
		return showImage;
	}

	const showImage = createImg();
	showImage.src = "/public/icons/visibility.svg";
	showImage.alt = "Show password";
	button.appendChild(showImage);

	const hideImage = createImg();
	hideImage.src = "/public/icons/visibility-off.svg";
	hideImage.alt = "Show password";
	if (target.type === "password") {
		hideImage.style.display = "none";
	}
	button.appendChild(hideImage);

	button.addEventListener("click", handleClick);

	function handleClick() {
		if (target.type === "password") {
			target.type = "text";
			showImage.style.display = "none";
			hideImage.style.display = "initial";
		} else {
			target.type = "password";
			showImage.style.display = "initial";
			hideImage.style.display = "none";
		}
	}
}

// I actually tried to do this with web components, but they have some problems:
//  1. safari doesn't support customized built-in components
//  2. if using autonomous, then css is not inherited, so I need to duplicate button styles
//  3. it's god-damn annoying to write even if it worked

/** @type {NodeListOf<HTMLButtonElement>} */
const buttons = document.querySelectorAll(".password-visibility");

for (const button of buttons) {
	initVisibilityButton(button);
}
