$(function() {
  $("#toggle-pin").click(function(){
    $.post("/control");
  })
});
