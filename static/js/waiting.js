const socket = io("/");

const admin = document.querySelector("#admin");
const usr = document.querySelector("#usr");
const usrName = document.querySelector("#usrName");
const usrBtn = document.querySelector("#usrBtn");
const waiting = document.querySelector("#waiting");
const adminTitle = document.querySelector("#adminTitle");

const urlParams = new URLSearchParams(window.location.search);
const room = urlParams.get("r");

if (localStorage.getItem("isAdmin")) {
	admin.style.display = "block";
	const name = localStorage.getItem("adminName");
	const t = localStorage.getItem("token");
	// const url = `${window.location.hostname}/waiting?r=${room}&t=${t}`;
	const url = `http://localhost:8000/waiting?r=${room}&t=${t}`;
	adminTitle.textContent = `Hello ${name}`;
	const link = document.querySelector("#link");
	link.value = url;
	link.addEventListener("click", () => {
		link.select();
		link.setSelectionRange(0, 99999);
		document.execCommand("copy");
		alert("Copied!");
	});
} else {
	usr.style.display = "block";
	usrBtn.addEventListener("click", () => {
		const name = usrName.value;
		const data = JSON.stringify({ room, name });
		localStorage.setItem("room", room);
		localStorage.setItem("name", name);
		socket.emit("joinRoom", data);
	});
}

const addJoin = data => {
	waiting.innerHTML = "";
	for (i in data) {
		console.log(data[i]);
		waiting.innerHTML += `<p>${data[i]}</p>`;
	}
};

//http://localhost:8000/waiting?r=chatRoom!!!!&t=aaabbb
console.log("loaded");
socket.on("joined", data => addJoin(data));
socket.on("getUsrs", data => console.log(data));

setTimeout(() => {
	let r = localStorage.getItem("room");
	let n = localStorage.getItem("name");
	if (r && n) {
		socket.emit("joinRoom", JSON.stringify({ room: r, name: n }));
		usr.style.display = "none";
	}
}, 500);
