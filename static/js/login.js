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
	let { status } = await axios.post("/auth", { name, token: t, nsp: room });
	if (status === 200) {
		localStorage.setItem("admin", JSON.stringify(t));
		//go to waiting room
	}
};

const invokeError = msg => {
	err.innerHTML = msg;
};

btn.addEventListener("click", checkToken);
