const socket = io("/");

const admin = document.querySelector("#admin");
const usr = document.querySelector("#usr");
const usrName = document.querySelector("#usrName");
const usrBtn = document.querySelector("#usrBtn");
const waiting = document.querySelector("#waiting");

if (localStorage.getItem("isAdmin")) {
	const name = localStorage.getItem("adminName");
	const room = localStorage.getItem("room");
	const t = localStorage.getItem("token");
	const url = `${window.location.hostname}?r=${room}&t=${t}`;
	admin.innerHTML = `
		<h3>Hello ${name}</h3>
		<p>Click the link below to copy the link to dashboard and send it to chat memebers</p>
		<input type="text" id="link" value="${url.toString()}"/>
		<button id="start">Start Chat</button>
	`;
	const link = document.querySelector("#link");
	link.addEventListener("click", () => {
		link.select();
		link.setSelectionRange(0, 99999);
		document.execCommand("copy");
		alert("GOod");
	});
} else {
	usr.style.display = "block";
	usrBtn.addEventListener("click", () => {
		const urlParams = new URLSearchParams(window.location.search);
		const room = urlParams.get("r");
		const name = usrName.value;
		const data = JSON.stringify({ room, name });
		socket.emit("joinRoom", data);
	});
}

const addJoin = data => {
	console.log(data);
	waiting.innerHTML += `<p>${data} has joined the room</p>`;
};

socket.on("joined", data => addJoin(data));
//http://localhost:8000/waiting?r=chatRoom!!!!&t=aaabbb
