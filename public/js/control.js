$(function() {
  $("#toggle-pin").click(function(){
    $.post("/arduino/control");
  })
});
