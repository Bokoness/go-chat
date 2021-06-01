const socket = io("/");

const t = {
	nsp: "ness",
	auth: "authtoken",
	isAdmin: "true",
};
// socket.emit("setRoom", JSON.stringify(t));

let btn = document.querySelector("#btn");
let usr = document.querySelector("#usr");
let usrMsg = document.querySelector("#usrMsg");
let board = document.querySelector("#board");
let typing = document.querySelector("#typing");
let login = document.querySelector("#login");

const loadChatData = data => {
	if (!data) return;
	board.innerHTML = "";
	data.forEach(val => {
		const msg = JSON.parse(val);
		board.innerHTML += `<p><strong>${msg.name} :</strong>${msg.message}</p>`;
	});
};

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

login.addEventListener("click", () => {
	socket.emit("setRoom", JSON.stringify(t));
});

socket.on("chatData", data => loadChatData(data));
socket.on("msg", data => apeendMsgToScreen(data));
socket.on("typing", data => appendTypingToScreen(data));
