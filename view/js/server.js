$(document).ready(function(){
	var username;
	$("#register").on("submit", function(e){
		e.preventDefault();
		username = $("#username").val();

		$.ajax({
			type:"POST",
			url:"http://localhost:8000/Validate",
			data:{
				"username": username
			},
			success: function(data){
				result(data);
			}
		})
	})

	function result(data){
		var obj = JSON.parse(data);
		if(obj.isvalid == true){
			console.log("Creating connection");
			createConnection();
		}else{
			console.log("No connection.");
		}
	}

	function createConnection(){
		var connection = new WebSocket("ws:://localhost:8000/Chat/" + username)
	}
})