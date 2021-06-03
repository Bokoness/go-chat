let adminName = document.querySelector("#adminName");
let token = document.querySelector("#token");
let btn = document.querySelector("#btn");
let roomName = document.querySelector("#roomName");
let err = document.querySelector("#err");

const checkToken = async () => {
	let name = adminName.value;
	let room = roomName.value;
	let t = token.value;
	if (!name || !t || !room) return invokeError("Please fill all fields");
	try {
		let { status } = await axios.post("/auth", {
			adminName: name,
			token: t,
			room: room,
		});
		if (status === 200) {
			localStorage.setItem("adminName", name);
			localStorage.setItem("room", room);
			localStorage.setItem("token", t);
			localStorage.setItem("isAdmin", true);
			window.location = `/waiting?t=${t}`;
		}
	} catch (e) {
		console.log(e.response);
	}
};

const invokeError = msg => {
	err.innerHTML = msg;
};

btn.addEventListener("click", checkToken);
