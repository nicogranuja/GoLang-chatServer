$(document).ready(function(){
	var username;
	var final_connection;

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
		$("#register").hide();
		$("#containerChat").show();

		var connection = new WebSocket("ws://localhost:8000/Chat/" + username)
		final_connection = connection;

		connection.onopen = function (response){
			connection.onmessage = function(response){
				console.log(response.data);
				var val = $("#chat_area").val();
				$("#chat_area").val(val + "\n" + response.data)
			}
		}

		$("#form_message").on("submit", function(e){
			e.preventDefault();
			var message = $("#msg").val();
			final_connection.send(message)
			$("#msg").val("");
		})

	}
})