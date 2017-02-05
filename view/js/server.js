$(document).ready(function(){

	$("#register").on("submit", function(e){
		e.preventDefault();
		var username = $("#username").val();

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

	}
})