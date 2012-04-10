function Errors(a, b, c) {
	alert("Cannot connect to server! It's probably closed :(");
}

//Fetch events
var FetchingEvents = false;
function GetEvents() {
	if (FetchingEvents) {
		return
	}
	
	FetchingEvents = true;
	$.ajax({
		url: "/events",
		success: HandleEvents,
		error: Errors,
		dataType: "json",
		type: "GET"
	});
}

function HandleEvents(evtobjs) {
	FetchingEvents = false
	if (evtobjs.error != undefined) {
		alert(k.error);
		return
	}
	for (i=0; i<evtobjs.length;i++) {
		handleEvent(evtobjs[i]);
	}	
}

function handleEvent(evt) {
	eval(evt.Javascript);
	if (evt.Reply) {
		$.ajax({
			url: "/reply",
			error: Errors,
			type: "GET",
			data: {"Id":evt.Id, "Reply": String(reply)}
		});
	}
}

function callHandler(id) {
	$.ajax({
		url: "/handler?id="+id,
		error: Errors,
	});
}
window.setInterval("GetEvents()", 10)

window.onbeforeunload = function() {
    return "Are you sure you want to quit?";
};
