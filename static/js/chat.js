const socket = io("/");

const t = {
	nsp: "ness",
	auth: "authtoken",
	isAdmin: "true",
};

let btn = document.querySelector("#btn");
let usr = document.querySelector("#usr");
let usrMsg = document.querySelector("#usrMsg");
let board = document.querySelector("#board");
let typing = document.querySelector("#typing");
let login = document.querySelector("#login");

const auth = connected => {
	if (connected) socket.emit("auth", JSON.stringify(t));
};

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
	auth(true);
});

socket.on("chatData", data => loadChatData(data));
socket.on("msg", data => apeendMsgToScreen(data));
socket.on("typing", data => appendTypingToScreen(data));
socket.on("auth", data => auth(data));
