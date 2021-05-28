const socket = io();

let btn = document.querySelector("#btn");
let usr = document.querySelector("#usr");
let usrMsg = document.querySelector("#usrMsg");
let board = document.querySelector("#board");

const postMsg = () => {
	let data = { name: usr.value, message: usrMsg.value };
	socket.emit("msg", JSON.stringify(data));
	usrMsg.value = "";
};

const apeendMsgToScreen = data => {
	board.innerHTML += `<p><strong>${data.name} :</strong>${data.message}</p>`;
};

btn.addEventListener("click", postMsg);
usrMsg.addEventListener("keyup", e => {
	if (event.keyCode === 13) {
		event.preventDefault();
		postMsg();
	}
});

socket.on("msg", data => apeendMsgToScreen(data));
