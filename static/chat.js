// axios.post("/auth/", "assda").then(res => {
// 	console.log(res);
// });

const socket = io("/");

let btn = document.querySelector("#btn");
let usr = document.querySelector("#usr");
let usrMsg = document.querySelector("#usrMsg");
let board = document.querySelector("#board");
let typing = document.querySelector("#typing");

const postMsg = () => {
	let data = { name: usr.value, message: usrMsg.value };
	socket.emit("msg", JSON.stringify(data));
	usrMsg.value = "";
};

const apeendMsgToScreen = data => {
	board.innerHTML += `<p><strong>${data.name} :</strong>${data.message}</p>`;
	typing.innerHTML = "";
};

const appendTypingToScreen = name => {
	typing.innerHTML = `<p>${name} is typing...</p>`;
};

const postTyping = () => {
	socket.emit("typing", usr.value);
};

btn.addEventListener("click", postMsg);
usrMsg.addEventListener("keyup", e => {
	if (e.keyCode === 13) {
		e.preventDefault();
		postMsg();
	} else {
		postTyping();
	}
});

socket.on("msg", data => apeendMsgToScreen(data));
socket.on("typing", data => appendTypingToScreen(data));
